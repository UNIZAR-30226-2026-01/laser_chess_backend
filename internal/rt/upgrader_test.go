package rt

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// newWSPair lanza un httptest.Server que llama a UpgradeConn y devuelve
// la conexión del lado cliente (la que usaría el navegador).
// El servidor queda bloqueado hasta que el test cierra clientConn.
func newWSPair(t *testing.T) (clientConn *websocket.Conn, cleanup func()) {
	t.Helper()

	ready := make(chan *websocket.Conn, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := UpgradeConn(w, r)
		if err != nil {
			t.Errorf("UpgradeConn falló en el servidor: %v", err)
			return
		}
		ready <- conn

		// Mantener la conexión abierta hasta que el cliente la cierre
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))

	// Convertir la URL http → ws
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		srv.Close()
		t.Fatalf("no se pudo conectar al servidor de test: %v", err)
	}

	// Esperar a que el servidor haya completado el upgrade
	<-ready

	cleanup = func() {
		clientConn.Close()
		srv.Close()
	}
	return clientConn, cleanup
}

// ---------------------------------------------------------------------------
// UpgradeConn
// ---------------------------------------------------------------------------

// TestUpgradeConn_OK comprueba que una petición HTTP normal se convierte
// correctamente en una conexión WebSocket funcional.
func TestUpgradeConn_OK(t *testing.T) {
	clientConn, cleanup := newWSPair(t)
	defer cleanup()

	if clientConn == nil {
		t.Fatal("la conexión WebSocket del cliente es nil")
	}
}

// TestUpgradeConn_CanSendAndReceive comprueba que la conexión resultante
// permite intercambiar mensajes en ambas direcciones.
func TestUpgradeConn_CanSendAndReceive(t *testing.T) {
	ready := make(chan *websocket.Conn, 1)
	msgs := make(chan string, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := UpgradeConn(w, r)
		if err != nil {
			t.Errorf("UpgradeConn falló: %v", err)
			return
		}
		ready <- conn

		// Leer un mensaje y devolverlo (echo)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Errorf("error leyendo mensaje: %v", err)
			return
		}
		msgs <- string(msg)
		conn.WriteMessage(websocket.TextMessage, msg)
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("no se pudo conectar: %v", err)
	}
	defer clientConn.Close()
	<-ready

	// Enviar mensaje desde el cliente
	payload := "hola"
	if err := clientConn.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
		t.Fatalf("error enviando mensaje: %v", err)
	}

	// Verificar que el servidor lo recibió
	received := <-msgs
	if received != payload {
		t.Errorf("el servidor recibió %q, se esperaba %q", received, payload)
	}

	// Verificar que el cliente recibe el echo
	_, echo, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("error leyendo echo: %v", err)
	}
	if string(echo) != payload {
		t.Errorf("el cliente recibió echo %q, se esperaba %q", string(echo), payload)
	}
}
