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
	onlineUsers    map[int64]bool
	fcm            *FirebaseManager
	mu             sync.RWMutex
}

func InitSSE(fcm *FirebaseManager) *EventSystem {
	return &EventSystem{
		clientChannels: make(map[int64][]chan Event),
		onlineUsers:    make(map[int64]bool),
		fcm:            fcm,
	}
}

func (es *EventSystem) MarkOnline(userID int64) {
	es.mu.Lock()
	es.onlineUsers[userID] = true
	es.mu.Unlock()

	fmt.Println("Usuario online:", userID)
}

func (es *EventSystem) MarkOffline(userID int64) {
	es.mu.Lock()
	delete(es.onlineUsers, userID)
	es.mu.Unlock()

	fmt.Println("Usuario offline:", userID)
}

func (es *EventSystem) IsOnline(userID int64) bool {
	es.mu.RLock()
	online := es.onlineUsers[userID]
	es.mu.RUnlock()

	return online
}

func (es *EventSystem) SendEvent(userID int64, event *Event, sendsFCM bool) {
	fmt.Println("SendEvent llamado:", userID, event.EventType, event.Data, "sendsFCM=", sendsFCM)
	
	if !es.IsOnline(userID) {
		fmt.Println("Usuario offline, enviando FCM: ", userID, event.EventType)

		if sendsFCM && es.fcm != nil {
			if err := es.fcm.SendNotification(userID, event); err != nil {
				fmt.Println("Error enviando FCM:", err)
			}
		}
		return
	}

	fmt.Println("Usuario online, enviando SSE:", userID, event.EventType)

	es.mu.RLock()
	chSlice := append([]chan Event(nil), es.clientChannels[userID]...)
	es.mu.RUnlock()

	if len(chSlice) == 0 {
		fmt.Println("Online pero sin canal SSE, fallback a FCM:", userID)

		if sendsFCM && es.fcm != nil {
			if err := es.fcm.SendNotification(userID, event); err != nil {
				fmt.Println("Error enviando FCM:", err)
			}
		}
		return
	}

	// Envío no bloqueante
	for _, ch := range chSlice {
		select {
		case ch <- *event:
			fmt.Println("Evento entregado por SSE:", userID, event.EventType)
		default:
			// canal bloqueado → lo ignoramos
			fmt.Println("Canal SSE bloqueado:", userID, event.EventType)
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
