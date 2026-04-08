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

func (f *FirebaseManager) InitFirebase() error {
	ctx := context.Background()

	firebasePath := os.Getenv("FIREBASE_CONFIG_PATH")
	opt := option.WithCredentialsFile(firebasePath)

	var err error
	f.App, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	f.Messaging, err = f.App.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("error initializing messaging service: %v", err)
	}

	return nil
}

func (f *FirebaseManager) SendNotification(app *firebase.App, token string,
	message *messaging.Message) (string, error) {

	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		//
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		//
	}
	return response, nil
}

// Creacion del servicio y su interfaz

var F *FirebaseManager

func InitFirebase() error {
	return F.InitFirebase()
}

func SendChallenge(app *firebase.App, token string) (string, error) {
	message := &messaging.Message{
		Token: token, // El token que guardaste del usuario
		Notification: &messaging.Notification{
			Title: "¡Hola!",
			Body:  "Tienes un nuevo mensaje.",
		},
	}
	return F.SendNotification(app, token, message)
}
