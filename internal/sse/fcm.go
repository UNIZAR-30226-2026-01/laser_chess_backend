package sse

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/device"
	"google.golang.org/api/option"
)

// Definición del servicio

type FirebaseManager struct {
	App           *firebase.App
	Messaging     *messaging.Client
	deviceService *db.DeviceService
}

func InitFirebase(devices *db.DeviceService) (*FirebaseManager, error) {
	ctx := context.Background()

	firebasePath := os.Getenv("FIREBASE_CONFIG_PATH")
	if firebasePath == "" {
		return nil, fmt.Errorf("FIREBASE_CONFIG_PATH no definido")
	}

	credentialsJSON, err := os.ReadFile(firebasePath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo credenciales Firebase: %v", err)
	}

	opt := option.WithCredentialsJSON(credentialsJSON)

	var f *FirebaseManager = &FirebaseManager{}

	f.App, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	f.Messaging, err = f.App.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing messaging service: %v", err)
	}

	f.deviceService = devices
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
	tokens, err := f.deviceService.GetDevicesById(ctx, userID)
	if err != nil {
		fmt.Println("FCM error obteniendo tokens:", err)
		return err
	}

	fmt.Println("FCM tokens para user", userID, ":", tokens)

	if len(tokens) == 0 {
		fmt.Println("FCM: usuario sin tokens registrados:", userID)
		return nil
	}

	// Enviamos a todos los dispositivos la notificacion
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Data: map[string]string{
			"event_type": event.EventType,
			"data":       fmt.Sprintf("%v", event.Data),
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
	}

	responses, err := client.SendEachForMulticast(ctx, message)
	if err != nil {
		fmt.Println("FCM error enviando multicast:", err)
		return err
	}

	fmt.Println("FCM resultado: success=", responses.SuccessCount, " failure=", responses.FailureCount)

	// Filtramos los dispositivos que ya no estan registrados y los borramos
	for i, response := range responses.Responses {
		if !response.Success {
			fmt.Println("FCM error en token:", tokens[i], "error:", response.Error)

			if messaging.IsUnregistered(response.Error) {
				_, _ = f.deviceService.DeleteDevice(ctx, tokens[i])
			}
		} else {
			fmt.Println("FCM mensaje enviado correctamente al token:", tokens[i])
		}
	}

	return nil
}
