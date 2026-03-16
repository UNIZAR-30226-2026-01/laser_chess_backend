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
func (g *LaserChessGame) InitLaserChessGame(uidRedPlayer int64, uidBluePlayer int64) {
	g.redPlayer = uidRedPlayer
	g.bluePlayer = uidBluePlayer
	g.turn = uidRedPlayer
	g.gameBoard = InitBoard(ACE)
	// newBoard := InitBoard(ACE)
	// newGame := LaserChessGame{
	// 	redPlayer:  uidRedPlayer,
	// 	bluePlayer: uidBluePlayer,
	// 	turn:       uidRedPlayer,
	// 	gameBoard:  newBoard}
}

func (g *LaserChessGame) Run() {

}
