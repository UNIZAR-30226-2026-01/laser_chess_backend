package rt

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/match"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/rating"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/elo"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

// fichero que gestiona las rooms
// una room se encarga de gestionar los mensajes de los dos
// jugadores de una partida :)
// es el intermediario entre el front y el juego

type Room struct {
	Player1    *Client
	Player2    *Client
	isNewMatch bool
	GameInfo   *game.GameInfo
	Game       *game.LaserChessGame
	Registry   *MatchRegistry

	P1Pause bool
	P2Pause bool

	Broadcast chan game.ResponseToRoom

	matchService  *match.MatchService
	ratingService *rating.RatingService
}

func (r *Room) InitRoom(Player1 *Client, Player2 *Client,
	matchService *match.MatchService, ratingService *rating.RatingService,
	isNewMatch bool, GameInfo *game.GameInfo, registry *MatchRegistry) {

	r.Player1 = Player1
	r.Player2 = Player2
	r.P1Pause = false
	r.P2Pause = false
	r.Broadcast = make(chan game.ResponseToRoom, 2)

	r.matchService = matchService
	r.ratingService = ratingService

	r.isNewMatch = isNewMatch
	r.GameInfo = GameInfo

	fmt.Println(GameInfo.BoardType)

	r.Game = &game.LaserChessGame{}
	r.Game.InitLaserChessGame(r.Player1.AccountID, r.Player2.AccountID,
		GameInfo.BoardType, GameInfo.Log, GameInfo.TimeBase, GameInfo.TimeIncrement)

	r.Registry = registry

	r.Registry.RegisterMatch(r.Player1.AccountID, r.Player2.AccountID, r)

	go r.run()
}

func (r *Room) end() {
	fmt.Println("Cierre de la room")

	// Vaciar Broadcast
VaciarBroadcast:
	for {
		select {
		case message := <-r.Broadcast:
			fmt.Println("Broadcast: ", message)
			r.Player1.Send <- message
			r.Player2.Send <- message
		default:
			break VaciarBroadcast
		}
	}

	// TODO: actualizar experiencia

	actualizarElo := r.GameInfo.MatchType != "FRIENDLY" &&
		r.GameInfo.Winner != "NONE"

	// Guardar en BD
	if actualizarElo {
		newP1Rating, newP2Rating, err := r.getUpdatedPlayerRatings()
		fmt.Println("P1Rating: ", newP1Rating, ", P2Rating: ", newP2Rating)
		if err != nil {
			// TODO: gestionar estos errores
		}

		// Llamamos a la transacción pasándole la partida y los Elos
		matchData := match.MatchSaveDTO{
			IsNewMatch: r.isNewMatch,
			GameInfo:   r.GameInfo,
			P1ID:       r.Player1.AccountID,
			P2ID:       r.Player2.AccountID,
			P1Elo:      int32(newP1Rating.Value),
			P2Elo:      int32(newP2Rating.Value),
			Date:       time.Now(),
		}
		fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAA")
		err = r.matchService.SaveMatchResultTx(context.Background(), matchData, *newP1Rating, *newP2Rating)
		if err != nil {
			fmt.Println("EL ERROR: ", err.Error())
		}
	} else {
		p1Elo, p2Elo, err := r.getPlayerRatingValues()
		if err != nil {
			// TODO: gestionar estos errores
		}

		// Guardado simple para amistosas con los elos normales
		matchData := match.MatchSaveDTO{
			IsNewMatch: r.isNewMatch,
			GameInfo:   r.GameInfo,
			P1ID:       r.Player1.AccountID,
			P2ID:       r.Player2.AccountID,
			P1Elo:      p1Elo,
			P2Elo:      p2Elo,
			Date:       time.Now(),
		}
		err = r.matchService.SaveMatch(context.Background(), matchData)
	}

	// TODO: mandar al cliente la info de elo y exp
	fmt.Println("Antes de enviar el EOC a los clientes")
	// Avisar y cerrar los clientes
	r.Player1.mu.RLock()
	if r.Player1.Online {
		r.Player1.Send <- game.ResponseToRoom{Type: game.EOC}
	}
	r.Player1.mu.RUnlock()

	r.Player2.mu.RLock()
	if r.Player2.Online {
		r.Player2.Send <- game.ResponseToRoom{Type: game.EOC}
	}
	r.Player2.mu.RUnlock()

	fmt.Println("Despues de enviar el EOC a los clientes")

	r.Registry.RemoveMatch(r.Player1.AccountID, r.Player2.AccountID)

}

