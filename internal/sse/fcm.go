package sse

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"google.golang.org/api/option"
)

// Definicion del servicio

type FirebaseManager struct {
	App            *firebase.App
	Messaging      *messaging.Client
	accountService *db.AccountService
}

func InitFirebase(accounts *db.AccountService) (*FirebaseManager, error) {
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

func (f *FirebaseManager) SendNotification(userID int64,
	event *Event) error {

	// Obtenemos el cliente
	ctx := context.Background()
	client, err := f.App.Messaging(ctx)
	if err != nil {
		return err
	}

	// Obtenemos los tokens de los dispositivos del cliente
	tokens, err := f.accountService.GetDevicesById(ctx, userID)

	// Enviamos a todos los dispositivos la notificacion
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
		return err
	}

	// Filtramos los dispositivos que ya no estan registrados y los borramos
	if responses.FailureCount > 0 {
		for i, response := range responses.Responses {
			if !response.Success {
				if messaging.IsUnregistered(response.Error) {
					f.accountService.DeleteDevice(ctx, tokens[i])
				}
			}
		}
	}

	return nil
}
