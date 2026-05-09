package auth

import (
	"sync"
	"time"
)

type TicketData struct {
	AccountID int64
	ExpiresAt time.Time
}

type WsTicketStore struct {
	mu      sync.Mutex
	tickets map[string]TicketData
}

func NewWsTicketStore() *WsTicketStore {
	store := &WsTicketStore{
		tickets: make(map[string]TicketData),
	}

	// Goroutine que limpia tickets expirados cada 5 segundos
	go func() {
		for {
			time.Sleep(5 * time.Second)
			store.cleanup()
		}
	}()

	return store
}

func (s *WsTicketStore) Store(ticket string, accountID int64, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tickets[ticket] = TicketData{
		AccountID: accountID,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Devuelve el accountID y un booleano indicando si era válido. Lo borra siempre.
func (s *WsTicketStore) GetAndDelete(ticket string) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.tickets[ticket]
	if !exists {
		return -1, false
	}

	// Borramos el ticket inmediatamente
	delete(s.tickets, ticket)

	// Comprobamos si ha caducado
	if time.Now().After(data.ExpiresAt) {
		return -1, false
	}

	return data.AccountID, true
}

func (s *WsTicketStore) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for ticket, data := range s.tickets {
		if now.After(data.ExpiresAt) {
			delete(s.tickets, ticket)
		}
	}
}
