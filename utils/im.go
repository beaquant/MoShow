package utils

import (
	"net/url"
	"path/filepath"

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

//SendP2PSysImageMessage 发送图片系统消息
func SendP2PSysImageMessage(URL string, toID []string) (string, error) {
	return SendP2PImageMessage(URL, ImSysAdminID, toID)
}

//SendP2PSysVoiceMessage duration单位:毫秒
func SendP2PSysVoiceMessage(URL string, duration uint, toID []string) (string, error) {
	return SendP2PVoiceMessage(URL, duration, ImSysAdminID, toID)
}

//SendP2PSysVideoMessage .
func SendP2PSysVideoMessage(URL string, duration uint, toID []string) (string, error) {
	return SendP2PVideoMessage(URL, duration, ImSysAdminID, toID)
}

//SendP2PImageMessage .
func SendP2PImageMessage(URL string, fromID string, toID []string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	return ImClient.SendBatchImageMessage(fromID, toID, &netease.ImageMessage{URL: URL, Md5: ShaHashToHexStringFromString(URL), Extension: filepath.Ext(u.Path)}, nil)
}

//SendP2PVoiceMessage duration单位:毫秒
func SendP2PVoiceMessage(URL string, duration uint, fromID string, toID []string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	return ImClient.SendBatchVoiceMessage(fromID, toID, &netease.VoiceMessage{URL: URL, Md5: ShaHashToHexStringFromString(URL), Duration: duration, Extension: filepath.Ext(u.Path)}, nil)
}

//SendP2PVideoMessage .
func SendP2PVideoMessage(URL string, duration uint, fromID string, toID []string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	return ImClient.SendBatchVideoMessage(fromID, toID, &netease.VideoMessage{URL: URL, Duration: duration, Md5: ShaHashToHexStringFromString(URL), Extension: filepath.Ext(u.Path)}, nil)
}

//SendDIYSysMessage 发送自定义系统消息
func SendDIYSysMessage(content *ImSysNotifyMessage, toIds []string) error {
	msgStr, err := JSONMarshalToString(content)
	if err != nil {
		return err
	}
	return ImClient.SendBatchAttachMsg(ImSysAdminID, msgStr, toIds, nil)
}
