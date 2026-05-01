package sse

import (
	"fmt"
	"net/http"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	md "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func (es *EventSystem) EventHandler(ctx *gin.Context) {
	// Cabeceras SSE
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	clientGone := ctx.Request.Context().Done()

	eventChan := make(chan Event, 10)

	userID, err := md.ExtractAccountID(ctx)
	if err != nil {
		apierror.DetectAndSendError(ctx, err)
		return
	}

	es.SaveChan(userID, eventChan)
	defer es.removeChan(userID, eventChan)

	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Evento inicial
	ctx.SSEvent("Init", "connected")
	flusher.Flush()

	fmt.Println("Cliente conectado:", userID)

	for {
		select {
		case <-clientGone:
			fmt.Println("Cliente desconectado:", userID)
			return

		case event := <-eventChan:
			ctx.SSEvent(event.EventType, event.Data)
			flusher.Flush()
		}
	}
}
