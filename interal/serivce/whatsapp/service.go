package whatsapp

import (
	"api_wappi/interal/accountPool"
	"api_wappi/interal/clientPool"
	"api_wappi/interal/devicePool"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type Service struct {
	accPool    *accountPool.AccountPool
	devicePool *devicePool.DevicePool
	clientPool *clientPool.ClientPool

	container *sqlstore.Container
}

// New Инстанцирует сервис по работе с WhatsApp.
// Принимает:
// Клиент для работы с API WhatsApp.
// Пул аккаунтов WhatsApp.
func New(container *sqlstore.Container) Service {
	cp := clientPool.NewClientPool()
	ap := accountPool.NewAccountPool()
	dp := devicePool.NewDevicePool()

	return Service{
		accPool:    ap,
		devicePool: dp,
		clientPool: cp,
		container:  container,
	}
}
