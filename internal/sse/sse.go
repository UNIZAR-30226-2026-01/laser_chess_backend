package sse

import (
	"fmt"
	"sync"
)

type Event struct {
	EventType string      `json:"event_type"`
	Data      interface{} `json:"data"`
}

type EventSystem struct {
	clientChannels map[int64][]chan Event
	fcm            *FirebaseManager
	mu             sync.RWMutex
}

func InitSSE(fcm *FirebaseManager) *EventSystem {
	return &EventSystem{
		clientChannels: make(map[int64][]chan Event),
		fcm:            fcm,
	}
}

func (es *EventSystem) SendEvent(userID int64, event *Event, sendsFCM bool) {
	es.mu.RLock()
	chSlice, exists := es.clientChannels[userID]
	es.mu.RUnlock()

	// Si no hay clientes conectados → FCM
	if !exists && sendsFCM && es.fcm != nil {
		err := es.fcm.SendNotification(userID, event)
		if err != nil {
			fmt.Println("Error enviando FCM:", err)
		}
		return
	}

	// Envío no bloqueante
	for _, ch := range chSlice {
		select {
		case ch <- *event:
		default:
			// canal bloqueado → lo ignoramos
			fmt.Println("Canal bloqueado, evento descartado")
		}
	}
}

func (es *EventSystem) SaveChan(userID int64, eventChan chan Event) {
	es.mu.Lock()
	es.clientChannels[userID] = append(es.clientChannels[userID], eventChan)
	es.mu.Unlock()
}

func (es *EventSystem) removeChan(userID int64, ch chan Event) {
	es.mu.Lock()
	defer es.mu.Unlock()

	channels := es.clientChannels[userID]

	for i, c := range channels {
		if c == ch {
			es.clientChannels[userID] = append(channels[:i], channels[i+1:]...)
			break
		}
	}

	// limpiar si no quedan canales
	if len(es.clientChannels[userID]) == 0 {
		delete(es.clientChannels, userID)
	}
}
