package model

import "go.mau.fi/whatsmeow/types"

type UserInfo struct {
	types.IsOnWhatsAppResponse
	types.UserInfo
}
