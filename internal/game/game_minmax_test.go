package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	boardtemplates "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game/boardTemplates"
)

func TestInitBoard_ACE(t *testing.T) {
	// 1. Probamos que el parser CSV funciona bien usando el template real exportado
	board, err := InitBoard(boardtemplates.ACE)
	require.NoError(t, err, "El CSV del tablero ACE debería parsearse sin errores")
	require.NotNil(t, board, "El tablero no debe ser nulo tras inicializarse")

	// 2. Verificamos casillas clave según tu documento

	// Casilla a8 (x=0, y=0) -> Laser Azul Apuntando Abajo (LAD)
	laserAzul, ok := board.cells[0][0].(*BoardPieceLaser)
	assert.True(t, ok, "Debería haber un Láser en a8")
	if ok {
		assert.Equal(t, BLUE_TEAM, laserAzul.team)
		assert.Equal(t, DOWN, laserAzul.pointing)
	}

	// Casilla f8 (x=5, y=0) -> Rey Azul (KA)
	reyAzul, ok := board.cells[5][0].(*BoardPieceKing)
	assert.True(t, ok, "Debería haber un Rey en f8")
	if ok {
		assert.Equal(t, BLUE_TEAM, reyAzul.team)
	}

	// Casilla d1 (x=3, y=7) -> Escudo Rojo apuntando Arriba (ERU)
	escudoRojo, ok := board.cells[3][7].(*BoardPieceShield)
	assert.True(t, ok, "Debería haber un Escudo en d1")
	if ok {
		assert.Equal(t, RED_TEAM, escudoRojo.team)
		assert.Equal(t, UP, escudoRojo.pointing)
	}
}

func TestMinMax_GetBestMove(t *testing.T) {
	// Vamos a probar que tu IA es capaz de evaluar el tablero y devolver un movimiento
	// Usamos profundidad 1 (lvl 1) para que el test sea súper rápido y no ralentice el CI/CD

	// ACE = 0. Le pasamos un log vacío para simular el turno 1.
	bestMove := GetBestMove(ACE, "", 1)

	// La IA debería devolver un movimiento válido (ej: "Ta2:a3" o "Ra1")
	assert.NotEmpty(t, bestMove, "El bot MinMax debería devolver un movimiento")
	assert.GreaterOrEqual(t, len(bestMove), 3, "El movimiento devuelto debe tener sentido (ej: 'Ra1' mínimo)")
}

func TestEngineApplyLogToBoard(t *testing.T) {
	// 1. Usamos el tablero ACE real para que los movimientos sean coherentes
	board, _ := InitBoard(boardtemplates.ACE)

	engine := GameEngine{
		gameBoard: board,
		boardType: ACE,
		// Log con 2 movimientos siguiendo tu BNF:
		// Turno 1 (Rojo): Rotar pieza en c1 (DRR) -> tiempo queda en 25000ms
		// Turno 2 (Azul): Rotar pieza en a8 (LAD) -> tiempo queda en 24000ms
		gameLog: "Rc1%c1%{25000};Ra8%a8%{24000};",
	}

	// 2. Ejecutamos la lógica (timeBase es el tiempo inicial, ej: 30 segundos)
	nextTeam, redTime, blueTime := engine.EngineApplyLogToBoard(30000)

	// 3. Verificaciones
	// Según engine.go, el primer movimiento se asigna al equipo que empieza (RED_TEAM)
	// Log 1: Rojo mueve -> redTime = 25000, nextTeam = BLUE
	// Log 2: Azul mueve -> blueTime = 24000, nextTeam = RED

	assert.Equal(t, RED_TEAM, nextTeam, "Después de 2 movimientos, vuelve a ser turno de Rojo")

	// Verificamos que los tiempos coincidan con el último valor registrado para cada equipo
	assert.Equal(t, float64(25000), redTime, "El tiempo de Rojo debe ser el del primer movimiento")
	assert.Equal(t, float64(24000), blueTime, "El tiempo de Azul debe ser el del segundo movimiento")
}
