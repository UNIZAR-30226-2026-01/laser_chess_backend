package rt

// Fichero que define la informacion que hay que guardar
// de un usuario en un hub

import (
	"time"

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

	// Canal para avisar de fin
	Done chan struct{}
}

func (c *Client) InitClient(AccountID int64, Conn *websocket.Conn) {
	c.AccountID = AccountID
	c.Conn = Conn
	c.Send = make(chan game.ResponseToRoom)
	c.ToRoom = make(chan ClientSocketMessage)

	c.Done = make(chan struct{})

	go c.ReadPump()
	go c.WritePump()
}

// lee mensajes del socket y los manda a la Room
func (c *Client) ReadPump() error {
	defer c.Close()

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
		close(c.Done)
		c.Close()
	}()

	for message := range c.Send {
		err := c.Conn.WriteJSON(message)
		if err != nil {
			return err
		}

		if message.Type == game.End || message.Type == game.Paused {

			return nil
		}
	}

	return nil
}

func (c *Client) Close() {
	deadline := time.Now().Add(time.Second)
	c.Conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		deadline,
	)
	c.Conn.Close()
}
