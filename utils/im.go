package utils

import (
	netease "github.com/MrSong0607/netease-im"
)

var (
	neteaseAppKey    = "b2c60dbed0ae2d3c48e6c85664836dc9"
	neteaseAppSecret = "1ed04f7d7085"
	imClient         = netease.CreateImClient(neteaseAppKey, neteaseAppSecret, "") //http://127.0.0.1:8889
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

func genImClient() *netease.ImClient {
	return imClient
}

//ImCreateUser .
func ImCreateUser(user *netease.ImUser) (*netease.TokenInfo, error) {
	return genImClient().CreateImUser(user)
}

//ImRefreshToken .
func ImRefreshToken(id string) (*netease.TokenInfo, error) {
	return genImClient().RefreshToken(id)
}

//SendP2PMessage .
func SendP2PMessage(fromID, toID, content string) error {
	return imClient.SendTextMessage(fromID, toID, &netease.TextMessage{Message: content}, nil)
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
	return imClient.SendBatchAttachMsg(ImSysAdminID, msgStr, toIds, nil)
}
