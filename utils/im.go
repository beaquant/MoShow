package utils

import (
	netease "github.com/MrSong0607/netease-im"
)

var imClient = netease.CreateImClient("36bb3190572f691d3b180fc099a1b4f1", "99c220190258", "")

//ImCreateUser .
func ImCreateUser(user *netease.ImUser) (*netease.TokenInfo, error) {
	return imClient.CreateImUser(user)
}

//ImRefreshToken .
func ImRefreshToken(id string) (*netease.TokenInfo, error) {
	return imClient.RefreshToken(id)
}
