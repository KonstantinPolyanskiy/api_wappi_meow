package handler

import (
	"api_wappi/interal/handler/account"
	"api_wappi/interal/serivce"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type AccountHandler interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	AccountHandler
}

func New(service serivce.AggregatorService) Handler {
	return Handler{
		AccountHandler: account.NewHandler(service.WhatsAppService),
	}
}

func (h Handler) Init() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/api/whatsapp/{id}/authenticate", h.Authenticate)

	return r
}
