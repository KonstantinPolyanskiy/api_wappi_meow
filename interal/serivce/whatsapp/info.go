package whatsapp

import (
	"api_wappi/interal/model"
	"fmt"
	"go.mau.fi/whatsmeow/types"
	"strconv"
)

func (s Service) InfoByNumber(accountId int, phoneNumber string) (model.UserInfo, error) {
	id := strconv.Itoa(accountId)

	client, exists := s.clientPool.GetClient(id)
	if !exists {
		return model.UserInfo{}, fmt.Errorf("клиент не существует")
	}

	isOnResp, err := client.IsOnWhatsApp([]string{phoneNumber})
	if err != nil {
		return model.UserInfo{}, err
	}

	info, err := client.GetUserInfo([]types.JID{isOnResp[0].JID})
	if err != nil {
		return model.UserInfo{}, err
	}

	result, exists := info[isOnResp[0].JID]
	if !exists {
		return model.UserInfo{}, fmt.Errorf("информация по номеру телефона не найдена")
	}

	return model.UserInfo{IsOnWhatsAppResponse: isOnResp[0], UserInfo: result}, nil

}
