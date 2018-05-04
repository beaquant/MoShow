package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 45 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 15 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

const (
	wsMessageTypeText         = iota //文本消息内容
	wsMessageTypeChannelStart        //房间初始化(双方进入房间，视频开始，计费开始)
	wsMessageTypeAllocateFund        //扣费信息(每分钟发送用户余额，消费信息，聊天时长)
	wsMessageTypeException           //异常挂断(扣费异常或者连接异常)
	wsMessageTypeChannelEnd          //房间结算(聊天结束后发送 结算信息)
	wsMessageTypeSystem              //系统消息
	wsMessageTypeJoinSuccess         //加入房间成功
	wsMessageTypeJoinFail            //关闭连接
)

var (
	upgrader     = websocket.Upgrader{}
	chatChannels = make(map[uint64]*ChatChannel)
	chattingUser = make(map[uint64]interface{})
)

//WebsocketController websocket相关
type WebsocketController struct {
	beego.Controller
}

//ChatChannel .
type ChatChannel struct {
	ID               uint64
	DialID           uint64
	NIMChannelID     uint64 //网易云信房间ID
	DstID            uint64
	Inited           bool
	Stoped           bool
	ChannelStartTime int64
	StartTime        int64
	Timelong         uint64 //聊天时长,单位:秒
	StopTime         int64
	Price            uint64
	Amount           uint64
	logger           *logrus.Entry
	Src              *ChatClient
	Dst              *ChatClient
	Send             chan *WsMessage
	Join             chan *ChatClient
	Exit             chan []error
}

//ChatClient .
type ChatClient struct {
	User    *models.UserProfile
	Channel *ChatChannel
	Conn    *websocket.Conn
	Request *http.Request
	Send    chan *WsMessage
}

//WsMessage .
type WsMessage struct {
	Content     string `json:"content"`
	MessageType int    `json:"type"`
	DialID      uint64 `json:"dial_id"`
}

//VideoCost 用户消费信息
type VideoCost struct {
	Balance      uint64 `json:"balance" description:"用户余额"`
	Cost         uint64 `json:"cost" description:"用户花费"`
	Income       uint64 `json:"income,omitempty" description:"主播收益"`
	Timelong     uint64 `json:"timelong" description:"聊天时长"`
	NIMChannelID uint64 `json:"NIMChannelID,omitempty" description:"网易云房间ID"`
}

func closeConnWithMessage(conn *websocket.Conn, ws *WsMessage) {
	ws.MessageType = wsMessageTypeJoinFail
	b, _ := utils.JSONMarshal(ws)

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	conn.WriteJSON(ws)
	conn.WriteMessage(websocket.CloseNormalClosure, b)
	beego.Info("服务器主动挂断,并推送消息")
	conn.Close()
}

//Create .
// @Title 创建聊天通道[websocket]
// @Description 创建聊天通道[websocket]
// @Param   parterid     path    int  true        "聊天对象的ID"
// @router /:parterid/create [get]
func (c *WebsocketController) Create() {
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		beego.Error(err)
		return
	}

	parter, err := strconv.ParseUint(c.Ctx.Input.Param(":parterid"), 10, 64)
	if err != nil {
		beego.Error(err)
		return
	}

	ws := &WsMessage{MessageType: wsMessageTypeException}
	tk := GetToken(c.Ctx)
	if _, ok := chattingUser[tk.ID]; ok { //如果用户正在聊天，拒绝创建聊天通道
		ws.Content = "您已在通话中"
		beego.Error(ws.Content, tk.ID, c.Ctx.Request.UserAgent())
		closeConnWithMessage(conn, ws)
		return
	}

	if _, ok := chattingUser[parter]; ok { //如果主播正在聊天，拒绝创建聊天通道
		ws.Content = "对方正在通话中"
		beego.Error(ws.Content, tk.ID, c.Ctx.Request.UserAgent())
		closeConnWithMessage(conn, ws)
		return
	}

	up := &models.UserProfile{ID: tk.ID}
	if err := up.Read(); err != nil {
		ws.Content = "获取用户信息失败\t" + err.Error()
		beego.Error("获取用户信息失败", err, tk.ID, c.Ctx.Request.UserAgent())
		closeConnWithMessage(conn, ws)
		return
	}

	pp := &models.UserProfile{ID: parter}
	if err := pp.Read(); err != nil {
		ws.Content = "获取用户信息失败\t" + err.Error()
		beego.Error("获取用户信息失败", err, parter, c.Ctx.Request.UserAgent())
		closeConnWithMessage(conn, ws)
		return
	}

	if pp.UserType != models.UserTypeAnchor {
		ws.Content = "对方不是主播,不能直播"
		closeConnWithMessage(conn, ws)
		return
	}

	if up.Balance+up.Income < pp.Price {
		ws.Content = "余额不足本次通话"
		beego.Error(ws.Content, tk.ID, c.Ctx.Request.UserAgent())
		closeConnWithMessage(conn, ws)
		return
	}

	dl := &models.Dial{FromUserID: tk.ID, ToUserID: parter}
	if err := dl.Add(); err != nil {
		ws.Content = "创建通话记录失败"
		beego.Error(ws.Content, tk.ID, c.Ctx.Request.UserAgent())
		closeConnWithMessage(conn, ws)
		return
	}

	//Exit通道要设置缓冲区，不然会在写Exit的时候死锁导致无法读Exit
	channel := &ChatChannel{Join: make(chan *ChatClient), Send: make(chan *WsMessage), Exit: make(chan []error, 1), ID: tk.ID, DstID: parter, DialID: dl.ID}
	go channel.Run()

	client := &ChatClient{User: up, Conn: conn, Send: make(chan *WsMessage), Request: c.Ctx.Request}
	client.Channel = channel
	chatChannels[tk.ID] = channel
	chattingUser[tk.ID] = nil

	go client.Read()
	go client.Write()

	channel.Join <- client //加入频道
}

