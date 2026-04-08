package sse

import "sync"

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

func (es *EventSystem) SendEvent(userID int64, event *Event) {
	es.mu.RLock()
	chSlice, exists := es.clientChannels[userID]
	es.mu.RUnlock()

	if !exists {
		es.fcm.SendNotification([]string{""}, event) // TODO: coger el token de dispositivo
		return
	}

	for _, ch := range chSlice {
		ch <- *event
	}
}

func (es *EventSystem) SaveChan(userID int64, eventChan chan Event) {
	es.mu.Lock()
	es.clientChannels[userID] = append(es.clientChannels[userID], eventChan)
	es.mu.Unlock()
}
