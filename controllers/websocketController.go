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
	ID    uint64
	DstID uint64
	Src   *ChatClient
	Dst   *ChatClient
	Send  chan *WsMessage
	Close chan *ChatClient
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

	client := &ChatClient{User: up, Conn: conn, Send: make(chan *WsMessage)}
	channel := &ChatChannel{Src: client, Send: make(chan *WsMessage), Close: make(chan *ChatClient), ID: tk.ID, DstID: parter}

	client.Channel = channel
	chatChannels[tk.ID] = channel
	go channel.Run()
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
	if !ok || cn.DstID != tk.ID || cn.Dst != nil {
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

	client := &ChatClient{User: up, Channel: cn, Conn: conn, Send: make(chan *WsMessage)}
	cn.Dst = client

	go client.Read()
	go client.Write()
}

//Run .
func (c *ChatChannel) Run() {
	defer func() {
		close(c.Send)
		close(c.Close)
		delete(chatChannels, c.ID)

		if c.Src != nil {
			close(c.Src.Send)
		}

		if c.Dst != nil {
			close(c.Dst.Send)
		}
	}()

	for {
		select {
		case msg := <-c.Send:
			beego.Info(utils.JSONMarshalToString(msg))

			if msg.FromID == c.DstID {
				c.Src.Send <- msg
			} else if c.Dst != nil {
				c.Dst.Send <- msg
			}
		case <-c.Close:
			return
		}
	}
}

func (c *ChatClient) Read() {
	defer func() {
		c.Conn.Close()
		c.Channel.Close <- c
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
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
				// The hub closed the channel.
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
