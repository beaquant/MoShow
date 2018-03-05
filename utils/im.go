package utils

import (
	netease "github.com/MrSong0607/netease-im"
)

func genImClient() *netease.ImClient {
	return netease.CreateImClient("36bb3190572f691d3b180fc099a1b4f1", "99c220190258", "")
}

//ImCreateUser .
func ImCreateUser(user *netease.ImUser) (*netease.TokenInfo, error) {
	return genImClient().CreateImUser(user)
}

//ImRefreshToken .
func ImRefreshToken(id string) (*netease.TokenInfo, error) {
	return genImClient().RefreshToken(id)
}
