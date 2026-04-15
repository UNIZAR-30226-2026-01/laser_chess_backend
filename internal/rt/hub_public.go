package rt

// fichero que se encarga de gestionar las conexiones iniciales en
// partidas privadas
// guarda la info de que invitaciones a partidas se han hecho
// y conecta a los dos clientes cuando se acepta la partida

// Para poder acceder y consultar de forma concurrente los datos desde muchos
// puntos a la vez, utilizamos sharding.

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
)

// Variable del fichero

var bracketLenght int = 50

// Info inicial de la partida que se creará

type MatchRequest struct {
	PlayerClient *Client
	PlayerELO    int32
	PlayerRD     int32
	ELOBracket   int64 // Asigando aqui
	GameMode     int
	Ranked       int
	ResponseChan chan *MatchmakingFound
	ErrorChan    chan error
	FoundChan    chan bool // Creado aqui
	CancelChan   chan bool
	ListElement  *list.Element // Asignado aqui
	Board		 int
}

type MatchmakingFound struct {
	Client *Client
	Board  int
}

type MatchmakingQueue struct {
	players *list.List
}

type GameModeQueue struct {
	// Mapa para dividir en bloques de ratings
	matchmakingQueues map[int64]*MatchmakingQueue

	// Mutex para sincronizar la concurrencia
	mu sync.RWMutex
}

type PublicHub struct {

	// Dos mapas para los modos de juego
	rankingOrCasualQueues [2]map[int]*GameModeQueue

	// Players in queue
	playersInQueue map[int64]bool

	// Mutex para la creacion de mapas
	mu sync.RWMutex
}

// Crea un hub para partidas privadas
func NewPublicHub() *PublicHub {
	ph := &PublicHub{}
	ph.rankingOrCasualQueues[0] = make(map[int]*GameModeQueue)
	ph.rankingOrCasualQueues[1] = make(map[int]*GameModeQueue)
	ph.playersInQueue = make(map[int64]bool)
	return ph
}

// Inicia el matchmaking para una partida de tipo Rapid
func (ph *PublicHub) AddPlayerToMatchmaking(request *MatchRequest) {

	// Comprobar que el jugador no está en matchmaking
	if !ph.EnterPlayersInQueue(request.PlayerClient.AccountID) {
		fmt.Println("ERROR: El jugador esta ya dentro de la cola")
		request.ErrorChan <- apierror.ErrAlreadyInQueue
		return
	}

	// Buscar en el rango de ELO efectivo del jugador
	k := int32(2)
	effectiveELO := request.PlayerELO - k*request.PlayerRD
	request.ELOBracket = int64(effectiveELO / int32(bracketLenght))
	request.ListElement = nil
	request.FoundChan = make(chan bool, 1)
	ph.CheckCreatedQueue(request)

	// Si ya encontramos oponente terminamos
	if ph.CheckBracket(request, request.ELOBracket) {
		return
	}

	// Si no nos metemos en la lista de espera
	ph.AddToQueue(request)

	fmt.Println("Player1 dentro con elo: ", request.PlayerELO)

	// Y buscamos
	go ph.Search(request)

	return
}

// Busqueda de oponente por ELO
func (ph *PublicHub) Search(request *MatchRequest) {
	radius := 1
	// Cada 3 segundos aumentamos el radio hasta 3 niveles de diferencia
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {

		case <-request.FoundChan:
			ph.RemoveFromQueue(request)
			return
		case <-request.CancelChan:
			ph.RemoveFromQueue(request)
			return
		case <-ticker.C:
			if radius <= 6 {
				radius++
			}
			ph.FindOpponentInRange(request, radius)

		}
	}
}

func (ph *PublicHub) FindOpponentInRange(request *MatchRequest, radius int) {

	for i := range radius {
		eloDiff := int64(i)
		upBracket := request.ELOBracket + eloDiff
		downBracket := request.ELOBracket - eloDiff

		// Comprobamos el rango superior
		ph.CheckCreatedQueue(request, eloDiff)
		ph.CheckBracket(request, upBracket)

		// Comprobamos el rango inferior
		ph.CheckCreatedQueue(request, -eloDiff)
		ph.CheckBracket(request, downBracket)

	}
}

