package game

import (
	"fmt"

	boardtemplates "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game/boardTemplates"
)

type RoomMsg struct {
	PlayerUid  int64
	MsgType    GameMessageType
	MsgContent string
}

type ResponseToRoom struct {
	Type    GameMessageType
	Content string
	Laser   string
}

type LaserChessGame struct {
	redPlayer  int64
	bluePlayer int64

	turn int64

	FromRoom chan RoomMsg
	ToRoom   chan ResponseToRoom

	gameBoard Board
	boardType Board_T
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
	BoardType Board_T) {
	g.redPlayer = UidRedPlayer
	g.bluePlayer = UidBluePlayer
	g.turn = UidRedPlayer
	g.gameBoard, _ = InitBoard("boardTemplates/ace.csv") // TODO: NO poner el nombre del csv, poner un switch case

	g.FromRoom = make(chan RoomMsg)
	g.ToRoom = make(chan ResponseToRoom)

	go g.Run()

	fmt.Println("Game inicializado")
}

func (g *LaserChessGame) getInitialState() string {
	switch g.boardType {
	case ACE:
		return boardtemplates.ACE
	case CURIOSITY:
		return boardtemplates.CURIOSITY
	case GRAIL:
		return boardtemplates.GRAIL
		// poner el resto
	default:
		return ""
	}
}

func extraerEsquinas(laserPath []vector2_T) []vector2_T {
	// En caso de ser menor a 2 no hace falta extaer esquinas
	if len(laserPath) <= 2 {
		return laserPath
	}

	var l []vector2_T

	// Agregamos el primer punto
	l = append(l, laserPath[0])

	// Vamos agregando todas las esquinas
	for i := 1; i < len(laserPath)-1; i++ {
		anterior := laserPath[i-1]
		actual := laserPath[i]
		siguiente := laserPath[i+1]

		// Vector 1 (del punto anterior al actual)
		dx1 := actual.x - anterior.x
		dy1 := actual.y - anterior.y

		// Vector 2 (del punto actual al siguiente)
		dx2 := siguiente.x - actual.x
		dy2 := siguiente.y - actual.y

		// Se hace el producto cruzado para ver si es una esquina
		if (dx1*dy2)-(dy1*dx2) != 0 {
			l = append(l, actual)
		}
	}

	// Agregamos el último punto
	l = append(l, laserPath[len(laserPath)-1])

	return l
}

func (g *LaserChessGame) Run() {
	for message := range g.FromRoom {
		switch message.MsgType {
		case Move:
			switch g.turn {
			case g.redPlayer:
				if message.PlayerUid == g.redPlayer {
					fmt.Println("RED:", message.MsgContent)
					resul, laser, _, err := g.gameBoard.ProcessTurn(message.MsgContent, RED_TEAM)
					fmt.Println("ANSWER:", resul)
					g.ToRoom <- ResponseToRoom{
						Type:    Move,
						Content: resul,
						Laser:   fmt.Sprint(laser),
					}

					// Si el moviento es correcto se pasa el turno
					if err == nil {
						g.turn = g.bluePlayer
					}

				} else {
					// TODO: Esto habrá tratarlo mejor
					g.ToRoom <- ResponseToRoom{
						Type:    Error,
						Content: "no es tu turno",
						Laser:   "",
					}
				}
			case g.bluePlayer:
				if message.PlayerUid == g.bluePlayer {
					fmt.Println("BLUE:", message.MsgContent)
					resul, laser, _, err := g.gameBoard.ProcessTurn(message.MsgContent, BLUE_TEAM)
					fmt.Println("ANSWER:", resul)
					g.ToRoom <- ResponseToRoom{
						Type:    Move,
						Content: resul,
						Laser:   fmt.Sprint(extraerEsquinas(laser)),
					}

					// Si el moviento es correcto se pasa el turno
					if err == nil {
						g.turn = g.redPlayer
					}
				} else {
					// TODO: Esto habrá tratarlo mejor
					g.ToRoom <- ResponseToRoom{
						Type:    Error,
						Content: "no es tu turno",
						Laser:   "",
					}
				}
			}
		case GetState:
			// state := g.gameBoard.GetState()
			// g.ToRoom <- ResponseToRoom{MsgContent: state}
		case GetInitialState:
			initialState := g.getInitialState()
			g.ToRoom <- ResponseToRoom{
				Type:    InitialState,
				Content: initialState,
			}
		}
	}
}