//Join .
// @Title 加入聊天通道[websocket]
// @Description 加入聊天通道[websocket]
// @Param   channelid     path    int  true        "聊天频道的ID"
// @router /:channelid/join [get]
func (c *WebsocketController) Join() {
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		beego.Error(err)
		return
	}

	channelid, err := strconv.ParseUint(c.Ctx.Input.Param(":channelid"), 10, 64)
	if err != nil {
		beego.Error(err)
		return
	}

	ws := &WsMessage{MessageType: wsMessageTypeException}
	tk := GetToken(c.Ctx)
	cn, ok := chatChannels[channelid]
	if !ok || (cn.DstID != tk.ID && cn.ID != tk.ID) {
		ws.Content = "用户不属于该房间，禁止加入"
		beego.Error(ws.Content, tk.ID, c.Ctx.Request.UserAgent())
		closeConnWithMessage(conn, ws)
		return
	}

	if cn.StopTime != 0 {
		ws.Content = "房间已关闭,加入失败"
		beego.Error(ws.Content, tk.ID, c.Ctx.Request.UserAgent())
		closeConnWithMessage(conn, ws)
		return
	}

	var client *ChatClient
	if tk.ID == cn.ID {
		beego.Info("用户重连成功,尝试挂断旧连接")
		cn.Src.Conn.Close()
		cn.Src.Conn = conn
		client = cn.Src
	} else if tk.ID == cn.DstID {
		if cn.Dst == nil {
			up := &models.UserProfile{ID: tk.ID}
			if err := up.Read(); err != nil {
				if err != nil {
					beego.Error(err)
					return
				}
			}

			cn.Price = up.Price
			client = &ChatClient{User: up, Channel: cn, Conn: conn, Send: make(chan *WsMessage), Request: c.Ctx.Request}

			cn.Join <- client
		} else {
			beego.Info("主播重连成功,尝试挂断旧连接")
			cn.Dst.Conn.Close()
			cn.Dst.Conn = conn
			client = cn.Dst
		}
	}

	chattingUser[tk.ID] = nil

	go client.Read()
	go client.Write()
}

//Reject .
// @Title 拒绝聊天请求
// @Description 拒绝聊天请求
// @Param   channelid     path    int  true        "聊天对象的ID"
// @router /:channelid/reject [get]
func (c *WebsocketController) Reject() {
	dto := &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	channelid, err := strconv.ParseUint(c.Ctx.Input.Param(":channelid"), 10, 64)
	if err != nil {
		beego.Error(err)
		return
	}

	tk := GetToken(c.Ctx)
	cn, ok := chatChannels[channelid]
	if !ok || cn.DstID != tk.ID {
		dto.Message = "未找到指定频道，或者当前用户不在指定频道中"
		return
	}

	if !cn.Stoped { //未结算
		cn.Exit <- nil
		(&models.UserProfile{ID: tk.ID}).AddDialReject(nil)
	}

	dto.Sucess = true
}

