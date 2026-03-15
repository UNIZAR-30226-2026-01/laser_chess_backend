package game

type RoomMsg struct {
	PlayerUid  int64
	MsgType    string
	MsgContent string
}

type ResponseToRoom struct {
	PlayerUid  int64
	MsgType    string
	MsgContent string
}

type LaserChessGame struct {
	redPlayer  int64
	bluePlayer int64

	FromRoom chan RoomMsg
	ToRoom   chan ResponseToRoom

	gameBoard Board

	piecesTakenByRed  []BoardPiece
	piecesTakenByBlue []BoardPiece
}
