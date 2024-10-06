package model

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"time"
)

// Account аккаунт пользователя WhatsApp
type Account struct {
	// ID - выступает номер телефона
	ID         string
	Client     *whatsmeow.Client
	QRChannel  <-chan whatsmeow.QRChannelItem
	LastActive time.Time
	Device     *store.Device
}
