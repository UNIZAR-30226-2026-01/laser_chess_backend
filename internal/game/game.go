package game

import "fmt"

type RoomMsg struct {
	PlayerUid  int64
	MsgType    string
	MsgContent string
}

type ResponseToRoom struct {
	MsgContent string
}

type LaserChessGame struct {
	redPlayer  int64
	bluePlayer int64

	turn int64

	FromRoom chan RoomMsg
	ToRoom   chan ResponseToRoom

	gameBoard Board
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
	g.gameBoard.InitBoard("boardTemplates/ace.csv")

	g.FromRoom = make(chan RoomMsg)
	g.ToRoom = make(chan ResponseToRoom)

	go g.Run()

	fmt.Println("Game inicializado")
}

func (g *LaserChessGame) Run() {
	for message := range g.FromRoom {
		switch message.MsgType {
		case "Move":
			switch g.turn {
			case g.redPlayer:
				if message.PlayerUid == g.redPlayer {
					fmt.Println("RED:", message.MsgContent)
					resul, _, _, err := g.gameBoard.ProcessTurn(message.MsgContent, RED_TEAM)
					fmt.Println("ANSWER:", resul)
					g.ToRoom <- ResponseToRoom{MsgContent: resul}

					// Si el moviento es correcto se pasa el turno
					if err == nil {
						g.turn = g.bluePlayer
					}

				} else {
					// TODO: Esto habrá tratarlo mejor
					g.ToRoom <- ResponseToRoom{MsgContent: "no es tu turno"}
				}
			case g.bluePlayer:
				if message.PlayerUid == g.bluePlayer {
					fmt.Println("BLUE:", message.MsgContent)
					resul, _, _, err := g.gameBoard.ProcessTurn(message.MsgContent, BLUE_TEAM)
					fmt.Println("ANSWER:", resul)
					g.ToRoom <- ResponseToRoom{MsgContent: resul}

					// Si el moviento es correcto se pasa el turno
					if err == nil {
						g.turn = g.redPlayer
					}
				} else {
					// TODO: Esto habrá tratarlo mejor
					g.ToRoom <- ResponseToRoom{MsgContent: "no es tu turno"}
				}
			}
		case "GetState":
			// state := g.gameBoard.GetState()
			// g.ToRoom <- ResponseToRoom{MsgContent: state}
		}
	}
}
