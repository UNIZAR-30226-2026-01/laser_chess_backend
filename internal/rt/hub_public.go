package rt

// fichero que se encarga de gestionar las conexiones iniciales en
// partidas privadas
// guarda la info de que invitaciones a partidas se han hecho
// y conecta a los dos clientes cuando se acepta la partida

// Para poder acceder y consultar de forma concurrente los datos desde muchos
// puntos a la vez, utilizamos sharding.

import (
	"container/list"
	"sync"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

// Info inicial de la partida que se creará
type NewMatchInfo struct {
	ChallengerClient *Client
	ChallengedId     int64
	Board            game.Board_T
	StartingTime     int32
	TimeIncrement    int32
}

type MatchRequest struct {
	PlayerID     int64
	PlayerELO    int64
	GameMode     int
	ResponseChan chan int64
	FoundChan    chan int64
	CancelChan   chan bool
	ListElement  *list.Element
}

type MatchmakingQueue struct {
	mu      sync.RWMutex
	players *list.List
}

type GameModeQueue struct {
	// Mapa para dividir en bloques de ratings
	matchmakingQueues map[int64]*MatchmakingQueue
}

type PublicHub struct {
	// registro de partidas activas
	registry *MatchRegistry

	// Un mapa para cada modo de juego
	gameModeQueues map[int]*GameModeQueue
}

// Crea un hub para partidas privadas
func NewPublicHub(r *MatchRegistry) *PublicHub {
	ph := &PublicHub{registry: r}
	ph.gameModeQueues = make(map[int]*GameModeQueue)
	return ph
}

// Inicia el matchmaking para una partida de tipo Rapid
func (ph *PublicHub) AddPlayerToMatchmaking(request *MatchRequest) {

	// Buscar en el rango de ELO del jugador

	// Comprobar si existe el mapa del modo de juego
	if ph.gameModeQueues[request.GameMode].matchmakingQueues == nil {
		ph.gameModeQueues[request.GameMode].matchmakingQueues =
			make(map[int64]*MatchmakingQueue)
	}

	bracket := request.PlayerELO / 100

	ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].mu.Lock()
	defer ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].mu.Unlock()
	// Comprobar si existe la lista para este rango de ELO específico
	if ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].players == nil {
		ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].players = list.New()
	}

	if ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].players.Len() == 0 {
		request.ListElement = ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].players.PushBack(request)
	} else {
		// Oponente encontrado
		opponent := ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].players.Front().Value.(*MatchRequest)
		ph.NotifyMatch(request, opponent)
	}

	go ph.Search(request)
}

// Elimina un reto, borra la info del challenger y lo quita de la lista del challenged
func (ph *PublicHub) RemoveFromMatchmaking(request *MatchRequest) {
	bracket := request.PlayerELO / 100
	ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].mu.Lock()
	defer ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].mu.Unlock()

	ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].players.Remove(request.ListElement)
}

// Busqueda de oponente por elo
func (ph *PublicHub) Search(request *MatchRequest) {
	radius := 0
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if radius < 3 {
				radius++
			}
		case opponentID := <-request.FoundChan:
			request.ResponseChan <- opponentID
			return
		case <-request.CancelChan:
			ph.RemoveFromMatchmaking(request)
			return
		}

		ph.FindOpponentInRange(request, radius)

	}
}

func (ph *PublicHub) FindOpponentInRange(request *MatchRequest, radius int) {
	bracket := request.PlayerELO / 100
	for i := range radius {
		eloDiff := int64(i) * 100
		if ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket+eloDiff].players != nil {
			opponent := ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].players.Front().Value.(*MatchRequest)
			ph.NotifyMatch(request, opponent)
		}
		if ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket-eloDiff].players != nil {
			opponent := ph.gameModeQueues[request.GameMode].matchmakingQueues[bracket].players.Front().Value.(*MatchRequest)
			ph.NotifyMatch(request, opponent)
		}
	}
}

func (ph *PublicHub) NotifyMatch(player1 *MatchRequest, player2 *MatchRequest) {
	player1.ResponseChan <- player2.PlayerID
	player2.ResponseChan <- player1.PlayerID
}
