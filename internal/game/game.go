package game

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
	g.gameBoard = InitBoard(BoardType)
	go g.Run()
}

func (g *LaserChessGame) Run() {
	for message := range g.FromRoom {
		switch message.MsgType {
		case "Move":
			switch g.turn {
			case g.redPlayer:
				resul, _, _, _ := g.gameBoard.ProcessTurn(message.MsgContent, RED_TEAM)
				g.ToRoom <- ResponseToRoom{MsgContent: resul}
			case g.bluePlayer:
				resul, _, _, _ := g.gameBoard.ProcessTurn(message.MsgContent, BLUE_TEAM)
				g.ToRoom <- ResponseToRoom{MsgContent: resul}
			}
		case "GetState":
			// state := g.gameBoard.GetState()
			// g.ToRoom <- ResponseToRoom{MsgContent: state}
		}
	}
}
