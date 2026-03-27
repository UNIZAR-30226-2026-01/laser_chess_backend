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
	EndChan   chan interface{}
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

	// Mandar estado inicial
	r.Broadcast <- r.getInitialGameState()

	for {
		select {
		case message := <-r.Broadcast:
			fmt.Println("Broadcast: ", message)
			r.Player1.Send <- message
			r.Player2.Send <- message

		case message := <-r.Player1.ToRoom:
			r.filterMessage(r.Player1, message)
		case message := <-r.Player2.ToRoom:
			r.filterMessage(r.Player2, message)
		case <-r.EndChan:
			return
		}
	}
}

func (r *Room) filterMessage(player *Client, message ClientSocketMessage) {
	// debug
	fmt.Println("Type: ", message.Type)
	fmt.Println("Content: ", message.Content)

	switch game.GameMessageType(message.Type) {
	case game.Move:
		result := r.sendMoveToGame(player.AccountID, message.Content)
		switch result.Type {
		case game.Move:
			r.Broadcast <- result
		case game.End:
			r.Broadcast <- result
			r.End()
		case game.Error:
			player.Send <- result
		}
	case game.GetState:
		state := r.getGameState()
		player.Send <- state
	}
}

// FUNCIONES DE COMUNICACIÓN CON EL JUEGO

func (r *Room) sendMoveToGame(accountID int64, move string) game.ResponseToRoom {
	r.Game.FromRoom <- game.RoomMsg{
		PlayerUid:  accountID,
		MsgType:    game.Move,
		MsgContent: move,
	}
	response := <-r.Game.ToRoom
	return response
}

func (r *Room) getGameState() game.ResponseToRoom {
	r.Game.FromRoom <- game.RoomMsg{
		PlayerUid:  0,
		MsgType:    game.GetState,
		MsgContent: "",
	}
	response := <-r.Game.ToRoom
	return response
}

func (r *Room) getInitialGameState() game.ResponseToRoom {
	r.Game.FromRoom <- game.RoomMsg{
		PlayerUid:  0,
		MsgType:    game.GetInitialState,
		MsgContent: "",
	}
	response := <-r.Game.ToRoom
	return response
}

func (r *Room) End() {
	r.Player1.Close()
	r.Player2.Close()
	r.EndChan <- ""
}
