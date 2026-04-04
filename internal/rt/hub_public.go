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
)

// Info inicial de la partida que se creará

type MatchRequest struct {
	PlayerClient *Client
	PlayerELO    int64
	ELOBracket   int64 // Asigando aqui
	GameMode     int
	ResponseChan *chan *Client
	FoundChan    chan bool // Creado aqui
	CancelChan   chan bool
	ListElement  *list.Element // Asignado aqui
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
	request.ELOBracket = request.PlayerELO / 100
	ph.CheckCreatedQueue(request)

	ph.gameModeQueues[request.GameMode].mu.Lock()
	defer ph.gameModeQueues[request.GameMode].mu.Unlock()

	if ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket].players.Len() == 0 {
		request.FoundChan = make(chan bool)
		request.ListElement = ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket].players.PushBack(request)
	} else {
		// Oponente encontrado
		opponent := ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket].players.Front().Value.(*MatchRequest)
		ph.NotifyMatch(request, opponent)
	}

	fmt.Println("Player1 dentro con elo: ", request.PlayerELO)

	go ph.Search(request)
}

// Busqueda de oponente por elo
func (ph *PublicHub) Search(request *MatchRequest) {
	radius := 1
	// Cada 3 segundos aumentamos el radio hasta 3 niveles de diferencia
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if radius <= 3 {
				radius++
			}
		case <-request.FoundChan:
			return
		case <-request.CancelChan:
			ph.RemoveFromMatchmaking(request)
			return
		}

		ph.FindOpponentInRange(request, radius)

	}
}

func (ph *PublicHub) FindOpponentInRange(request *MatchRequest, radius int) {
	for i := range radius {
		eloDiff := int64(i) * 100

		// Comprobamos el rango superior
		ph.CheckCreatedQueue(request, eloDiff)
		ph.gameModeQueues[request.GameMode].mu.Lock()
		if ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket+eloDiff].players.Len() != 0 {
			opponent := ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket].players.Front().Value.(*MatchRequest)
			if opponent.PlayerClient.AccountID == request.PlayerClient.AccountID {
				ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket].players.PushFront(opponent)
			} else {
				ph.NotifyMatch(request, opponent)
				ph.gameModeQueues[request.GameMode].mu.Unlock()
				return
			}
		}
		ph.gameModeQueues[request.GameMode].mu.Unlock()

		// Comprobamos el rango inferior
		ph.CheckCreatedQueue(request, -eloDiff)
		ph.gameModeQueues[request.GameMode].mu.Lock()
		if ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket-eloDiff].players.Len() != 0 {
			opponent := ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket].players.Front().Value.(*MatchRequest)
			if opponent.PlayerClient.AccountID == request.PlayerClient.AccountID {
				ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket].players.PushFront(opponent)
			} else {
				ph.NotifyMatch(request, opponent)
				ph.gameModeQueues[request.GameMode].mu.Unlock()
				return
			}
		}
		ph.gameModeQueues[request.GameMode].mu.Unlock()
	}
}

func (ph *PublicHub) CheckCreatedQueue(request *MatchRequest, eloDiff_optional ...int64) {

	var eloDiff int64 = 0
	if len(eloDiff_optional) > 0 {
		eloDiff = eloDiff_optional[0]
	}

	// Comprobar si existe el mapa del modo de juego
	if ph.gameModeQueues[request.GameMode] == nil {
		ph.gameModeQueues[request.GameMode] = &GameModeQueue{
			matchmakingQueues: make(map[int64]*MatchmakingQueue),
		}
	}

	ph.gameModeQueues[request.GameMode].mu.Lock()
	defer ph.gameModeQueues[request.GameMode].mu.Unlock()

	// Comprobar si existe el mapa del rando de ELO especifico
	if ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket+eloDiff] == nil {
		ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket+eloDiff] =
			&MatchmakingQueue{
				players: list.New(),
			}
	}

	// Comprobar si existe la lista para este rango de ELO específico
	if ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket+eloDiff].players == nil {
		ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket+eloDiff].players = list.New()
	}
}

func (ph *PublicHub) NotifyMatch(request *MatchRequest, opponent *MatchRequest) {
	*request.ResponseChan <- opponent.PlayerClient
	*opponent.ResponseChan <- request.PlayerClient
	opponent.FoundChan <- true

	//ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket].players.Remove(request.ListElement)
}

func (ph *PublicHub) RemoveFromMatchmaking(request *MatchRequest) {
	ph.gameModeQueues[request.GameMode].mu.Lock()
	defer ph.gameModeQueues[request.GameMode].mu.Unlock()

	ph.gameModeQueues[request.GameMode].matchmakingQueues[request.ELOBracket].players.Remove(request.ListElement)
}