func (c *ChatChannel) genVideoCost() (*VideoCost, error) {
	vc := &VideoCost{}
	up := (&models.UserProfile{ID: c.ID})
	if err := up.Read(); err != nil {
		return nil, err
	}

	vc.Balance = up.Balance + up.Income
	vc.Cost = c.Amount
	vc.Timelong = c.Timelong
	return vc, nil
}

//Run .
func (c *ChatChannel) Run() {
	defer c.CloseChannel()

	c.ChannelStartTime = time.Now().Unix()
	//初始化日志模块
	c.logger = logrus.WithFields(logrus.Fields{"dial_id": c.DialID})
	file, err := os.OpenFile(path.Join("logs", strconv.FormatUint(c.ID, 10)+"_ws.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		c.logger.Error("打开日志文件失败", err)
	} else {
		c.logger.Logger.Out = file
		defer file.Close()
	}

	for {
		select {
		case client := <-c.Join:
			if client.User.ID == c.ID {
				c.Src = client
			} else if client.User.ID == c.DstID {
				c.Dst = client
			}

			client.Send <- &WsMessage{MessageType: wsMessageTypeJoinSuccess, Content: "加入房间成功"}
			c.logger.Info("[uid:", client.User.ID, "]用户加入房间")
		case msg := <-c.Send:
			go c.wsMsgDeal(msg)
		case exp := <-c.Exit:
			c.Stoped = true

			if c.StartTime == 0 {
				c.StartTime = c.ChannelStartTime
			}

			if c.Dst == nil || (c.StopTime == 0 && !c.Inited) {
				c.logger.Info("视频未接通:用户挂断,用户ID:", c.ID, "主播ID:", c.DstID)
				return //主播未加入房间，直接退出 || 没有初始化聊天，直接退出
			}

			if c.StopTime == 0 { //未进入正常结算流程
				//房间初始化成功，但是没有进入正常结费流程(未收到 wsMessageTypeChannelEnd)
				c.StopTime = time.Now().Unix()
				c.Timelong = uint64(c.StopTime - c.StartTime)
			} else { //收到 wsMessageTypeChannelEnd 正常结算流程
				if !c.Inited { //没有正常开始聊天，需要补扣用户的钱
					if err := videoAllocateFund(c.Src.User, c.Dst.User, c.Amount); err != nil { //扣费失败关闭聊天频道
						c.logger.Error("扣费失败", err)
						exp = append(exp, err)
					}
				}
			}

			c.logger.Info("开始结算,房间ID:", c.ID, "主播ID:", c.DstID)
			if exp == nil || len(exp) == 0 { //没有异常的情况下再给主播结费，否则只生成异常通话记录
				if err := videoDone(c.Src.User, c.Dst.User, &models.VideoChgInfo{TimeLong: c.Timelong, Price: c.Price, DialID: c.DialID}, c.Amount); err != nil {
					c.logger.Error("[websocket结算异常]视频结费错误", err, "发起人:", c.ID, "接受人:", c.DstID, "金额:", c.Amount, "通话时长:", c.Timelong)
					exp = append(exp, errors.New("[websocket结算异常]视频结费错误:"+err.Error()))
				}
			}

			income, _, _ := computeIncome(c.Amount)
			ciStr, _ := utils.JSONMarshalToString(&models.ClearingInfo{NIMChannelID: c.NIMChannelID, Cost: c.Amount, Income: uint64(income), Price: c.Price, Timelong: c.Timelong})

			//生成通话记录
			dl, dt := &models.Dial{ID: c.DialID}, &models.DialTag{}
			if exp != nil {
				dl.Status = models.DialStatusException
				for _, val := range exp {
					dt.ErrorMsg = append(dt.ErrorMsg, val.Error())
				}
				dl.Tag, _ = utils.JSONMarshalToString(dt)
			} else {
				dl.Status = models.DialStatusSuccess
			}

			trans := models.TransactionGen()
			if err := dl.Update(map[string]interface{}{"duration": c.Timelong, "create_at": c.ChannelStartTime, "status": dl.Status, "clearing": ciStr}, trans); err != nil {
				js, _ := utils.JSONMarshalToString(dl)
				c.logger.Error("[websocket结算异常]通话记录更新失败", err, js)
				models.TransactionRollback(trans)
				return
			}

			if err := c.Src.User.AddDialDuration(c.Timelong, trans); err != nil {
				c.logger.Error("[websocket结算异常]用户增加通话时长失败", err, c.ID)
				models.TransactionRollback(trans)
				return
			}
			if err := c.Dst.User.AddDialDuration(c.Timelong, trans); err != nil {
				c.logger.Error("[websocket结算异常]用户增加通话时长失败", err, c.DstID)
				models.TransactionRollback(trans)
				return
			}

			models.TransactionCommit(trans)
			ms := &WsMessage{MessageType: wsMessageTypeChannelEnd, DialID: c.DialID}
			vc, _ := c.genVideoCost()
			vc.Income, vc.NIMChannelID = uint64(income), c.NIMChannelID
			ms.Content, _ = utils.JSONMarshalToString(vc)
			c.Src.Send <- ms
			c.Dst.Send <- ms
			c.logger.Info("房间结算成功")
			return
		}
	}
}

func (c *ChatChannel) wsMsgDeal(msg *WsMessage) {
	defer func() {
		if err := recover(); err != nil {
			c.logger.Error(err)
			debug.PrintStack()
		}
	}()

	switch msg.MessageType {
	case wsMessageTypeChannelStart:
		if c.Dst == nil {
			c.Src.Send <- &WsMessage{MessageType: wsMessageTypeSystem, Content: "主播还未进入房间"}
			break
		}

		vcp := &VideoCost{}
		if err := utils.JSONUnMarshal(msg.Content, vcp); err == nil {
			if c.NIMChannelID == 0 {
				c.NIMChannelID = vcp.NIMChannelID
			}
		} else {
			c.logger.Error("解析客户端消息参数失败", msg.Content, err)
		}

		if !c.Inited { //双方进入房间，初始化房间，开始视频
			c.Inited = true
			c.StartTime = time.Now().Unix()

			c.logger.Info("视频开始,开始扣费")
			c.Amount += c.Price
			//扣费
			if err := videoAllocateFund(c.Src.User, c.Dst.User, c.Price); err != nil {
				c.logger.Error("扣费失败", err)
				if !c.Stoped {
					c.Exit <- nil
				}
			}

			go c.ticktokPay()
		}

		vc, err := c.genVideoCost()
		if err != nil {
			c.logger.Error("生成消费信息失败", err)
			if !c.Stoped {
				c.Exit <- nil
			}
		}

		m := &WsMessage{MessageType: wsMessageTypeChannelStart}
		m.Content, _ = utils.JSONMarshalToString(vc)
		c.Src.Send <- m
		c.Dst.Send <- m
	case wsMessageTypeChannelEnd:
		if c.StopTime != 0 {
			break
		}

		c.StopTime = time.Now().Unix()

		var errs []error
		vcp := &VideoCost{}
		if err := utils.JSONUnMarshal(msg.Content, vcp); err == nil {
			if c.NIMChannelID == 0 {
				c.NIMChannelID = vcp.NIMChannelID
			}

			if !c.Inited { //未收到房间初始化信息，直接进入结费,按客户端传入的时长计算结费信息
				if vcp.Timelong == 0 || vcp.Timelong > uint64((c.StopTime-c.ChannelStartTime)*2) {
					vcp.Timelong = uint64(c.StopTime - c.ChannelStartTime)
				}
				c.StartTime = c.StopTime - int64(vcp.Timelong)
				c.Timelong = vcp.Timelong
				c.Amount = c.Price * ((vcp.Timelong + 59) / 60)
			} else {
				c.Timelong = uint64(c.StopTime - c.StartTime)
			}
		} else {
			errs = append(errs, errors.New("解析结费请求参数错误:"+msg.Content))
			c.logger.Error("解析结费请求参数错误", msg.Content, msg.MessageType, msg.Content)
		}
		c.Exit <- errs
	}
}

func (c *ChatChannel) ticktokPay() {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		if err := recover(); err != nil {
			c.logger.Error(err)
			debug.PrintStack()
		}
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			if c.StopTime == 0 { //若频道未关闭,每隔60秒扣费一次
				c.Timelong = uint64(time.Now().Unix() - c.StartTime)

				//扣费
				if err := videoAllocateFund(c.Src.User, c.Dst.User, c.Price); err != nil { //扣费失败关闭聊天频道
					c.logger.Error(err)
					m := &WsMessage{MessageType: wsMessageTypeSystem, Content: "用户扣费失败"}
					c.Src.Send <- m
					c.Dst.Send <- m

					time.Sleep(time.Second)
					if !c.Stoped {
						c.Exit <- nil
					}
					return
				}
				c.Amount += c.Price

				vc, err := c.genVideoCost()
				if err != nil {
					c.logger.Error("生成消费信息失败", err)
					if !c.Stoped {
						c.Exit <- nil
					}
					return
				}

				m := &WsMessage{MessageType: wsMessageTypeAllocateFund}
				m.Content, _ = utils.JSONMarshalToString(vc)
				c.Src.Send <- m
				c.Dst.Send <- m
			} else {
				break
			}
		}
	}
}

