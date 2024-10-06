package devicePool

import (
	"go.mau.fi/whatsmeow/store"
	"sync"
)

// DevicePool предоставляет методы для работы с множеством устройств
type DevicePool struct {
	devices map[string]*store.Device
	mu      sync.RWMutex
}

// NewDevicePool создает новый пул устройств
func NewDevicePool() *DevicePool {
	return &DevicePool{
		devices: make(map[string]*store.Device),
	}
}

// AddDevice добавляет устройство в пул
func (dp *DevicePool) AddDevice(id string, device *store.Device) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.devices[id] = device
}

// GetDevice получает устройство из пула по его ID
func (dp *DevicePool) GetDevice(id string) (*store.Device, bool) {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	device, exists := dp.devices[id]
	return device, exists
}

// RemoveDevice удаляет устройство из пула по его ID
func (dp *DevicePool) RemoveDevice(id string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	delete(dp.devices, id)
}

// ListDevices возвращает список всех устройств
func (dp *DevicePool) ListDevices() []*store.Device {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	devices := make([]*store.Device, 0, len(dp.devices))
	for _, device := range dp.devices {
		devices = append(devices, device)
	}
	return devices
}

func (dp *DevicePool) UpdateDeviceID(oldID, newID string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	if device, exists := dp.devices[oldID]; exists {
		delete(dp.devices, oldID)
		dp.devices[newID] = device
	}
}
