package rt

import (
	"context"
	"fmt"
	"strconv"
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

	Broadcast chan game.ResponseToRoom

	MatchService *match.MatchService
}

func (r *Room) InitRoom(Player1 *Client, Player2 *Client,
	MatchService *match.MatchService, IsNewMatch bool, GameInfo *game.GameInfo) {

	r.Player1 = Player1
	r.Player2 = Player2
	r.P1Pause = false
	r.P2Pause = false
	r.Broadcast = make(chan game.ResponseToRoom, 1)
	r.MatchService = MatchService
	r.IsNewMatch = IsNewMatch
	r.GameInfo = GameInfo

	r.Game = &game.LaserChessGame{}
	r.Game.InitLaserChessGame(r.Player1.AccountID, r.Player2.AccountID,
		GameInfo.BoardType, GameInfo.Log, GameInfo.TimeBase, GameInfo.TimeIncrement)

	go r.Run()
}

func (r *Room) End() {
	fmt.Println("Cierre de la room")

	// Guardar la partida en BD
	err := r.SaveMatchInDB()

	if err != nil {
		fmt.Println(err)
	}

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
			MovementHistory: r.Game.GetCurrentState(),
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
			MovementHistory: r.Game.GetCurrentState(),
			TimeBase:        int32(r.GameInfo.TimeBase),
			TimeIncrement:   int32(r.GameInfo.TimeIncrement),
			MatchID:         r.GameInfo.MatchID,
		})
	}
	return err
}

func (r *Room) Run() {

	fmt.Println("La partida ha iniciado :)")

	// Pedir estado inicial
	r.Game.FromRoom <- game.RoomMsg{MsgType: game.GetInitialState}

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

		case gameMsg := <-r.Game.ToRoom:
			if r.HandleGameMessage(gameMsg) {
				// Cerrar la room si acaba la partida
				return
			}
		}
	}
}

func (r *Room) FilterMessage(player *Client, message ClientSocketMessage) {
	// debug
	fmt.Println("Type: ", message.Type)
	fmt.Println("Content: ", message.Content)

	switch game.GameMessageType(message.Type) {
	case game.Move:
		r.Game.FromRoom <- game.RoomMsg{
			PlayerUid:  player.AccountID,
			MsgType:    game.Move,
			MsgContent: message.Content,
		}
	case game.GetState:
		r.Game.FromRoom <- game.RoomMsg{
			PlayerUid: player.AccountID,
			MsgType:   game.GetState,
		}
	case game.Pause:
		r.ManagePause(player)
	}
}

// El valor que devuelve indica si hay que cerrar la room
func (r *Room) HandleGameMessage(response game.ResponseToRoom) bool {
	switch response.Type {
	case game.Move, game.InitialState:
		r.Broadcast <- response

	case game.State, game.Error:
		// Mandar exclusivamente al jugador correspondiente usando el Extra
		if strconv.FormatInt(r.Player1.AccountID, 10) == response.Extra {
			r.Player1.Send <- response
		} else if strconv.FormatInt(r.Player2.AccountID, 10) == response.Extra {
			r.Player2.Send <- response
		}

	case game.End:
		r.Player2.Send <- response
		r.Player1.Send <- response

		r.GameInfo.Winner = response.Content
		r.GameInfo.Termination = response.Extra

		r.End()
		return true

	case game.Paused:
		r.Player2.Send <- response
		r.Player1.Send <- response

		r.GameInfo.Winner = "NONE"
		r.GameInfo.Termination = "UNFINISHED"

		r.End()
		return true
	}

	return false
}

func (r *Room) ManagePause(player *Client) {
	switch player.AccountID {
	case r.Player1.AccountID:
		r.P1Pause = true
		if !r.P2Pause {
			r.Player2.Send <- game.ResponseToRoom{Type: game.PauseRequest, Content: ""}
		}
	case r.Player2.AccountID:
		r.P2Pause = true
		if !r.P1Pause {
			r.Player1.Send <- game.ResponseToRoom{Type: game.PauseRequest, Content: ""}
		}
	}

	if r.P1Pause && r.P2Pause {
		r.Game.FromRoom <- game.RoomMsg{
			PlayerUid:  0,
			MsgType:    game.Pause,
			MsgContent: "",
		}
	}
}
