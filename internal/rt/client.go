package rt

// Fichero que define la informacion que hay que guardar
// de un usuario en un hub

import (
	"fmt"
	"strconv"
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

	isAI              bool
	initStateReceived bool
	ToAI              chan ClientSocketMessage
	FromAI            chan ClientSocketMessage

	mu sync.RWMutex

	// Canal para avisar de fin
	Done chan struct{}
}

func (c *Client) InitClient(AccountID int64, Conn *websocket.Conn, isAI bool) {
	c.AccountID = AccountID
	c.Conn = Conn
	c.Send = make(chan game.ResponseToRoom, 1)
	c.ToRoom = make(chan ClientSocketMessage, 1)

	c.Done = make(chan struct{})
	c.Reconnect = make(chan bool, 1)

	c.mu.Lock()
	c.Online = true
	c.mu.Unlock()

	if isAI {
		c.initStateReceived = false
		c.ToAI = make(chan ClientSocketMessage)
		c.FromAI = make(chan ClientSocketMessage)
		go c.RunAIClient()
	} else {
		go c.ReadPump()
		go c.WritePump()
	}

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
			fmt.Println("MENSAJE RECIBIDO DE LA ROOM AL CLIENTE")
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

// IMPLEMENTACION PARA IA

// lee mensajes del socket y los manda a la Room
func (c *Client) RunAIClient() error {
	defer func() {
		c.ToAI <- ClientSocketMessage{Type: "EOC", Content: ""}
		select {
		case <-c.Done:
		default:
			close(c.Done)
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			fmt.Println()
			fmt.Println("IA RECIBE MENSAJE: ")
			fmt.Println("Tipo: ", message.Type)
			fmt.Println("Content", message.Content)
			fmt.Println("Extra: ", message.Extra)
			fmt.Println()
			if !ok {
				return nil
			}

			switch message.Type {
			case game.EOC:
				return nil
			case game.InitialState:
				// Filtramos por si es un mensaje causado
				// por una reconexion
				if c.initStateReceived {
					continue
				}
				if message.Extra == strconv.FormatInt(c.AccountID, 10) {
					c.ToAI <- ClientSocketMessage{
						Type:    "Move",
						Content: "",
					}
					response := <-c.FromAI
					c.ToRoom <- response

					// Enviamos el mensaje del movimiento a la IA para que
					// aplique su log
					log := <-c.Send
					c.ToAI <- ClientSocketMessage{
						Type:    "Move",
						Content: log.Content,
					}
					c.initStateReceived = true
				}
			case game.Move:
				fmt.Println("Calculando movimiento")

				c.ToAI <- ClientSocketMessage{
					Type:    "Move",
					Content: message.Content,
				}
				response := <-c.FromAI
				fmt.Println("Enviando mensaje a room: ", response.Content)

				c.ToRoom <- response
				fmt.Println("Mensaje enviado a room")

				// Filtramos el mensaje del movimiento
				log := <-c.Send
				c.ToAI <- ClientSocketMessage{
					Type:    "Move",
					Content: log.Content,
				}
			default:

			}

		case <-c.Done:
			fmt.Println("IA TERMINADA")
			// Si se detecta un error salimos
			return nil
		}
	}
}
