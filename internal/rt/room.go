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

	P1Pause bool
	P2Pause bool

	Broadcast chan interface{}
}

func (r *Room) InitRoom(Player1 *Client, Player2 *Client, BoardType game.Board_T) {
	r.Player1 = Player1
	r.Player2 = Player2
	r.P1Pause = false
	r.P2Pause = false
	r.Broadcast = make(chan interface{}, 1)

	r.Game = &game.LaserChessGame{}
	r.Game.InitLaserChessGame(r.Player1.AccountID, r.Player2.AccountID, BoardType)

	go r.Run()
}

func (r *Room) End() {
	fmt.Println("Cierre de la room")
	r.Player1.Close()
	r.Player2.Close()
}

func (r *Room) Run() {

	fmt.Println("La partida ha iniciado :)")

	// Mandar estado inicial
	r.Broadcast <- r.GetInitialGameState()

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

	switch game.GameMessageType(message.Type) {
	case game.Move:
		result := r.SendMoveToGame(player.AccountID, message.Content)
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
		state := r.GetGameState()
		player.Send <- state
	case game.Pause:
		r.ManagePause(player)
	}
}

func (r *Room) ManagePause(player *Client) {
	switch player.AccountID {
	case r.Player1.AccountID:
		r.P1Pause = true
		if !r.P2Pause {
			r.Player2.Send <- game.ResponseToRoom{
				Type:    game.PauseRequest,
				Content: "",
				Laser:   "",
			}
		}
	case r.Player2.AccountID:
		r.P2Pause = true
		if !r.P1Pause {
			r.Player1.Send <- game.ResponseToRoom{
				Type:    game.PauseRequest,
				Content: "",
				Laser:   "",
			}
		}
	}

	if r.P1Pause && r.P2Pause {
		result := r.PauseGame()
		r.Broadcast <- result
		// TODO: Sera necesario guardar la partida en BD
		r.End()
	}
}

// FUNCIONES DE COMUNICACIÓN CON EL JUEGO

func (r *Room) SendMoveToGame(accountID int64, move string) game.ResponseToRoom {
	r.Game.FromRoom <- game.RoomMsg{
		PlayerUid:  accountID,
		MsgType:    game.Move,
		MsgContent: move,
	}
	response := <-r.Game.ToRoom
	return response
}

func (r *Room) GetGameState() game.ResponseToRoom {
	r.Game.FromRoom <- game.RoomMsg{
		PlayerUid:  0,
		MsgType:    game.GetState,
		MsgContent: "",
	}
	response := <-r.Game.ToRoom
	return response
}

func (r *Room) GetInitialGameState() game.ResponseToRoom {
	r.Game.FromRoom <- game.RoomMsg{
		PlayerUid:  0,
		MsgType:    game.GetInitialState,
		MsgContent: "",
	}
	response := <-r.Game.ToRoom
	return response
}

func (r *Room) PauseGame() game.ResponseToRoom {
	r.Game.FromRoom <- game.RoomMsg{
		PlayerUid:  0,
		MsgType:    game.Pause,
		MsgContent: "",
	}
	response := <-r.Game.ToRoom
	return response
}
