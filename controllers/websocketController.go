package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

const (
	wsMessageTypeText = iota
	wsMessageTypeGift
	wsMessageTypeSystem
	wsMessageTypeClose
)

var (
	upgrader     = websocket.Upgrader{}
	chatChannels = make(map[uint64]*ChatChannel)
)

//WebsocketController websocket相关
type WebsocketController struct {
	beego.Controller
}

//ChatChannel .
type ChatChannel struct {
	ID        uint64
	DstID     uint64
	StartTime int64
	Timelong  uint64 //聊天时长,单位:秒
	StopTime  int64
	Price     uint64
	Src       *ChatClient
	Dst       *ChatClient
	Send      chan *WsMessage
	Join      chan *ChatClient
	Exit      chan *ChatClient
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
	FromID      uint64
	Content     string
	MessageType int
}

//Create .
// @Title 创建聊天通道
// @Description 创建聊天通道
// @Param   parterid     path    int  true        "聊天对象的ID"
// @router /:parterid/create [get]
func (c *WebsocketController) Create() {
	parter, err := strconv.ParseUint(c.Ctx.Input.Param(":parterid"), 10, 64)
	if err != nil {
		beego.Error(err)
		return
	}

	tk := GetToken(c.Ctx)
	_, ok := chatChannels[tk.ID]
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

	channel := &ChatChannel{Join: make(chan *ChatClient), Send: make(chan *WsMessage), Exit: make(chan *ChatClient, 1), ID: tk.ID, DstID: parter}
	go channel.Run()

	client := &ChatClient{User: up, Conn: conn, Send: make(chan *WsMessage)}
	client.Channel = channel
	chatChannels[tk.ID] = channel

	channel.Join <- client //加入频道

	go client.Read()
	go client.Write()
}

//Join .
// @Title 加入聊天通道
// @Description 加入聊天通道
// @Param   channelid     path    int  true        "聊天对象的ID"
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

	go client.Read()
	go client.Write()
}

//Reject .
// @Title 拒绝聊天请求
// @Description 拒绝聊天请求
// @Param   channelid     path    int  true        "聊天对象的ID"
// @router /:channelid/reject [post]
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

	dto.Sucess = true
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

				if c.StartTime == 0 { //开始计时
					c.StartTime = time.Now().Unix()

					ticker := time.NewTicker(60 * time.Second)
					//扣费
					if err := videoAllocateFund(c.Src.User, c.Dst.User, c.Price); err != nil {
						beego.Error(err)
						c.Exit <- nil
					}

					go func() {
						defer ticker.Stop()

						for {
							if c.StopTime == 0 { //若频道未关闭,每隔60秒扣费一次
								<-ticker.C
								c.Timelong += 60
								//扣费
								if err := videoAllocateFund(c.Src.User, c.Dst.User, c.Price); err != nil { //扣费失败关闭聊天频道
									beego.Error(err)
									c.Exit <- nil
								}
							}
						}
					}()
				}
			}
		case msg := <-c.Send:
			beego.Info(utils.JSONMarshalToString(msg))

			if msg.MessageType == wsMessageTypeClose {
				c.Exit <- nil
				return
			}

			if msg.FromID == c.DstID && c.Src != nil {
				c.Src.Send <- msg
			} else if c.Dst != nil {
				c.Dst.Send <- msg
			}
		case <-c.Exit:
			c.StopTime = time.Now().Unix()

			//生成变动
			if err := videoDone(c.Src.User, c.Dst.User, &models.VideoChgInfo{TimeLong: c.Timelong, Price: c.Price}); err != nil {
				beego.Error(errors.New("视频结算生成余额变动错误\t" + err.Error()))
			}

			//生成通话记录
			dl := &models.Dial{FromUserID: c.ID, ToUserID: c.DstID, Duration: int(c.Timelong), CreateAt: c.StartTime, Status: models.DialStatusSuccess}
			if err := dl.Add(); err != nil {
				beego.Error(errors.New("通话记录生成失败\t" + err.Error()))
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

	if c.Src != nil {
		close(c.Src.Send)
	}

	if c.Dst != nil {
		close(c.Dst.Send)
	}
}

func (c *ChatClient) Read() {
	defer func() {
		if c.Channel.StopTime == 0 { //聊天通道未结束
			oldConnection := c.Conn
			time.Sleep(30 * time.Second) //等待30秒,如过连接没有恢复,执行退出

			if oldConnection != c.Conn { //重连成功
				return
			}
		}

		c.Channel.Exit <- c //发送退出信号,关闭通道后write方法会立即退出
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
		m := &WsMessage{FromID: c.User.ID}
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
