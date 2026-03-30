package rt

import (
	"context"
	"fmt"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/match"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db"
	sqlc "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
	"github.com/jackc/pgx/v5/pgtype"
)

// fichero que gestiona las rooms
// una room se encarga de gestionar los mensajes de los dos
// jugadores de una partida :)
// es el intermediario entre el front y el juego

type Room struct {
	Player1    *Client
	Player2    *Client
	IsNewMatch bool
	GameInfo   *game.GameInfo
	Game       *game.LaserChessGame

	P1Pause bool
	P2Pause bool

	Broadcast chan interface{}

	MatchService *match.MatchService
}

func (r *Room) InitRoom(Player1 *Client, Player2 *Client,
	MatchService *match.MatchService, IsNewMatch bool, GameInfo *game.GameInfo) {

	r.Player1 = Player1
	r.Player2 = Player2
	r.P1Pause = false
	r.P2Pause = false
	r.Broadcast = make(chan interface{}, 1)
	r.MatchService = MatchService
	r.IsNewMatch = IsNewMatch
	r.GameInfo = GameInfo

	r.Game = &game.LaserChessGame{}
	r.Game.InitLaserChessGame(r.Player1.AccountID, r.Player2.AccountID,
		GameInfo.BoardType, GameInfo.Log)

	go r.Run()
}

func (r *Room) End() {
	fmt.Println("Cierre de la room")

	// Guardar la partida en BD
	err := r.SaveMatchInDB()

	if err != nil {
		fmt.Println(err)
	}

	r.Player1.Close()
	r.Player2.Close()
}

func (r *Room) SaveMatchInDB() error {
	var err error
	if r.IsNewMatch {
		fmt.Println("Cierre de la room con match nueva")
		_, err = r.MatchService.Create(context.Background(), &match.MatchDTO{
			P1ID:            r.Player1.AccountID,
			P2ID:            r.Player2.AccountID,
			P1Elo:           0, // cambiar
			P2Elo:           0, // cambiar
			Date:            pgtype.Timestamptz{Time: time.Now(), Valid: true},
			Winner:          sqlc.Winner(r.GameInfo.Winner),
			Termination:     sqlc.Termination(r.GameInfo.Termination),
			MatchType:       sqlc.MatchType(r.GameInfo.MatchType),
			Board:           db.IntToBoard[r.GameInfo.BoardType],
			MovementHistory: r.GetGameState().Content,
			TimeBase:        int32(r.GameInfo.TimeBase),
			TimeIncrement:   int32(r.GameInfo.TimeIncrement),
		})
	} else {
		_, err = r.MatchService.UpdateMatch(context.Background(), &sqlc.UpdateMatchParams{
			P1ID:            r.Player1.AccountID,
			P2ID:            r.Player2.AccountID,
			P1Elo:           0, // cambiar
			P2Elo:           0, // cambiar
			Date:            pgtype.Timestamptz{Time: time.Now(), Valid: true},
			Winner:          sqlc.Winner(r.GameInfo.Winner),
			Termination:     sqlc.Termination(r.GameInfo.Termination),
			MatchType:       sqlc.MatchType(r.GameInfo.MatchType),
			Board:           db.IntToBoard[r.GameInfo.BoardType],
			MovementHistory: r.GetGameState().Content,
			TimeBase:        int32(r.GameInfo.TimeBase),
			TimeIncrement:   int32(r.GameInfo.TimeIncrement),
			MatchID:         r.GameInfo.MatchID,
		})
	}
	return err
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
			}
		}
	case r.Player2.AccountID:
		r.P2Pause = true
		if !r.P1Pause {
			r.Player1.Send <- game.ResponseToRoom{
				Type:    game.PauseRequest,
				Content: "",
			}
		}
	}

	if r.P1Pause && r.P2Pause {
		result := r.PauseGame()
		r.Broadcast <- result
		r.GameInfo.Winner = "NONE"
		r.GameInfo.Termination = "UNFINISHED"
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
