package whatsapp

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
)

func (s Service) getOrCreateDeviceAndClient(id string) (*store.Device, *whatsmeow.Client) {
	var client *whatsmeow.Client
	var device *store.Device

	d, exists := s.devicePool.GetDevice(id)
	if exists {
		device = d
	} else {
		device = s.container.NewDevice()
		s.devicePool.AddDevice(id, device)
	}

	c, exists := s.clientPool.GetClient(id)
	if exists {
		client = c
	} else {
		client = whatsmeow.NewClient(device, nil)
		s.clientPool.AddClient(id, client)
	}

	return device, client
}
