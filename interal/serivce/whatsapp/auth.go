package whatsapp

import (
	"api_wappi/interal/model"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	"os"
	"strconv"
	"time"
)

func (s Service) InitiateAuthentication(accountId int) (string, error) {
	tempId := fmt.Sprintf("%d", time.Now().UnixNano())

	deviceStore := s.container.NewDevice()

	client := whatsmeow.NewClient(deviceStore, nil)

	qrChan, err := client.GetQRChannel(context.Background())
	if err != nil {
		return "", err
	}

	account := &model.Account{
		ID:         strconv.Itoa(accountId),
		Client:     client,
		QRChannel:  qrChan,
		Device:     deviceStore,
		LastActive: time.Now(),
	}

	s.accPool.AddAccount(tempId, account)

	err = client.Connect()
	if err != nil {
		s.devicePool.RemoveDevice(tempId)
		s.clientPool.RemoveClient(tempId)
		return "", fmt.Errorf("не удалось подключиться к клиенту: %v", err)
	}

	var qrCode string

	timeout := time.After(30 * time.Second)

QRLoop:
	for {
		select {
		case evt := <-qrChan:
			if evt.Event == "code" {
				qrCode = evt.Code
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				break QRLoop
			} else {
				fmt.Printf("Событие QR - : %s\n", evt.Event)
			}
		case <-timeout:
			// Таймаут ожидания QR-кода
			s.devicePool.RemoveDevice(tempId)
			s.clientPool.RemoveClient(tempId)
			return "", fmt.Errorf("конец ожидания QR кода")
		}
	}

	pngData, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		s.devicePool.RemoveDevice(tempId)
		s.clientPool.RemoveClient(tempId)
		return "", fmt.Errorf("не удалось сгенерировать QR: %v", err)
	}

	base64QRCode := base64.StdEncoding.EncodeToString(pngData)

	s.devicePool.AddDevice(tempId, deviceStore)
	s.clientPool.AddClient(tempId, client)

	client.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.Connected:
			if client.Store.ID != nil {
				accountID := client.Store.ID.User

				s.devicePool.UpdateDeviceID(tempId, accountID)
				s.clientPool.UpdateClientID(tempId, accountID)
				s.accPool.AddAccount(strconv.Itoa(accountId), account)

				fmt.Printf("Пользовать %s\n авторизирован", accountID)
			}
		case *events.Disconnected:
			s.devicePool.RemoveDevice(tempId)
			s.clientPool.RemoveClient(tempId)
			s.accPool.RemoveAccount(strconv.Itoa(accountId))
		}
	})

	return base64QRCode, nil
}
