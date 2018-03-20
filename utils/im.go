package utils

import (
	netease "github.com/MrSong0607/netease-im"
)

func genImClient() *netease.ImClient {
	return netease.CreateImClient("b2c60dbed0ae2d3c48e6c85664836dc9", "1ed04f7d7085", "")
}

//ImCreateUser .
func ImCreateUser(user *netease.ImUser) (*netease.TokenInfo, error) {
	return genImClient().CreateImUser(user)
}

//ImRefreshToken .
func ImRefreshToken(id string) (*netease.TokenInfo, error) {
	return genImClient().RefreshToken(id)
}
