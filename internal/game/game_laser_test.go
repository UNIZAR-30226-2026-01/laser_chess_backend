package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShootLaser_OutOfBounds(t *testing.T) {
	b := createEmptyBoard()
	// Un láser en (0,0) apuntando hacia ARRIBA (UP=2).
	// El vector UP es {0, -1}, por lo que el siguiente paso es (0, -1), saliéndose del tablero.
	laser := &BoardPieceLaser{team: RED_TEAM, pointing: UP}

	path, result := laser.shootLaser(0, 0, b)

	assert.Equal(t, OUT_OF_BOUNDS, result, "El láser debería salirse del tablero")
	assert.NotEmpty(t, path, "El path debería contener al menos la posición de disparo")
}

func TestShootLaser_HitKing(t *testing.T) {
	b := createEmptyBoard()
	laser := &BoardPieceLaser{team: RED_TEAM, pointing: DOWN}

	// Ponemos un Rey azul justo debajo del láser (y=1)
	b.cells[0][1] = &BoardPieceKing{team: BLUE_TEAM}

	path, result := laser.shootLaser(0, 0, b)

	assert.Equal(t, HIT_BLUE_KING, result, "El láser debería destruir al rey azul")
	assert.Contains(t, path, vector2_T{0, 1}, "El path debe registrar la casilla del impacto")
}

func TestShootLaser_Shield(t *testing.T) {
	b := createEmptyBoard()
	laser := &BoardPieceLaser{team: RED_TEAM, pointing: DOWN}

	t.Run("Escudo bloquea de frente", func(t *testing.T) {
		// El láser va hacia ABAJO (DOWN=0). El escudo mira hacia ARRIBA (UP=2).
		// Se encuentran cara a cara.
		b.cells[0][1] = &BoardPieceShield{team: BLUE_TEAM, pointing: UP}

		_, result := laser.shootLaser(0, 0, b)
		assert.Equal(t, STOP, result, "El escudo debería detener el láser de frente sin morir")
	})

	t.Run("Escudo muere por la espalda", func(t *testing.T) {
		// El escudo mira hacia ABAJO (DOWN=0), igual que el láser. Le da por la espalda.
		b.cells[0][1] = &BoardPieceShield{team: BLUE_TEAM, pointing: DOWN}

		_, result := laser.shootLaser(0, 0, b)
		assert.Equal(t, HIT, result, "El láser debería destruir el escudo si le da por detrás o de lado")
	})
}

func TestShootLaser_DeflectorAndSwitch(t *testing.T) {
	b := createEmptyBoard()
	laser := &BoardPieceLaser{team: RED_TEAM, pointing: DOWN}

	t.Run("Rebote en Deflector", func(t *testing.T) {
		// Ponemos un deflector en (0,1).
		// Dependiendo de tu implementación exacta, apuntar a un lado u otro causará el rebote.
		// En cualquier caso, si le da en el espejo, el resultado NO debe ser HIT ni STOP,
		// debe rebotar y acabar saliéndose del tablero (OUT_OF_BOUNDS) o chocando con otra cosa.
		b.cells[0][1] = &BoardPieceDeflector{team: BLUE_TEAM, pointing: RIGHT}

		_, result := laser.shootLaser(0, 0, b)
		assert.NotEqual(t, HIT, result, "El deflector no debería morir si se le da en el espejo")
		assert.NotEqual(t, STOP, result, "El deflector no debe detener el rayo, sino desviarlo")
	})

	t.Run("Switch siempre desvía o muere", func(t *testing.T) {
		// El Switch tiene espejos en ambas caras diagonales, así que un rayo directo
		b.cells[0][1] = &BoardPieceSwitch{team: BLUE_TEAM, pointing: DOWN}

		_, result := laser.shootLaser(0, 0, b)
		assert.NotEqual(t, STOP, result, "El Switch nunca detiene el láser en seco")
	})
}
