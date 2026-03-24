package rt

// Fichero que define la informacion que hay que guardar
// de un usuario en un hub

import (
	"github.com/gorilla/websocket"
)

type ClientSocketMessage struct {
	Type    string `json:"type" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type Client struct {
	AccountID int64
	Conn      *websocket.Conn
	Send      chan interface{}
	ToRoom    chan ClientSocketMessage

	// Canal para avisar de fin
	Done chan struct{}
}

func (c *Client) InitClient(AccountID int64, Conn *websocket.Conn) {
	c.AccountID = AccountID
	c.Conn = Conn
	c.Send = make(chan interface{})
	c.ToRoom = make(chan ClientSocketMessage)

	c.Done = make(chan struct{})

	go c.ReadPump()
	go c.WritePump()
}

// lee mensajes del socket y los manda a la Room
func (c *Client) ReadPump() error {
	defer func() {
		close(c.Done)
		c.Conn.Close()
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
	defer c.Conn.Close()

	for message := range c.Send {
		err := c.Conn.WriteJSON(message)
		if err != nil {
			return err
		}
	}

	return nil
}

// Cierra la conexion de un cliente
func (c *Client) Close() error {
	return c.Conn.Close()
}
