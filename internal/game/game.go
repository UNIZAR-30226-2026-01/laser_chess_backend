package game

import (
	"fmt"
	"strconv"
	"time"
)

type RoomMsg struct {
	PlayerUid  int64
	MsgType    GameMessageType
	MsgContent string
}

type ResponseToRoom struct {
	Type    GameMessageType `json:"Type"`
	Content string          `json:"Content"`
	Extra   string          `json:"Extra,omitempty"` //campo extra, contiene o el laser, o que jugador eres
}

type GameInfo struct {
	Log           string
	BoardType     Board_T
	TimeBase      int32
	TimeIncrement int32
	Winner        string
	Termination   string
	MatchType     string
	MatchID       int64
}

type LaserChessGame struct {
	redPlayer  int64
	bluePlayer int64

	turn int64

	timerRed  *GameTimer
	timerBlue *GameTimer

	FromRoom   chan RoomMsg
	ToRoom     chan ResponseToRoom
	gameEngine GameEngine
}

/*
* Desc: Esta funcion realiza el procesamiento del recorrido del haz laser en el tablero
*
* --- Parametros ---
* uidRedPlayer int64 - Es el uid del jugador rojo.
* uidBluePlayer int64 - Es el uid del jugador azul.
* --- Resultados ---
* LaserChessGame - Es la nueva instancia del juego inicializada para comenzar a jugar
 */
func (g *LaserChessGame) InitLaserChessGame(UidRedPlayer int64, UidBluePlayer int64,
	BoardType Board_T, Log string, timeBase int32, timeInc int32) {
	//Rellenan los datos relevantes
	g.redPlayer = UidRedPlayer
	g.bluePlayer = UidBluePlayer

	g.gameEngine.gameLog = Log

	//Estado inicial de la partida
	g.gameEngine.InitEngine(BoardType)

	//si el log no está vacío hay que reconstruir el estado
	if g.gameEngine.gameLog != "" {
		team, _, _ := g.gameEngine.ApplyLogToBoard()

		switch team {
		case RED_TEAM:
			g.turn = UidRedPlayer
		case BLUE_TEAM:
			g.turn = UidBluePlayer
		}

	} else {
		g.turn = UidRedPlayer
	}

	//Se crean los canales de comunicacón
	g.FromRoom = make(chan RoomMsg, 2)
	g.ToRoom = make(chan ResponseToRoom, 2)

	// Inicializar timers
	g.timerRed = NewGameTimer(time.Duration(timeBase)*time.Second, time.Duration(timeInc)*time.Second)
	g.timerBlue = NewGameTimer(time.Duration(timeBase)*time.Second, time.Duration(timeInc)*time.Second)
	// El timer no empieza hasta el primer movimiento

	go g.Run()

	fmt.Println("Game inicializado")
}

func (g *LaserChessGame) getTurn() team_T {
	switch g.turn {
	case g.bluePlayer:
		return BLUE_TEAM
	case g.redPlayer:
		return RED_TEAM
	default:
		// Imposible
		fmt.Println("Error al calcular el turno")
		return RED_TEAM
	}
}

func (g *LaserChessGame) changeTurn() {
	switch g.turn {
	case g.bluePlayer:
		g.timerBlue.Stop()
		g.timerRed.Start()
		fmt.Println(g.timerBlue.Remaining)

		g.turn = g.redPlayer
	case g.redPlayer:
		g.timerRed.Stop()
		g.timerBlue.Start()
		fmt.Println(g.timerRed.Remaining)

		g.turn = g.bluePlayer
	}
	fmt.Println(g.turn)
}

// Devuelve true si ha acabado la partida
func (g *LaserChessGame) processMove(message RoomMsg) bool {

	turno := g.getTurn()

	if message.PlayerUid == g.turn {
		// Si es tu turno

		fmt.Println(message.PlayerUid, ":", message.MsgContent)
		fmt.Println(message.PlayerUid, ":", turno)
		var timestamp time.Duration
		switch g.turn {
		case g.bluePlayer:
			timestamp = g.timerBlue.Remaining
		case g.redPlayer:
			timestamp = g.timerRed.Remaining
		}
		resul, laser, laserInteractionRes, err := g.gameEngine.ProcessTurn(message.MsgContent, turno, time.Duration(timestamp.Abs().Seconds()))
		g.gameEngine.gameBoard.printlaser(laser)
		fmt.Println("ANSWER:", resul)

		// Si hay un error, se notifica de este
		if err != nil {
			g.ToRoom <- ResponseToRoom{
				Type:    Error,
				Content: err.Error(),
				Extra:   strconv.FormatInt(message.PlayerUid, 10),
			}
			return false
		}

		g.ToRoom <- ResponseToRoom{
			Type:    Move,
			Content: resul,
			Extra:   fmt.Sprint(formatearLaserPath(laser)),
		}

		// Si se ha terminado la partida se notifica de esto
		switch laserInteractionRes {
		case HIT_BLUE_KING:
			g.ToRoom <- ResponseToRoom{
				Type:    End,
				Content: "P1_WINS",
				Extra:   "LASER",
			}
			fmt.Println("END:", resul)
			return true
		case HIT_RED_KING:
			g.ToRoom <- ResponseToRoom{
				Type:    End,
				Content: "P2_WINS",
				Extra:   "LASER",
			}
			fmt.Println("END:", resul)
			return true
		}

		g.changeTurn()

	} else {
		// No es tu turno
		g.ToRoom <- ResponseToRoom{
			Type:    Error,
			Content: "no es tu turno",
		}
	}

	return false
}

// Devuelve true si ha acabado o pausa la partida
func (g *LaserChessGame) HandleRoomMsg(message RoomMsg) bool {
	switch message.MsgType {
	case Move:
		return g.processMove(message)
	case GetState:
		g.ToRoom <- ResponseToRoom{
			Type:    State,
			Content: g.gameEngine.GetState(),
			Extra:   strconv.FormatInt(message.PlayerUid, 10),
		}
	case GetInitialState:
		initialState := g.gameEngine.getInitialState()
		g.ToRoom <- ResponseToRoom{
			Type:    InitialState,
			Content: initialState,
			Extra:   strconv.FormatInt(g.redPlayer, 10),
		}
	case Pause:
		//gestionar pausa del juego
		g.ToRoom <- ResponseToRoom{
			Type:    Paused,
			Content: "", // quizas manda algo aqui
			Extra:   "",
		}

		return true
	}

	return false
}

func (g *LaserChessGame) GetCurrentState() string {
	return g.gameEngine.GetState()
}

func (g *LaserChessGame) Run() {
	defer func() {
		g.timerRed.Stop()
		g.timerBlue.Stop()
	}()

	for {
		select {
		case message := <-g.FromRoom:
			if g.HandleRoomMsg(message) {

				return
			}

		case <-g.timerRed.Expired:
			g.ToRoom <- ResponseToRoom{
				Type:    End,
				Content: "P2_WINS",
				Extra:   "OUT_OF_TIME",
			}
			return
		case <-g.timerBlue.Expired:
			g.ToRoom <- ResponseToRoom{
				Type:    End,
				Content: "P1_WINS",
				Extra:   "OUT_OF_TIME",
			}
			return
		}
	}
}