func (r *Room) getPlayerRatingValues() (int32, int32, error) {
	ctx := context.Background()

	p1Data, err := r.ratingService.GetEloByID(ctx, r.Player1.AccountID, r.GameInfo.TimeBase)
	if err != nil {
		return 0, 0, err
	}

	p2Data, err := r.ratingService.GetEloByID(ctx, r.Player1.AccountID, r.GameInfo.TimeBase)
	if err != nil {
		return 0, 0, err
	}

	return p1Data.Value, p2Data.Value, nil
}

func (r *Room) getUpdatedPlayerRatings() (*elo.Rating, *elo.Rating, error) {
	ctx := context.Background()

	// Obtener los ratings actuales de la base de datos
	p1Data, err := r.ratingService.GetEloByID(ctx, r.Player1.AccountID, r.GameInfo.TimeBase)
	if err != nil {
		return nil, nil, err
	}

	p2Data, err := r.ratingService.GetEloByID(ctx, r.Player1.AccountID, r.GameInfo.TimeBase)
	if err != nil {
		return nil, nil, err
	}

	var scoreP1 float64
	switch r.GameInfo.Winner {
	case "P1_WINS":
		scoreP1 = 1.0
	case "P2_WINS":
		scoreP1 = 0.0
	}

	p1Rating := elo.Rating{
		Value:      float64(p1Data.Value),
		Deviation:  float64(p1Data.Deviation),
		Volatility: p1Data.Volatility,
	}
	p2Rating := elo.Rating{
		Value:      float64(p2Data.Value),
		Deviation:  float64(p2Data.Deviation),
		Volatility: p2Data.Volatility,
	}

	p1Rating = elo.ApplyInactivity(p1Rating, p1Data.LastUpdatedAt)
	p2Rating = elo.ApplyInactivity(p2Rating, p2Data.LastUpdatedAt)

	// Calcular nuevos elos
	newP1Rating, newP2Rating := elo.ProcessMatch(p1Rating, p2Rating, scoreP1)

	return &newP1Rating, &newP2Rating, nil
}

func (r *Room) run() {

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
			r.filterMessage(r.Player1, message)
		case message := <-r.Player2.ToRoom:
			r.filterMessage(r.Player2, message)

		case gameMsg := <-r.Game.ToRoom:
			if r.handleGameMessage(gameMsg) {
				// Cerrar la room si acaba la partida
				return
			}
		}
	}
}

func (r *Room) filterMessage(player *Client, message ClientSocketMessage) {
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
		r.managePause(player)
	case game.EOC:
		go r.manageDisconnection(player)
	}
}

func (r *Room) manageDisconnection(player *Client) {
	timeout := time.NewTimer(3 * time.Second)
	defer timeout.Stop()

	fmt.Println("Desconexion detectada")
	if player.AccountID == r.Player1.AccountID {
		r.Player2.Send <- game.ResponseToRoom{Type: game.Disconnection}
	}

	select {
	case <-timeout.C:
		fmt.Println("Desconexion confirmada")
		r.Game.FromRoom <- game.RoomMsg{
			PlayerUid: player.AccountID,
			MsgType:   game.Disconnection,
		}
	case <-player.Reconnect:
	}
}

// El valor que devuelve indica si hay que cerrar la room
func (r *Room) handleGameMessage(response game.ResponseToRoom) bool {
	switch response.Type {
	case game.Move, game.InitialState:
		r.Broadcast <- response

	case game.State, game.Error:
		// Mandar exclusivamente al jugador correspondiente usando el Extra
		if strconv.FormatInt(r.Player1.AccountID, 10) == response.Extra {
			response.Extra = ""
			r.Player1.Send <- response
		} else if strconv.FormatInt(r.Player2.AccountID, 10) == response.Extra {
			response.Extra = ""
			r.Player2.Send <- response
		}

	case game.End:
		r.Broadcast <- response

		r.GameInfo.Winner = response.Content
		r.GameInfo.Termination = response.Extra

		fmt.Println("Winner: ", r.GameInfo.Winner, ", Termination: ", r.GameInfo.Termination)
		r.end()
		return true

	case game.Paused:
		r.Broadcast <- response

		r.GameInfo.Winner = "NONE"
		r.GameInfo.Termination = "UNFINISHED"

		r.end()
		return true
	}

	return false
}

func (r *Room) managePause(player *Client) {
	if r.GameInfo.MatchType != "PRIVATE" {
		player.Send <- game.ResponseToRoom{Type: game.Error,
			Content: "You can't pause a public game"}
		return
	}

	switch player.AccountID {
	case r.Player1.AccountID:
		r.P1Pause = true
		if !r.P2Pause {
			r.Player2.Send <- game.ResponseToRoom{Type: game.PauseRequest}
		}
	case r.Player2.AccountID:
		r.P2Pause = true
		if !r.P1Pause {
			r.Player1.Send <- game.ResponseToRoom{Type: game.PauseRequest}
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