func (ph *PublicHub) CheckBracket(request *MatchRequest, bracket int64) bool {
	queue := ph.rankingOrCasualQueues[request.Ranked][request.GameMode]

	queue.mu.Lock()
	defer queue.mu.Unlock()

	if queue.matchmakingQueues[bracket].players.Len() != 0 {
		opponent :=
			queue.matchmakingQueues[bracket].players.Front().Value.(*MatchRequest)
		if opponent.PlayerClient.AccountID != request.PlayerClient.AccountID {
			// Sacar a los jugadores de la cola
			ph.ExitPlayersInQueue(request.PlayerClient.AccountID)
			ph.ExitPlayersInQueue(opponent.PlayerClient.AccountID)
			ph.NotifyMatch(request, opponent)
			return true
		}
	}

	return false
}

func (ph *PublicHub) CheckCreatedQueue(request *MatchRequest, eloDiff_optional ...int64) {

	var eloDiff int64 = 0
	if len(eloDiff_optional) > 0 {
		eloDiff = eloDiff_optional[0]
	}

	// Comprobar si existe el mapa del modo de juego
	ph.mu.Lock()

	queue := ph.rankingOrCasualQueues[request.Ranked][request.GameMode]

	if queue == nil {
		queue = &GameModeQueue{
			matchmakingQueues: make(map[int64]*MatchmakingQueue),
		}
		ph.rankingOrCasualQueues[request.Ranked][request.GameMode] = queue
	}

	ph.mu.Unlock()

	queue.mu.Lock()
	defer queue.mu.Unlock()

	// Comprobar si existe el mapa del rando de ELO especifico
	if queue.matchmakingQueues[request.ELOBracket+eloDiff] == nil {
		queue.matchmakingQueues[request.ELOBracket+eloDiff] =
			&MatchmakingQueue{
				players: list.New(),
			}
	}

	// Comprobar si existe la lista para este rango de ELO específico
	if queue.matchmakingQueues[request.ELOBracket+eloDiff].players == nil {
		queue.matchmakingQueues[request.ELOBracket+eloDiff].players = list.New()
	}
}

func (ph *PublicHub) NotifyMatch(request *MatchRequest, opponent *MatchRequest) {
	request.ResponseChan <- &MatchmakingFound{
		Client: opponent.PlayerClient,
		Board:  opponent.Board,
	}
	opponent.ResponseChan <- &MatchmakingFound{
		Client: request.PlayerClient,
		Board:  request.Board,
	}
	request.FoundChan <- true
	opponent.FoundChan <- true

	fmt.Println("Notificacion enviada")

}

func (ph *PublicHub) AddToQueue(request *MatchRequest) {
	queue := ph.rankingOrCasualQueues[request.Ranked][request.GameMode]

	queue.mu.Lock()
	defer queue.mu.Unlock()

	request.ListElement =
		queue.matchmakingQueues[request.ELOBracket].players.PushBack(request)
}

func (ph *PublicHub) RemoveFromQueue(request *MatchRequest) {
	queue := ph.rankingOrCasualQueues[request.Ranked][request.GameMode]

	queue.mu.Lock()
	defer queue.mu.Unlock()

	queue.matchmakingQueues[request.ELOBracket].players.Remove(request.ListElement)
}

func (ph *PublicHub) EnterPlayersInQueue(id int64) bool {
	ph.mu.Lock()
	defer ph.mu.Unlock()
	if ph.playersInQueue[id] {
		return false
	}
	ph.playersInQueue[id] = true
	return true
}

func (ph *PublicHub) ExitPlayersInQueue(id int64) {
	ph.mu.Lock()
	defer ph.mu.Unlock()
	ph.playersInQueue[id] = false
}
