package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 45 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 15 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

const (
	wsMessageTypeText         = iota //文本消息内容
	wsMessageTypeChannelStart        //房间初始化
	wsMessageTypeAllocateFund        //扣费信息
	wsMessageTypeException           //异常挂断
	wsMessageTypeChannelEnd          //房间结算
	wsMessageTypeSystem              //系统消息
	wsMessageTypeClose               //关闭房间
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
	ID        uint64
	DstID     uint64
	Inited    bool
	StartTime int64
	Timelong  uint64 //聊天时长,单位:秒
	StopTime  int64
	Price     uint64
	Amount    uint64
	Src       *ChatClient
	Dst       *ChatClient
	Send      chan *WsMessage
	Join      chan *ChatClient
	Exit      chan error
}

//ChatClient .
type ChatClient struct {
	User    *models.UserProfile
	Channel *ChatChannel
	Conn    *websocket.Conn
	Send    chan *WsMessage
}

//WsMessage .
type WsMessage struct {
	Content     string `json:"content"`
	MessageType int    `json:"type"`
}

//VideoCost 用户消费信息
type VideoCost struct {
	Balance  uint64 `json:"balance" description:"用户余额"`
	Cost     uint64 `json:"cost" description:"用户花费"`
	Income   uint64 `json:"income,omitempty" description:"主播收益"`
	Timelong uint64 `json:"timelong" description:"聊天时长"`
}

//WsTraficInfo .
type WsTraficInfo struct {
	NIMChannelID uint64 `json:"NIMChannelID" description:"网易云房间ID"`
	Timelong     uint64 `json:"income" description:"聊天时长"`
}

//Create .
// @Title 创建聊天通道[websocket]
// @Description 创建聊天通道[websocket]
// @Param   parterid     path    int  true        "聊天对象的ID"
// @router /:parterid/create [get]
func (c *WebsocketController) Create() {
	parter, err := strconv.ParseUint(c.Ctx.Input.Param(":parterid"), 10, 64)
	if err != nil {
		beego.Error(err)
		return
	}

	tk := GetToken(c.Ctx)
	_, ok := chattingUser[tk.ID]
	if ok {
		return
	}

	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		beego.Error(err)
		return
	}

	up := &models.UserProfile{ID: tk.ID}
	if err := up.Read(); err != nil {
		if err != nil {
			beego.Error(err)
			return
		}
	}

	channel := &ChatChannel{Join: make(chan *ChatClient), Send: make(chan *WsMessage), Exit: make(chan error, 1), ID: tk.ID, DstID: parter}
	go channel.Run()

	client := &ChatClient{User: up, Conn: conn, Send: make(chan *WsMessage)}
	client.Channel = channel
	chatChannels[tk.ID] = channel
	chattingUser[tk.ID] = nil

	channel.Join <- client //加入频道

	go client.Read()
	go client.Write()
}

//Join .
// @Title 加入聊天通道[websocket]
// @Description 加入聊天通道[websocket]
// @Param   channelid     path    int  true        "聊天频道的ID"
// @router /:channelid/join [get]
func (c *WebsocketController) Join() {
	channelid, err := strconv.ParseUint(c.Ctx.Input.Param(":channelid"), 10, 64)
	if err != nil {
		beego.Error(err)
		return
	}

	tk := GetToken(c.Ctx)
	cn, ok := chatChannels[channelid]
	if !ok || (cn.DstID != tk.ID && cn.ID != tk.ID) {
		return
	}

	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		beego.Error(err)
		return
	}

	var client *ChatClient
	if tk.ID == cn.ID {
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

			client = &ChatClient{User: up, Channel: cn, Conn: conn, Send: make(chan *WsMessage)}

			cn.Join <- client
		} else {
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

	cn.Exit <- nil
	dl := &models.Dial{FromUserID: channelid, ToUserID: tk.ID, Duration: 0, CreateAt: time.Now().Unix(), Status: models.DialStatusFail}
	if err := dl.Add(); err != nil {
		dto.Message = "websocket关闭成功，添加通话记录失败\t" + err.Error()
		return
	}

	(&models.UserProfile{ID: tk.ID}).AddDialReject(nil)

	dto.Sucess = true
}

func (c *ChatChannel) genVideoCost() (*VideoCost, error) {
	vc := &VideoCost{}
	up := (&models.UserProfile{ID: c.ID})
	if err := up.Read(); err != nil {
		return nil, err
	}

	vc.Balance = up.Balance
	vc.Cost = c.Amount
	return vc, nil
}

