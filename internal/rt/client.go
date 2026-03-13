package rt

// Fichero que define la informacion que hay que guardar
// de un usuario en un hub

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	AccountID int64
	Conn      *websocket.Conn
	Send      chan interface{} // canal para mandar mensajes al front
	Room      *Room
	ToRoom    chan interface{}
}

type ClientSocketMessage struct {
	Type    string
	Content string
}

func (c *Client) InitClient(AccountID int64, Conn *websocket.Conn,
	Room *Room, ToRoom chan interface{}) {
	c.AccountID = AccountID
	c.Conn = Conn
	c.Send = make(chan interface{})
	c.Room = Room

	go c.ReadPump()
	go c.WritePump()
}

// lee mensajes del socket y los manda a la Room
func (c *Client) ReadPump() error {
	// aqui ira un bucle for que escucha c.Conn.ReadJSON()
	// y escribe en algun canal de Room
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
	for {
		select {
		case message := <-c.Send:
			c.Conn.WriteJSON(message)
		}
	}
}

// cierra la conexion de un cliente
func (c *Client) Close() {
	c.Conn.Close()
}

//TODO: rellenar las funciones
