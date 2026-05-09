package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper: Crea un tablero totalmente vacío para los tests
func createEmptyBoard() *Board {
	b := &Board{}
	for x := 0; x < XDIM; x++ {
		for y := 0; y < YDIM; y++ {
			b.cells[x][y] = &BoardPieceVacant{}
		}
	}
	return b
}

func TestVacantPiece(t *testing.T) {
	b := createEmptyBoard()
	vacant := &BoardPieceVacant{}

	errMove := vacant.canMoveTo(5, 5, b, RED_TEAM)
	assert.Error(t, errMove, "Una casilla vacía no debería poder moverse")

	errRot := vacant.canRotate('L', RED_TEAM)
	assert.Error(t, errRot, "Una casilla vacía no debería poder rotar")
}

func TestKingPiece(t *testing.T) {
	b := createEmptyBoard()
	king := &BoardPieceKing{team: RED_TEAM}

	// 1. Movimiento válido a casilla neutral (x=5 no pertenece a nadie)
	err := king.canMoveTo(5, 5, b, RED_TEAM)
	assert.NoError(t, err, "El rey debería poder moverse a una casilla vacía neutral")

	// 2. Intentar moverlo con el equipo contrario
	err = king.canMoveTo(5, 6, b, BLUE_TEAM)
	assert.ErrorContains(t, err, "ficha del equipo opuesto")

	// 3. Intentar pisar una casilla del equipo azul (x=0 es la fila del BLUE_TEAM)
	err = king.canMoveTo(0, 5, b, RED_TEAM)
	assert.ErrorContains(t, err, "casilla destino del equipo opuesto")

	// 4. Intentar rotar (El rey no puede)
	err = king.canRotate('L', RED_TEAM)
	assert.ErrorContains(t, err, "rey no puede rotar")
}

func TestDeflectorPiece(t *testing.T) {
	b := createEmptyBoard()
	// Deflector apuntando hacia ARRIBA (UP = 2)
	def := &BoardPieceDeflector{team: BLUE_TEAM, pointing: UP}

	// 1. Rotación hacia la derecha (Clockwise) -> Debería acabar apuntando a RIGHT (3)
	err := def.canRotate('R', BLUE_TEAM)
	assert.NoError(t, err)
	assert.Equal(t, RIGHT, def.pointing)

	// 2. Rotación hacia la izquierda (Counterclockwise) -> Debería volver a UP (2)
	err = def.canRotate('L', BLUE_TEAM)
	assert.NoError(t, err)
	assert.Equal(t, UP, def.pointing)

	// 3. Intentar rotar una pieza del rival
	err = def.canRotate('R', RED_TEAM)
	assert.ErrorContains(t, err, "ficha del equipo opuesto")

	// 4. Casilla destino ocupada
	b.cells[5][5] = &BoardPieceKing{team: BLUE_TEAM} // Ponemos un rey en medio
	err = def.canMoveTo(5, 5, b, BLUE_TEAM)
	assert.ErrorContains(t, err, "casilla destino ocupada", "El deflector no puede pisar otras piezas")
}

func TestSwitchPiece(t *testing.T) {
	b := createEmptyBoard()
	sw := &BoardPieceSwitch{team: RED_TEAM, pointing: DOWN}

	// El Switch es la única pieza con reglas de permutación (puede pisar otras piezas para intercambiarse)

	t.Run("Mover a casilla vacía", func(t *testing.T) {
		err := sw.canMoveTo(5, 5, b, RED_TEAM)
		assert.NoError(t, err)
	})

	t.Run("Permutar con Shield", func(t *testing.T) {
		b.cells[5][6] = &BoardPieceShield{team: BLUE_TEAM, pointing: UP}
		err := sw.canMoveTo(5, 6, b, RED_TEAM)
		assert.NoError(t, err, "El Switch debería poder intercambiarse con un Shield")
	})

	t.Run("Permutar con Deflector", func(t *testing.T) {
		b.cells[5][7] = &BoardPieceDeflector{team: RED_TEAM, pointing: LEFT}
		err := sw.canMoveTo(5, 7, b, RED_TEAM)
		assert.NoError(t, err, "El Switch debería poder intercambiarse con un Deflector")
	})

	t.Run("Falla al intentar permutar con un King", func(t *testing.T) {
		b.cells[6][6] = &BoardPieceKing{team: BLUE_TEAM}
		err := sw.canMoveTo(6, 6, b, RED_TEAM)
		assert.ErrorContains(t, err, "casilla destino ocupada", "El Switch NO puede intercambiarse con un Rey")
	})
}

func TestTeamTileLogic(t *testing.T) {
	// Probamos la función independiente getTeamTile(x, y) de game_board.go
	assert.Equal(t, BLUE_TEAM, getTeamTile(0, 5), "Fila 0 entera debe ser AZUL")
	assert.Equal(t, RED_TEAM, getTeamTile(9, 5), "Fila 9 entera debe ser ROJA")

	assert.Equal(t, RED_TEAM, getTeamTile(1, 0), "Casilla especial (1,0) debe ser ROJA")
	assert.Equal(t, RED_TEAM, getTeamTile(1, 7), "Casilla especial (1,7) debe ser ROJA")

	assert.Equal(t, BLUE_TEAM, getTeamTile(8, 0), "Casilla especial (8,0) debe ser AZUL")
	assert.Equal(t, BLUE_TEAM, getTeamTile(8, 7), "Casilla especial (8,7) debe ser AZUL")

	assert.Equal(t, NONE, getTeamTile(5, 5), "Centro del tablero debe ser neutral")
}