//CloseChannel 关闭频道,通道,websocket链接
func (c *ChatChannel) CloseChannel() {
	if err := recover(); err != nil {
		beego.Error(err)
		debug.PrintStack()
	}

	close(c.Send)
	close(c.Exit)
	close(c.Join)
	delete(chatChannels, c.ID)
	delete(chattingUser, c.ID)
	delete(chattingUser, c.DstID)

	if c.Src != nil {
		close(c.Src.Send)
	}

	if c.Dst != nil {
		close(c.Dst.Send)
	}
}

func (c *ChatClient) Read() {
	curConnection := c.Conn

	defer func() {
		if err := recover(); err != nil {
			beego.Error(err)
			debug.PrintStack()
		}

		if c.Channel.StopTime == 0 { //聊天通道未结束
			// time.Sleep(30 * time.Second) //等待30秒,如过连接没有恢复,执行退出

			if curConnection != c.Conn { //重连成功
				c.Channel.logger.Info("重连成功,房间继续保留,房间ID:", c.Channel.ID)
				return
			}
		}

		c.Channel.logger.Info(fmt.Sprintf("[ws:%p]", c.Conn), "用户主动挂断或等待重连超时,执行退出流程")
		if !c.Channel.Stoped { //未结算
			c.Channel.Exit <- nil //发送退出信号,关闭通道后write方法会立即退出
		}
	}()
	curConnection.SetReadLimit(maxMessageSize)
	curConnection.SetReadDeadline(time.Now().Add(pongWait))
	curConnection.SetPongHandler(func(pong string) error {
		c.Channel.logger.Info("[uid:", c.User.ID, "]收到pong消息,刷新deadline")

		if _, ok := chatChannels[c.Channel.ID]; ok {
			curConnection.SetReadDeadline(time.Now().Add(pongWait))
		}
		return nil
	})

	for {
		m := &WsMessage{}
		_, message, err := curConnection.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.Channel.logger.Error("客户端异常挂断", err, "用户ID:", c.User.ID, "Agent:", c.Request.UserAgent())
				} else {
					c.Channel.logger.Info("客户端主动挂断", err, "用户ID:", c.User.ID, "Agent:", c.Request.UserAgent())
				}
			} else {
				// time.Sleep(pongWait - pingPeriod) //等待重连
			}

			break
		}

		err = utils.JSONUnMarshalFromByte(message, m)
		if err == nil {
			if c.Channel.StopTime == 0 {
				c.Channel.logger.Info("收到客户端消息:", string(message))
				c.Channel.Send <- m
			}
		} else {
			c.Channel.logger.Error("[异常(消息格式错误,忽略该消息)]收到客户端消息:", string(message), "错误信息:", err)
		}
	}
}

func (c *ChatClient) Write() {
	ticker := time.NewTicker(pingPeriod)
	curConnection := c.Conn
	di := c.Channel.DialID

	defer func() {
		if err := recover(); err != nil {
			beego.Error(err)
			debug.PrintStack()
		}

		ticker.Stop()

		if c.Channel.StopTime == 0 { //聊天通道未结束
			// time.Sleep(pongWait - pingPeriod) //等待重连
		}

		c.Channel.logger.Warn(fmt.Sprintf("[ws:%p]", c.Conn), "写入超时,尝试挂断,用户ID:", c.User.ID)
		curConnection.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if message != nil {
				message.DialID = di
			}

			curConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}

			ms, _ := utils.JSONMarshalToString(message)
			c.Channel.logger.Info("发送消息", ms)

			err := curConnection.WriteJSON(message)
			if err != nil {
				return
			}
		case <-ticker.C:
			curConnection.SetWriteDeadline(time.Now().Add(writeWait))

			c.Channel.logger.Info("[uid:", c.User.ID, "]执行ping操作")
			if err := curConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Channel.logger.Error("[uid:", c.User.ID, "]执行ping操作错误")
				return
			}
		}
	}
}
