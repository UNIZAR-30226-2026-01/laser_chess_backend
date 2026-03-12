package rt

// Fichero que define la informacion que hay que guardar
// de un usuario en un hub

import "github.com/gorilla/websocket"

type Client struct {
	AccountID int64
	Conn      *websocket.Conn
	Send      chan interface{} // canal para mandar mensajes al front
	Room      *Room
}

func (c *Client) InitClient(AccountID int64, Conn *websocket.Conn, 
			 Room *Room) {
	c.AccountID = AccountID
	c.Conn = Conn
	c.Send = make(chan interface{})
	c.Room = Room
}

// lee mensajes del socket y los manda a la Room
func (c *Client) ReadPump() {
	// aqui ira un bucle for que escucha c.Conn.ReadJSON()
	// y escribe en algun canal de Room
	// for {
	// 	v, err := c.Conn.ReadJSON()

		
	// }
}

// saca mensajes del canal c.Send y los escribe en el navegador
func (c *Client) WritePump() {
	// aquí ira un bucle for que escucha c.Send y hace c.Conn.WriteJSON()
	// for {

	// }
}

// cierra la conexion de un cliente
func (c *Client) Close() {
}

//TODO: rellenar las funciones
