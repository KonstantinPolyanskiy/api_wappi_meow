package account

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// Service сервис для работы с аккаунтами
type Service interface {
	InitiateAuthentication(accountId int) (string, error)
}

type Handler struct {
	service Service
}

// NewHandler возвращает обработчик для работы с аккаунтами
func NewHandler(service Service) Handler {
	return Handler{
		service: service,
	}
}

func (h Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения id: %v", err), http.StatusBadRequest)
	}

	qrCodeBase64, err := h.service.InitiateAuthentication(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка авторизацииы: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"id":      idStr,
		"qr_code": qrCodeBase64,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка формирования JSON ответа: %v", err), http.StatusInternalServerError)
		return
	}
}
