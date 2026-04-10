package rt

// Fichero que define la informacion que hay que guardar
// de un usuario en un hub

import (
	"fmt"
	"sync"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
	"github.com/gorilla/websocket"
)

type ClientSocketMessage struct {
	Type    string `json:"Type"`
	Content string `json:"Content"`
}

type Client struct {
	AccountID int64
	Conn      *websocket.Conn
	Send      chan game.ResponseToRoom
	ToRoom    chan ClientSocketMessage

	Reconnect chan bool
	Online    bool

	mu sync.RWMutex

	// Canal para avisar de fin
	Done chan struct{}
}

func (c *Client) InitClient(AccountID int64, Conn *websocket.Conn) {
	c.AccountID = AccountID
	c.Conn = Conn
	c.Send = make(chan game.ResponseToRoom, 1)
	c.ToRoom = make(chan ClientSocketMessage, 1)

	c.Done = make(chan struct{})
	c.Reconnect = make(chan bool, 1)

	c.mu.Lock()
	c.Online = true
	c.mu.Unlock()

	go c.ReadPump()
	go c.WritePump()
}

// lee mensajes del socket y los manda a la Room
func (c *Client) ReadPump() error {
	defer func() {
		select {
		case <-c.Done:
		default:
			close(c.Done)
		}
		c.Conn.Close()
		c.notifyDisconnection()
		fmt.Println("Cierre de la funcion ReadPump")
	}()

	for {
		var message ClientSocketMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			return err
		}

		c.ToRoom <- message
	}
}

// saca mensajes del canal c.Send y los escribe en el navegador
func (c *Client) WritePump() error {
	defer func() {
		select {
		case <-c.Done:
		default:
			close(c.Done)
		}
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				return nil
			}
			err := c.Conn.WriteJSON(message)
			if err != nil {
				return err
			}
			if message.Type == game.EOC {
				return nil
			}
		case <-c.Done:
			// Si ReadPump detecta un error salimos
			return nil
		}
	}
}

func (c *Client) notifyDisconnection() {
	fmt.Println("Desconexion identificadad desde client")
	c.mu.Lock()
	c.Online = false
	c.mu.Unlock()
	c.ToRoom <- ClientSocketMessage{Type: string(game.EOC), Content: ""}
}
