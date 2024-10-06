package phone

import (
	"go.mau.fi/whatsmeow/types"
	"net/http"
)

// Service сервис для работы с телефонами
type Service interface {
	FindByPhoneNumbers(numbers []string) (map[types.JID]types.UserInfo, error)
}

type Handler struct {
	service Service
}

// NewHandler возвращает обработчик для работы с телефонами
func NewHandler(service Service) Handler {
	return Handler{
		service: service,
	}
}

func (h Handler) FindByNumbers(w http.ResponseWriter, r *http.Request) {

}
