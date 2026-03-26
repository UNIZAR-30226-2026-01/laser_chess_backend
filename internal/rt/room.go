package rt

import (
	"fmt"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

// fichero que gestiona las rooms
// una room se encarga de gestionar los mensajes de los dos
// jugadores de una partida
// es el intermediario entre el front y el juego

type Room struct {
	Player1 *Client
	Player2 *Client
	Game    *game.LaserChessGame

	Broadcast chan interface{}
}

func (r *Room) InitRoom(Player1 *Client, Player2 *Client, BoardType game.Board_T) {
	r.Player1 = Player1
	r.Player2 = Player2
	r.Broadcast = make(chan interface{}, 1)

	r.Game = &game.LaserChessGame{}
	r.Game.InitLaserChessGame(r.Player1.AccountID, r.Player2.AccountID, BoardType)

	go r.Run()
}

func (r *Room) Run() {

	fmt.Println("La partida ha iniciado :)")

	for {
		select {
		case message := <-r.Broadcast:
			fmt.Println("Broadcast: ", message)
			r.Player1.Send <- message
			r.Player2.Send <- message

		case message := <-r.Player1.ToRoom:
			r.FilterMessage(r.Player1, message)
		case message := <-r.Player2.ToRoom:
			r.FilterMessage(r.Player2, message)
		}
	}
}

func (r *Room) FilterMessage(player *Client, message ClientSocketMessage) {
	// debug
	fmt.Println("Type: ", message.Type)
	fmt.Println("Content: ", message.Content)

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

func (r *Room) SendMoveToGame(accountID int64, move string) game.ResponseToRoom {
	r.Game.FromRoom <- game.RoomMsg{PlayerUid: accountID, MsgType: "Move", MsgContent: move}
	response := <-r.Game.ToRoom
	return response
}

func (r *Room) GetGameState() string {
	r.Game.FromRoom <- game.RoomMsg{PlayerUid: 0, MsgType: "GetState", MsgContent: ""}
	response := <-r.Game.ToRoom
	return response.MsgContent
}
