package rt

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/match"
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

	Broadcast   chan game.ResponseToRoom
	RefreshChan chan bool

	matchService *match.MatchService
}

func (r *Room) InitRoom(Player1 *Client, Player2 *Client,
	matchService *match.MatchService, isNewMatch bool,
	GameInfo *game.GameInfo, registry *MatchRegistry) {

	r.Player1 = Player1
	r.Player2 = Player2
	r.P1Pause = false
	r.P2Pause = false
	r.Broadcast = make(chan game.ResponseToRoom, 2)
	r.RefreshChan = make(chan bool)

	r.matchService = matchService

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
EmptyBroadcast:
	for {
		select {
		case message := <-r.Broadcast:
			fmt.Println("Broadcast: ", message)

			r.Player1.mu.RLock()
			if r.Player1.Online {
				r.Player1.Send <- message
			}
			r.Player1.mu.RUnlock()

			r.Player2.mu.RLock()
			if r.Player2.Online {
				r.Player2.Send <- message
			}
			r.Player2.mu.RUnlock()
		default:
			break EmptyBroadcast
		}
	}

	// Delegamos la logica de fin de partida en el matchService

	r.GameInfo.Log = r.Game.GetLog()

	summary := match.MatchSummaryDTO{
		IsNewMatch: r.isNewMatch,
		GameInfo:   r.GameInfo,
		P1ID:       r.Player1.AccountID,
		P2ID:       r.Player2.AccountID,
		Date:       time.Now(),
	}

	err := r.matchService.FinalizeMatch(context.Background(), summary)
	if err != nil {
		fmt.Println("ERROR finalizando la partida: ", err.Error())
		// TODO: gestionar error
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

func (r *Room) sendOpponentIds() {
	r.Player1.Send <- game.ResponseToRoom{
		Type: game.MatchStart, 
		Content: strconv.FormatInt(r.Player2.AccountID, 10),
	}

	r.Player2.Send <- game.ResponseToRoom{
		Type: game.MatchStart, 
		Content: strconv.FormatInt(r.Player1.AccountID, 10),
	}
}

func (r *Room) run() {

	fmt.Println("La partida ha iniciado :)")

	// Enviar id del oponente a cada jugador
	r.sendOpponentIds()

	// Pedir estado inicial
	r.Game.FromRoom <- game.RoomMsg{MsgType: game.GetInitialState}

	for {
		select {
		case message := <-r.Broadcast:
			fmt.Println("Broadcast: ", message)
			fmt.Println("ROOM: Enviando mensaje de broadcast a jugadores")
			r.Player1.Send <- message
			r.Player2.Send <- message
			fmt.Println("ROOM: Enviado mensaje de broadcast a jugadores")

		case message, ok := <-r.Player1.ToRoom:
			if !ok {
				continue
			}
			r.filterMessage(r.Player1, message)
		case message, ok := <-r.Player2.ToRoom:
			if !ok {
				continue
			}
			r.filterMessage(r.Player2, message)

		case gameMsg := <-r.Game.ToRoom:
			fmt.Println("ROOM: Mensaje recibido de game")
			if r.handleGameMessage(gameMsg) {
				// Cerrar la room si acaba la partida
				return
			}
		case <-r.RefreshChan:
			// No hace nada, es para refrescar las referencias
		}
	}
}

func (r *Room) filterMessage(player *Client, message ClientSocketMessage) {

	if player != r.Player1 && player != r.Player2 {
		return
	}

	// debug
	fmt.Println("Type: ", message.Type)
	fmt.Println("Content: ", message.Content)

	switch game.GameMessageType(message.Type) {
	case game.Move:
		fmt.Println("ROOM: Enviando movimiento a game")
		r.Game.FromRoom <- game.RoomMsg{
			PlayerUid:  player.AccountID,
			MsgType:    game.Move,
			MsgContent: message.Content,
		}
		fmt.Println("ROOM: Movimiento enviado a game")
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
	timeout := time.NewTimer(60 * time.Second)
	defer timeout.Stop()

	fmt.Println("Desconexion detectada")
	if player.AccountID == r.Player1.AccountID {
		r.Player2.Send <- game.ResponseToRoom{Type: game.Disconnection}
	} else {
		r.Player1.Send <- game.ResponseToRoom{Type: game.Disconnection}
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
	fmt.Println("ROOM: Dentro de handleGameMessage: ", response.Type, " ; ", response.Content)
	switch response.Type {
	case game.Move, game.InitialState:

		r.Broadcast <- response
		fmt.Println("ROOM: Enviado mensaje de movimiento a jugadores")

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

func (r *Room) NotifyReconnection(reconected *Client, opponent *Client) {

	// Mensajes al jugador no reconectado
	opponent.Send <- game.ResponseToRoom{Type: game.Reconnection}

	// Mensajes al jugador reconectado
	reconected.Send <- game.ResponseToRoom{
		Type:    game.Reconnection,
		Content: strconv.FormatInt(opponent.AccountID, 10),
	}

	r.Game.FromRoom <- game.RoomMsg{
		PlayerUid: reconected.AccountID,
		MsgType:   game.GetInitialState,
	}

	r.Game.FromRoom <- game.RoomMsg{
		PlayerUid: reconected.AccountID,
		MsgType:   game.GetState,
	}
}

func (r *Room) ReconnectProcedure(reconected *Client, opponent *Client,
	substituted **Client) {
	oldClient := *substituted
	oldClient.Reconnect <- true

	// Cerrar el cliente si esta online
	oldClient.mu.RLock()
	if oldClient.Online {
		oldClient.Send <- game.ResponseToRoom{Type: game.EOC}
	}
	oldClient.mu.RUnlock()

	// Sustituirlo
	*substituted = reconected
	r.RefreshChan <- true
	close(oldClient.Send)

	reconected.Send <- game.ResponseToRoom{
		Type:    game.Reconnection,
		Content: strconv.FormatInt(opponent.AccountID, 10),
	}

	r.NotifyReconnection(reconected, opponent)
}

func (r *Room) ReconnectPlayer(player *Client) {

	if player.AccountID == r.Player1.AccountID {
		r.ReconnectProcedure(player, r.Player2, &r.Player1)
	} else if player.AccountID == r.Player2.AccountID {
		r.ReconnectProcedure(player, r.Player1, &r.Player2)
	}

}
