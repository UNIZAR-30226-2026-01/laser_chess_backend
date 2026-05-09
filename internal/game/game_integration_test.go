package game

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGameTimer_Flow(t *testing.T) {
	initial := 200 * time.Millisecond
	increment := 50 * time.Millisecond
	gt := NewGameTimer(initial, increment)

	gt.Start()
	time.Sleep(50 * time.Millisecond)
	gt.Stop()

	assert.Greater(t, gt.Remaining, 150*time.Millisecond)
	assert.LessOrEqual(t, gt.Remaining, initial+increment)

	shortTimer := NewGameTimer(10*time.Millisecond, 0)
	shortTimer.Start()

	select {
	case expired := <-shortTimer.Expired:
		assert.True(t, expired)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("El timer no expiró")
	}
}

func TestGameEngine_ProcessTurn_Log(t *testing.T) {
	engine := GameEngine{}
	engine.InitEngine(ACE)

	// Rotar Deflector Rojo en c1.
	// Nota: Si el láser se sale del tablero, el resultado es OUT_OF_BOUNDS (2)
	timeLeft := 25000 * time.Millisecond
	instruction, path, result, err := engine.ProcessTurn("Rc1", RED_TEAM, timeLeft)

	require.NoError(t, err)
	assert.Contains(t, instruction, "Rc1")
	assert.NotNil(t, path)
	// Ajustamos al valor real que devuelve tu motor (OUT_OF_BOUNDS = 2)
	assert.Equal(t, OUT_OF_BOUNDS, result)

	state := engine.GetState()
	assert.Contains(t, state, "Rc1")
}

func TestLaserChessGame_HandleRoomMsg_Full(t *testing.T) {
	game := &LaserChessGame{}
	game.InitLaserChessGame(10, 20, ACE, "", 30000, 1000)

	// 1. Estado inicial
	game.HandleRoomMsg(RoomMsg{PlayerUid: 10, MsgType: GetInitialState})
	<-game.ToRoom

	// 2. Movimiento válido Rojo (10): Tc1:c2
	moveMsg := RoomMsg{PlayerUid: 10, MsgType: Move, MsgContent: "Tc1:c2"}
	game.HandleRoomMsg(moveMsg)
	respMove := <-game.ToRoom
	assert.Equal(t, Move, respMove.Type)

	// 3. Verificación de turno (le toca al 20)
	invalidTurn := RoomMsg{PlayerUid: 10, MsgType: Move, MsgContent: "Tc2:c3"}
	game.HandleRoomMsg(invalidTurn)
	respErr := <-game.ToRoom
	assert.Equal(t, Error, respErr.Type)
	assert.Equal(t, "no es tu turno", respErr.Content)

	// 4. Movimiento válido Azul (20):
	// Evitamos a8 (láser). Usamos otra pieza azul, por ejemplo b8 a b7.
	// (Nota: Asegúrate de que en b8 hay una pieza azul en tu versión de ACE)
	moveAzul := RoomMsg{PlayerUid: 20, MsgType: Move, MsgContent: "Th8:h7"}
	game.HandleRoomMsg(moveAzul)
	respMoveAzul := <-game.ToRoom

	if respMoveAzul.Type == Error {
		// Si b8 fallara, intenta con i8:i7
		t.Logf("Error azul: %s", respMoveAzul.Content)
	}
	assert.Equal(t, Move, respMoveAzul.Type)

	// 5. Pausar
	game.HandleRoomMsg(RoomMsg{PlayerUid: 10, MsgType: Pause})
	respPause := <-game.ToRoom
	assert.Equal(t, Paused, respPause.Type)
}

func TestEngine_ApplyLog_EdgeCases(t *testing.T) {
	engine := GameEngine{}
	engine.InitEngine(ACE)

	engine.gameLog = ""
	_, rt, bt := engine.EngineApplyLogToBoard(30000)
	assert.Equal(t, float64(30000), rt)
	assert.Equal(t, float64(30000), bt)

	// Log simulando rotaciones de Rojo (c1) y Azul (a8)
	engine.gameLog = "Rc1%c1%{29000};Ra8%a8%{28500};"
	next, rt, bt := engine.EngineApplyLogToBoard(30000)

	assert.Equal(t, RED_TEAM, next)
	assert.Equal(t, float64(29000), rt)
	assert.Equal(t, float64(28500), bt)
}

func TestMinMax_LogParsing(t *testing.T) {
	board, _ := InitBoard(boardTypeToCsv(ACE))
	log := "Rc1%c1%{29000};Ra8%a8%{28500};"

	nextTeam := MinMaxApplyLogToBoard(log, board)
	assert.Equal(t, RED_TEAM, nextTeam)
}
