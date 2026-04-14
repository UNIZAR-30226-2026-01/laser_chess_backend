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
	Content string          `json:"Content,omitempty"` // a veces tampoco se usa
	Extra   string          `json:"Extra,omitempty"`   //campo extra, contiene el laser, que jugador eres o contra que jugador estas
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
		team, timeLeftRed, timeLeftBlue := g.gameEngine.EngineApplyLogToBoard(timeBase)

		switch team {
		case RED_TEAM:
			g.turn = UidRedPlayer
		case BLUE_TEAM:
			g.turn = UidBluePlayer
		}

		g.timerRed = NewGameTimer(time.Duration(timeLeftRed)*time.Second, time.Duration(timeInc)*time.Second)
		g.timerBlue = NewGameTimer(time.Duration(timeLeftBlue)*time.Second, time.Duration(timeInc)*time.Second)

	} else {
		g.turn = UidRedPlayer
		g.timerRed = NewGameTimer(time.Duration(timeBase)*time.Millisecond, time.Duration(timeInc)*time.Second)
		g.timerBlue = NewGameTimer(time.Duration(timeBase)*time.Millisecond, time.Duration(timeInc)*time.Second)
	}

	//Se crean los canales de comunicacón
	g.FromRoom = make(chan RoomMsg, 2)
	g.ToRoom = make(chan ResponseToRoom, 2)

	if g.turn == UidRedPlayer {
		g.timerRed.Start()
		g.timerBlue.Stop()
	} else if g.turn == UidBluePlayer {
		g.timerBlue.Start()
		g.timerRed.Stop()
	}

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
		return RED_TEAM
	}
}

func (g *LaserChessGame) changeTurn() {
	switch g.turn {
	case g.bluePlayer:
		g.timerBlue.Stop()
		g.timerRed.Start()

		g.turn = g.redPlayer
	case g.redPlayer:
		g.timerRed.Stop()
		g.timerBlue.Start()

		g.turn = g.bluePlayer
	}
}

func (g *LaserChessGame) getCurrentTimeLeft() time.Duration {
	switch g.turn {
	case g.redPlayer:
		return g.timerRed.Remaining
	case g.bluePlayer:
		return g.timerBlue.Remaining
	default:
		// Imposible
		return g.timerRed.Remaining
	}
}

func (g *LaserChessGame) changeTimers() {
	switch g.turn {
	case g.bluePlayer:
		g.timerBlue.Stop()
		g.timerRed.Start()

	case g.redPlayer:
		g.timerRed.Stop()
		g.timerBlue.Start()
	}
	fmt.Println("Timer red:" + g.timerRed.Remaining.String())
	fmt.Println("Timer blue:" + g.timerBlue.Remaining.String())
}

// Devuelve true si ha acabado la partida
func (g *LaserChessGame) processMove(message RoomMsg) bool {

	turno := g.getTurn()

	if message.PlayerUid == g.turn {
		// Si es tu turno
		g.changeTimers()

		fmt.Println(message.PlayerUid, ":", message.MsgContent)
		fmt.Println(message.PlayerUid, ":", turno)

		timeLeft := g.getCurrentTimeLeft()
		result, laser, laserInteractionRes, err :=
			g.gameEngine.ProcessTurn(message.MsgContent, turno, timeLeft)

		g.gameEngine.gameBoard.printlaser(laser)
		fmt.Println("ANSWER:", result)

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
			Content: result,
		}

		// Si se ha terminado la partida se notifica de esto
		switch laserInteractionRes {
		case HIT_BLUE_KING:
			g.ToRoom <- ResponseToRoom{
				Type:    End,
				Content: "P1_WINS",
				Extra:   "LASER",
			}
			fmt.Println("END:", result)
			return true
		case HIT_RED_KING:
			g.ToRoom <- ResponseToRoom{
				Type:    End,
				Content: "P2_WINS",
				Extra:   "LASER",
			}
			fmt.Println("END:", result)
			return true
		}

		g.changeTurn()

	} else {
		// No es tu turno
		g.ToRoom <- ResponseToRoom{
			Type:    Error,
			Content: "no es tu turno",
			Extra:   strconv.FormatInt(message.PlayerUid, 10),
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
			Type: Paused,
		}

		return true

	case Disconnection:
		//gestionar pausa del juego
		var winner string
		if message.PlayerUid == g.redPlayer {
			winner = "P2_WINS"
		} else {
			winner = "P1_WINS"
		}
		g.ToRoom <- ResponseToRoom{
			Type:    End,
			Content: winner,
			Extra:   "DISCONNECTION",
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
		fmt.Println("Timer red:" + g.timerRed.Remaining.String())
		fmt.Println("Timer blue:" + g.timerBlue.Remaining.String())
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

func (g *LaserChessGame) GetLog() string {
	return g.gameEngine.gameLog
}
