package phone

import (
	"api_wappi/interal/model"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// Service сервис для работы с телефонами
type Service interface {
	InfoByNumber(accountId int, phoneNumber string) (model.UserInfo, error)
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

func (h Handler) SearchByPhoneNumber(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения id: %v", err), http.StatusBadRequest)
		return
	}

	phone := r.URL.Query().Get("phoneNumber")
	if phone == "" {
		http.Error(w, "Нет номера телефона", http.StatusBadRequest)
		return
	}

	info, err := h.service.InfoByNumber(id, phone)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения информации по номеру %s - %v", phone, err), http.StatusInternalServerError)
		return
	}

	isWhatsapp := strconv.FormatBool(info.IsIn)
	response := map[string]string{
		"is_whatsapp": isWhatsapp,
		"status":      info.Status,
		"devices":     info.Devices[0].User,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка формирования JSON ответа: %v", err), http.StatusInternalServerError)
		return
	}
}