//Run .
func (c *ChatChannel) Run() {
	defer c.CloseChannel()

	for {
		select {
		case client := <-c.Join:
			if client.User.ID == c.ID {
				c.Src = client
			} else if client.User.ID == c.DstID {
				c.Dst = client
			}
		case msg := <-c.Send:
			switch msg.MessageType {
			case wsMessageTypeChannelStart:
				if !c.Inited { //双方进入房间，初始化房间，开始视频
					c.Inited = true
					c.StartTime = time.Now().Unix()

					//扣费
					if err := videoAllocateFund(c.Src.User, c.Dst.User, c.Price); err != nil {
						beego.Error(err)
						c.Exit <- nil
					}

					go func() {
						ticker := time.NewTicker(60 * time.Second)
						defer ticker.Stop()

						for {
							select {
							case <-ticker.C:
								if c.StopTime == 0 { //若频道未关闭,每隔60秒扣费一次
									c.Timelong = uint64(time.Now().Unix() - c.StartTime)

									//扣费
									if err := videoAllocateFund(c.Src.User, c.Dst.User, c.Price); err != nil { //扣费失败关闭聊天频道
										beego.Error(err)
										c.Exit <- nil
									}

									c.Amount += c.Price
								} else {
									break
								}
							}
						}
					}()
				}

				vc, err := c.genVideoCost()
				if err != nil {
					beego.Error("生成消费信息失败", err)
					c.Exit <- nil
				}

				m := &WsMessage{MessageType: wsMessageTypeChannelStart}
				vcStr, err := utils.JSONMarshalToString(vc)
				if err != nil {
					beego.Error("解析消费信息失败", err)
					c.Exit <- nil
				}

				m.Content = vcStr
				c.Src.Send <- m
			case wsMessageTypeClose:
				c.StopTime = time.Now().Unix()
				c.Timelong = uint64(c.StopTime - c.StartTime)

				m := &WsMessage{MessageType: wsMessageTypeChannelEnd}
				vc := &VideoCost{Cost: c.Amount, Timelong: c.Timelong}

				income, _, _ := computeIncome(c.Amount)
				vc.Income = uint64(income)

				vcStr, err := utils.JSONMarshalToString(vc)
				if err != nil {
					beego.Error("解析消费信息失败", err)
				}

				m.Content = vcStr

				c.Src.Send <- m
				c.Dst.Send <- m

				c.Exit <- nil
				return
			}
			beego.Info(utils.JSONMarshalToString(msg))
		case exp := <-c.Exit:
			if !c.Inited { //在主播未加入之前取消，不做任何操作
				return
			}

			if c.StopTime == 0 {
				c.StopTime = time.Now().Unix()

			}
			c.Timelong = uint64(c.StopTime - c.StartTime)

			dt := &models.DialTag{}

			//生成变动
			if err := videoDone(c.Src.User, c.Dst.User, &models.VideoChgInfo{TimeLong: c.Timelong, Price: c.Price}, c.Amount); err != nil {
				beego.Error("[websocket结算异常]视频结算生成余额变动错误", err, "发起人:", c.ID, "接受人:", c.DstID, "金额:", c.Amount, "通话时长:", c.Timelong)
				dt.ErrorMsg = append(dt.ErrorMsg, err.Error())
			}

			//生成通话记录
			dl := &models.Dial{FromUserID: c.ID, ToUserID: c.DstID, Duration: c.Timelong, CreateAt: c.StartTime, Status: models.DialStatusSuccess}
			if exp != nil {
				dl.Status = models.DialStatusException
				dt.ErrorMsg = append(dt.ErrorMsg, exp.Error())
				dl.Tag, _ = utils.JSONMarshalToString(dt)
			}

			if err := dl.Add(); err != nil {
				js, _ := utils.JSONMarshalToString(dl)
				beego.Error("[websocket结算异常]通话记录生成失败", err, js)
			}

			if err := c.Src.User.AddDialDuration(c.Timelong, nil); err != nil {
				beego.Error("[websocket结算异常]用户增加通话时长失败", err, c.Src.User.ID)
			}
			if err := c.Dst.User.AddDialDuration(c.Timelong, nil); err != nil {
				beego.Error("[websocket结算异常]用户增加通话时长失败", err, c.Dst.User.ID)
			}

			return
		}
	}
}

//CloseChannel 关闭频道,通道,websocket链接
func (c *ChatChannel) CloseChannel() {
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
	oldConnection := c.Conn

	defer func() {
		if c.Channel.StopTime == 0 { //聊天通道未结束
			// time.Sleep(30 * time.Second) //等待30秒,如过连接没有恢复,执行退出

			if oldConnection != c.Conn { //重连成功
				return
			}
		}

		c.Channel.Exit <- nil //发送退出信号,关闭通道后write方法会立即退出
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		if _, ok := chatChannels[c.Channel.ID]; ok {
			c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		}
		return nil
	})

	for {
		m := &WsMessage{}
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				beego.Error(err)
			}
			break
		}

		err = utils.JSONUnMarshalFromByte(message, m)
		if err == nil {
			c.Channel.Send <- m
		}
	}
}

func (c *ChatClient) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.关闭聊天室
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteJSON(message)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
