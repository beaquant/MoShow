package utils

import (
	netease "github.com/MrSong0607/netease-im"
	"github.com/astaxie/beego"
)

var (
	//ImClient .
	ImClient = netease.CreateImClient(beego.AppConfig.String("neteaseAppKey"), beego.AppConfig.String("neteaseAppSecret"), "") //http://127.0.0.1:8889
	//ImSysAdminID 系统管理员ID
	ImSysAdminID = "1"
)

const (
	//ImSysNotifyMessageTypeText 文本通知
	ImSysNotifyMessageTypeText = iota
)

//ImSysNotifyMessage 系统消息通知
type ImSysNotifyMessage struct {
	Message string `json:"message"`
	Type    int    `json:"type"`
}

//ImCreateUser .
func ImCreateUser(user *netease.ImUser) (*netease.TokenInfo, error) {
	return ImClient.CreateImUser(user)
}

//ImRefreshToken .
func ImRefreshToken(id string) (*netease.TokenInfo, error) {
	return ImClient.RefreshToken(id)
}

//SendP2PMessage .
func SendP2PMessage(fromID, toID, content string) error {
	return ImClient.SendTextMessage(fromID, toID, &netease.TextMessage{Message: content}, nil)
}

//SendP2PSysMessage 发送点对点系统消息
func SendP2PSysMessage(content string, toID string) error {
	return SendP2PMessage(ImSysAdminID, toID, content)
}

//SendDIYSysMessage 发送自定义系统消息
func SendDIYSysMessage(content *ImSysNotifyMessage, toIds []string) error {
	msgStr, err := JSONMarshalToString(content)
	if err != nil {
		return err
	}
	return ImClient.SendBatchAttachMsg(ImSysAdminID, msgStr, toIds, nil)
}
