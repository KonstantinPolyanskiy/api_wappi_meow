package whatsapp

import (
	"api_wappi/interal/model"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow/types/events"
	"os"
	"strconv"
	"time"
)

func (s Service) InitiateAuthentication(accountId int) (string, error) {
	var account model.Account
	id := strconv.Itoa(accountId)

	device, client := s.getOrCreateDeviceAndClient(strconv.Itoa(accountId))

	qrChan, err := client.GetQRChannel(context.Background())
	if err != nil {
		return "", err
	}

	a, exist := s.accPool.GetAccount(strconv.Itoa(accountId))
	if exist {
		account = *a
	} else {
		account = model.Account{
			ID:         strconv.Itoa(accountId),
			Client:     client,
			QRChannel:  qrChan,
			Device:     device,
			LastActive: time.Now(),
		}
	}

	err = client.Connect()
	if err != nil {
		s.devicePool.RemoveDevice(id)
		s.clientPool.RemoveClient(id)
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
			s.devicePool.RemoveDevice(id)
			s.clientPool.RemoveClient(id)
			return "", fmt.Errorf("конец ожидания QR кода")
		}
	}

	pngData, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		s.devicePool.RemoveDevice(id)
		s.clientPool.RemoveClient(id)
		return "", fmt.Errorf("не удалось сгенерировать QR: %v", err)
	}

	base64QRCode := base64.StdEncoding.EncodeToString(pngData)

	client.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.Connected:
			if client.Store.ID != nil {
				fmt.Printf("Пользовать %s\n авторизирован", account.Client.Store.ID)
			}
		case *events.Disconnected:
			s.devicePool.RemoveDevice(id)
			s.clientPool.RemoveClient(id)
			s.accPool.RemoveAccount(id)
		}
	})

	return base64QRCode, nil
}
