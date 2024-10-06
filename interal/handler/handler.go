package handler

import (
	"api_wappi/interal/handler/account"
	"api_wappi/interal/handler/phone"
	"api_wappi/interal/serivce"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type AccountHandler interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

type PhoneHandler interface {
	SearchByPhoneNumber(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	AccountHandler
	PhoneHandler
}

func New(service serivce.AggregatorService) Handler {
	return Handler{
		AccountHandler: account.NewHandler(service.WhatsAppService),
		PhoneHandler:   phone.NewHandler(service.WhatsAppService),
	}
}

func (h Handler) Init() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Post("/whatsapp/{id}/authenticate", h.Authenticate)
		r.Get("/whatsapp/{id}/search", h.SearchByPhoneNumber)
	})

	return r
}
