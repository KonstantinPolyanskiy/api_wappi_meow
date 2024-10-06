package serivce

import (
	"api_wappi/interal/serivce/whatsapp"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type WhatsAppService interface {
	InitiateAuthentication(accountId int) (string, error)
}

// AggregatorService собирает в себе все сервисы и прочие зависимости для их работы
type AggregatorService struct {
	WhatsAppService
}

// New Инстанцирует слой сервисов
func New(container *sqlstore.Container) AggregatorService {
	return AggregatorService{
		WhatsAppService: whatsapp.New(container),
	}
}
