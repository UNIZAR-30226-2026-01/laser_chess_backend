package sse

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// Definicion del servicio

type FirebaseManager struct {
	App       *firebase.App
	Messaging *messaging.Client
}

func InitFirebase() (*FirebaseManager, error) {
	ctx := context.Background()

	firebasePath := os.Getenv("FIREBASE_CONFIG_PATH")
	opt := option.WithCredentialsFile(firebasePath)

	var err error
	var f *FirebaseManager = &FirebaseManager{}
	f.App, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	f.Messaging, err = f.App.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing messaging service: %v", err)
	}

	return f, nil
}

func (f *FirebaseManager) SendNotification(tokens []string,
	event *Event) error {

	ctx := context.Background()
	client, err := f.App.Messaging(ctx)
	if err != nil {
		//
	}

	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: event.EventType,
			Body:  "",
		},
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"extra_info":   "valor_personalizado",
		},
	}

	responses, err := client.SendEachForMulticast(ctx, message)
	if err != nil {
		//
	}

	if responses.FailureCount > 0 {
		// TODO: borrar los dispositivos que fallen
	}

	return nil
}
