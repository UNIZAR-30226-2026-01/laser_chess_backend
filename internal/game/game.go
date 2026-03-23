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

	piecesTakenByRed  []BoardPiece
	piecesTakenByBlue []BoardPiece
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
	// newBoard := InitBoard(ACE)
	// newGame := LaserChessGame{
	// 	redPlayer:  uidRedPlayer,
	// 	bluePlayer: uidBluePlayer,
	// 	turn:       uidRedPlayer,
	// 	gameBoard:  newBoard}
}

func (g *LaserChessGame) Run() {
	for {
		select {
		case message := <-g.FromRoom:

			switch message.MsgType {
			case "Move":
				if g.turn == g.redPlayer {
					resul, _, _, _ := g.gameBoard.ProcessTurn(message.MsgContent, RED_TEAM)
					g.ToRoom <- ResponseToRoom{MsgContent: resul}
				} else if g.turn == g.bluePlayer {
					resul, _, _, _ := g.gameBoard.ProcessTurn(message.MsgContent, BLUE_TEAM)
					g.ToRoom <- ResponseToRoom{MsgContent: resul}
				}

			case "GetState":
				// Funcion para coger el estado
			}

		}
	}
}
