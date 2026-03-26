package game

type GameMessageType string

const (
	Move            GameMessageType = "Move"
	GetState        GameMessageType = "GetState"
	GetInitialState GameMessageType = "GetInitialState"
	Pause           GameMessageType = "Pause"

	State        GameMessageType = "State"
	InitialState GameMessageType = "InitialState"
	PauseRequest GameMessageType = "PauseRequest"
	Paused       GameMessageType = "Paused"
	End          GameMessageType = "End"
	Error        GameMessageType = "Error"
)
