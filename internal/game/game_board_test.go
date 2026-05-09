package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoard_MovePiece(t *testing.T) {
	b := createEmptyBoard()

	// Lo colocamos en una zona neutral (centro del tablero) para evitar las filas de colores
	b.cells[5][5] = &BoardPieceDeflector{team: RED_TEAM, pointing: DOWN}

	// Movemos una casilla hacia abajo (x=5, y=6)
	err := b.movePiece(5, 5, 5, 6, RED_TEAM)
	assert.NoError(t, err, "El movimiento debería ser válido")

	// Verificamos que la pieza ha llegado a su destino
	destPiece, isDeflector := b.cells[5][6].(*BoardPieceDeflector)

	// Envolvemos el True en un if para evitar el panic si destPiece es nil
	if assert.True(t, isDeflector, "El deflector debería estar en la nueva casilla (5,6)") {
		assert.Equal(t, DOWN, destPiece.pointing, "El deflector no debería haber cambiado su orientación")
	}

	// Verificamos que la casilla original ahora está vacía
	_, isVacant := b.cells[5][5].(*BoardPieceVacant)
	assert.True(t, isVacant, "La casilla original (5,5) debería haberse quedado vacía (BoardPieceVacant)")
}

func TestBoard_RotatePiece(t *testing.T) {
	b := createEmptyBoard()

	// Colocamos un escudo azul en la esquina inferior derecha (x=9, y=7) que sería 'j1'
	shield := &BoardPieceShield{team: BLUE_TEAM, pointing: UP}
	b.cells[9][7] = shield

	// Rotamos hacia la derecha ('R')
	err := b.rotatePiece(9, 7, 'R', BLUE_TEAM)
	assert.NoError(t, err, "La rotación debería ser válida")

	// Verificamos que la rotación se ha aplicado a la pieza en el tablero
	assert.Equal(t, RIGHT, shield.pointing, "El escudo debería apuntar a la derecha (RIGHT = 3) tras rotar 'R'")
}

func TestBoard_MovePiece_FailsOnInvalidRule(t *testing.T) {
	b := createEmptyBoard()

	// Colocamos un Rey rojo
	b.cells[5][5] = &BoardPieceKing{team: RED_TEAM}

	// Intentamos moverlo con el equipo azul (debería fallar por la lógica interna del rey)
	err := b.movePiece(5, 5, 5, 6, BLUE_TEAM)
	assert.Error(t, err, "No se debería poder mover una pieza del rival")

	// La pieza no debe haberse movido
	_, isKing := b.cells[5][5].(*BoardPieceKing)
	assert.True(t, isKing, "El Rey debe seguir en su casilla original si el movimiento falla")
}

func TestBoard_ProcessTurn_MoveAndShoot(t *testing.T) {
	b := createEmptyBoard()

	b.blueTeamLaser = &BoardPieceLaser{team: BLUE_TEAM, pointing: DOWN}
	b.redTeamLaser = &BoardPieceLaser{team: RED_TEAM, pointing: UP}
	b.cells[0][0] = b.blueTeamLaser
	b.cells[9][7] = b.redTeamLaser

	b.cells[1][1] = &BoardPieceDeflector{team: BLUE_TEAM, pointing: DOWN}

	// USAMOS EL FORMATO BNF: "T" y ":" => Tb7:b6
	_, laserPath, laserResult, err := b.ProcessTurn("Tb7:b6", BLUE_TEAM)

	assert.NoError(t, err, "El comando ProcessTurn debería ejecutarse sin errores")

	_, originalIsVacant := b.cells[1][1].(*BoardPieceVacant)
	assert.True(t, originalIsVacant, "La casilla origen (b7) debe quedar vacía")

	_, destIsDeflector := b.cells[1][2].(*BoardPieceDeflector)
	assert.True(t, destIsDeflector, "El deflector debe estar en su nuevo destino (b6)")

	assert.NotEmpty(t, laserPath, "El path del láser no debe estar vacío")
	assert.Equal(t, OUT_OF_BOUNDS, laserResult, "El láser debería haberse salido del tablero")
}

func TestBoard_ProcessTurn_RotateAndShoot(t *testing.T) {
	b := createEmptyBoard()

	b.blueTeamLaser = &BoardPieceLaser{team: BLUE_TEAM, pointing: DOWN}
	b.redTeamLaser = &BoardPieceLaser{team: RED_TEAM, pointing: UP}
	b.cells[0][0] = b.blueTeamLaser
	b.cells[9][7] = b.redTeamLaser

	b.cells[4][4] = &BoardPieceKing{team: RED_TEAM}

	// USAMOS EL FORMATO BNF: Re4
	_, _, _, err := b.ProcessTurn("Re4", RED_TEAM)

	assert.Error(t, err, "El ProcessTurn debe propagar el error si el movimiento es ilegal")
	assert.ErrorContains(t, err, "rey no puede rotar")
}
