package rt

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// TestClient_ReadPump verifica que cuando el navegador envía un JSON,
// el objeto Client lo lee y lo envía correctamente al canal ToRoom.
func TestClient_ReadPump(t *testing.T) {
	// 1. Crear servidor de prueba para obtener una conexión WS
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := UpgradeConn(w, r)
		client := &Client{}
		// Inicializamos el cliente. El AccountID es 1 y no es IA.
		client.InitClient(1, conn, false)

		// En este test, nos quedamos esperando a ver qué llega a ToRoom
		msgReceived := <-client.ToRoom
		assert.Equal(t, "Move", msgReceived.Type)
		assert.Equal(t, "Tc1:c2", msgReceived.Content)
	}))
	defer srv.Close()

	// 2. Conectar cliente (simulando el navegador)
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	ws, _, _ := websocket.DefaultDialer.Dial(wsURL, nil)
	defer ws.Close()

	// 3. Enviar mensaje JSON desde el "navegador"
	msg := ClientSocketMessage{Type: "Move", Content: "Tc1:c2"}
	err := ws.WriteJSON(msg)
	assert.NoError(t, err)

	// Pequeña espera para que las gorutinas procesen
	time.Sleep(100 * time.Millisecond)
}

// TestClient_WritePump verifica que cuando la Room envía un mensaje al canal Send,
// el objeto Client lo escribe automáticamente en el WebSocket para el navegador.
func TestClient_WritePump(t *testing.T) {
	done := make(chan bool)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := UpgradeConn(w, r)
		client := &Client{}
		client.InitClient(1, conn, false)

		// Simulamos que la sala envía un mensaje al cliente
		client.Send <- game.ResponseToRoom{
			Type:    game.Move,
			Content: "Tc1:c2%laserpath%{30000}",
		}

		// Esperamos a que el test termine
		<-done
	}))
	defer srv.Close()

	// 2. Conectar "navegador"
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	ws, _, _ := websocket.DefaultDialer.Dial(wsURL, nil)
	defer ws.Close()

	// 3. Leer mensaje del WebSocket
	_, p, err := ws.ReadMessage()
	assert.NoError(t, err)

	var resp game.ResponseToRoom
	err = json.Unmarshal(p, &resp)
	assert.NoError(t, err)
	assert.Equal(t, game.Move, resp.Type)
	assert.Contains(t, resp.Content, "Tc1:c2")

	close(done)
}

// TestClient_Disconnection verifica que al cerrar la conexión,
// el cliente notifica a la sala con un mensaje EOC (End of Connection).
func TestClient_Disconnection(t *testing.T) {
	msgChan := make(chan ClientSocketMessage, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := UpgradeConn(w, r)
		client := &Client{}
		client.InitClient(1, conn, false)

		// Capturamos el mensaje de desconexión
		go func() {
			for msg := range client.ToRoom {
				if msg.Type == string(game.EOC) {
					msgChan <- msg
				}
			}
		}()

		// Esperamos a que la conexión se cierre fuera
		<-client.Done
	}))
	defer srv.Close()

	// 2. Conectar y cerrar inmediatamente
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	ws, _, _ := websocket.DefaultDialer.Dial(wsURL, nil)
	ws.Close()

	// 3. Verificar que llegó el mensaje de desconexión (EOC)
	select {
	case msg := <-msgChan:
		assert.Equal(t, string(game.EOC), msg.Type)
	case <-time.After(1 * time.Second):
		t.Fatal("No se recibió notificación de desconexión")
	}
}
