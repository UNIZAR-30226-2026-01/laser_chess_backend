package sse

import (
	"fmt"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	md "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func sseHandler(ctx *gin.Context, es *EventSystem) {
	// Cabeceras obligatorias para el SSE
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Canal para desconexion del cliente
	clientGone := ctx.Request.Context().Done()

	eventChan := make(chan Event)
	userID, err := md.ExtractAccountID(ctx)
	if err != nil {
		apierror.DetectAndSendError(ctx, err)
	}

	es.SaveChan(userID, eventChan)

	for {
		select {
		case <-clientGone:
			fmt.Println("Client disconnected")
			return
		case event := <-eventChan:
			ctx.SSEvent(event.EventType, event.Data)
		}
	}
}
