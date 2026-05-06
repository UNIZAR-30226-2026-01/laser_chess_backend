package game

type GameMessageType string

const (
	Move            GameMessageType = "Move"
	GetState        GameMessageType = "GetState"
	GetInitialState GameMessageType = "GetInitialState"
	Pause           GameMessageType = "Pause"

	State        GameMessageType = "State"
	MatchStart   GameMessageType = "MatchStart"
	InitialState GameMessageType = "InitialState"
	PauseRequest GameMessageType = "PauseRequest"
	PauseReject  GameMessageType = "PauseReject"
	Paused       GameMessageType = "Paused"
	End          GameMessageType = "End"
	Error        GameMessageType = "Error"
	EOC          GameMessageType = "EOC"

	Disconnection GameMessageType = "Disconnection"
	Reconnection  GameMessageType = "Reconnection"

	Rewards   GameMessageType = "Rewards"
	EloUpdate GameMessageType = "EloUpdate"
)
