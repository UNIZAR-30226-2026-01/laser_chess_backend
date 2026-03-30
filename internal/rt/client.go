package rt

// Fichero que define la informacion que hay que guardar
// de un usuario en un hub

import (
	"time"

	"github.com/gorilla/websocket"
)

type ClientSocketMessage struct {
	Type    string `json:"Type"`
	Content string `json:"Content"`
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
	}

	return nil
}

// Cierra la conexion de un cliente
func (c *Client) Close() error {
	deadline := time.Now().Add(time.Minute)
	err := c.Conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		deadline,
	)
	if err != nil {
		return err
	}

	err = c.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return err
	}

	for {
		_, _, err = c.Conn.NextReader()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			break
		}
		if err != nil {
			break
		}
	}

	err = c.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}
