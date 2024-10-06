package accountPool

import (
	"api_wappi/interal/model"
	"sync"
)

// AccountPool предоставляет методы для работы с множеством аккаунтов
type AccountPool struct {
	accounts map[string]*model.Account
	mu       sync.RWMutex
}

func NewAccountPool() *AccountPool {
	return &AccountPool{
		accounts: make(map[string]*model.Account),
	}
}

func (ap *AccountPool) AddAccount(id string, account *model.Account) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	ap.accounts[id] = account
}

func (ap *AccountPool) RemoveAccount(id string) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	delete(ap.accounts, id)
}
