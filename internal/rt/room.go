package rt

import (
	"fmt"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
	"github.com/gin-gonic/gin"
)

// fichero que gestiona las rooms
// una room se encarga de gestionar los mensajes de los dos
// jugadores de una partida
// es el intermediario entre el front y el juego

type Room struct {
	Player1 *Client
	Player2 *Client
	Game    *game.LaserChessGame

	FromP1 chan ClientSocketMessage
	FromP2 chan ClientSocketMessage

	ConP1     chan interface{}
	ConP2     chan interface{}
	Broadcast chan interface{}
}

func (r *Room) InitRoom(Player1 *Client, Player2 *Client, BoardType game.Board_T) {
	r.Player1 = Player1
	r.Player2 = Player2
	r.ConP1 = make(chan interface{})
	r.ConP2 = make(chan interface{})
	r.Broadcast = make(chan interface{})

	r.Game = &game.LaserChessGame{}
	r.Game.InitLaserChessGame(r.Player1.AccountID, r.Player2.AccountID, BoardType)

	go r.Run()
	// Notificar a ambos clientes que la partida ha empezado
	startMsg := gin.H{"type": "MatchStarted"}
	r.Broadcast <- startMsg
}

func (r *Room) Run() {

	fmt.Println("La partida ha iniciado :)")

	for {
		select {
		case message := <-r.ConP1:
			r.Player2.Send <- message

		case message := <-r.ConP2:
			r.Player2.Send <- message

		case message := <-r.Broadcast:
			r.Player1.Send <- message
			r.Player2.Send <- message
		case message := <-r.FromP1:
			r.FilterMessage(r.Player1, message)
		case message := <-r.FromP2:
			r.FilterMessage(r.Player2, message)
		}
	}
}

func (r *Room) FilterMessage(player *Client, message ClientSocketMessage) {
	switch message.Type {
	case "Move":
		r.MakeMove(player.AccountID, message.Content)
	case "GetState":
		state := r.GetGameState()
		player.Send <- state
	}
}

func (r *Room) MakeMove(accountID int64, move string) {
	if accountID != r.Player1.AccountID && accountID != r.Player2.AccountID {
		return
	}
	state := r.SendMoveToGame(accountID, move)
	r.Broadcast <- state
}

// FUNCIONES DE COMUNICACIÓN CON EL JUEGO

func (r *Room) SendMoveToGame(accountID int64, move string) string {
	r.Game.FromRoom <- game.RoomMsg{PlayerUid: accountID, MsgType: "Move", MsgContent: move}
	response := <-r.Game.ToRoom
	return response.MsgContent
}

func (r *Room) GetGameState() string {
	r.Game.FromRoom <- game.RoomMsg{PlayerUid: 0, MsgType: "GetState", MsgContent: ""}
	response := <-r.Game.ToRoom
	return response.MsgContent
}
