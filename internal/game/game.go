package game

import (
	"fmt"
	"strconv"
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

type LaserChessGame struct {
	redPlayer  int64
	bluePlayer int64

	turn int64

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
	BoardType Board_T, Log string) {
	//Rellenan los datos relevantes
	g.redPlayer = UidRedPlayer
	g.bluePlayer = UidBluePlayer

	g.gameEngine.gameLog = Log
	
	//Estado inicial de la partida
	g.gameEngine.InitEngine(BoardType)

	//si el log no está vacío hay que reconstruir el estado
	if g.gameEngine.gameLog != "" {
		team, _, _ := g.gameEngine.ApplyLogToBoard()
		
		switch team{
		case RED_TEAM:
			g.turn = UidRedPlayer
		case BLUE_TEAM:
			g.turn = UidBluePlayer
		}

	}else{
		g.turn = UidRedPlayer
	}


	//Se crean los canales de comunicacón
	g.FromRoom = make(chan RoomMsg)
	g.ToRoom = make(chan ResponseToRoom)
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
		g.turn = g.redPlayer
	case g.redPlayer:
		g.turn = g.bluePlayer
	}
}

func (g *LaserChessGame) processMove(message RoomMsg) {

	turno := g.getTurn()

	if message.PlayerUid == g.turn {
		// Si es tu turno

		fmt.Println(message.PlayerUid, ":", message.MsgContent)
		fmt.Println(message.PlayerUid, ":", turno)
		resul, laser, laserInteractionRes, err := g.gameEngine.ProcessTurn(message.MsgContent, turno)
		g.gameEngine.gameBoard.printlaser(laser)
		fmt.Println("ANSWER:", resul)

		// Si hay un error, se notifica de este
		if err != nil {
			g.ToRoom <- ResponseToRoom{
				Type:    Error,
				Content: err.Error(),
			}
			return
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
				Content: "P2_WINS",
			}
			return
		case HIT_RED_KING:
			g.ToRoom <- ResponseToRoom{
				Type:    End,
				Content: "P1_WINS",
			}
			return
		}

		g.changeTurn()

	} else {
		// No es tu turno
		g.ToRoom <- ResponseToRoom{
			Type:    Error,
			Content: "no es tu turno",
		}
	}
}

func (g *LaserChessGame) Run() {
	for message := range g.FromRoom {
		switch message.MsgType {
		case Move:
			g.processMove(message)
		case GetState:
			g.ToRoom <- ResponseToRoom{
				Type:    State,
				Content: g.gameEngine.GetState(),
				Extra:   "",
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
		}
	}
}
