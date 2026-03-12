package game

import "testing"

func reiniciarTablero(tablero *Board) {

	for y := 0; y < 8; y++ {
		for x := 0; x < 10; x++ {
			tablero.cells[x][y] = &BoardPieceVacant{NONE}
		}
	}

	for y := 0; y < 8; y++ {
		tablero.cells[0][y] = &BoardPieceVacant{BLUE_TEAM}
		tablero.cells[9][y] = &BoardPieceVacant{RED_TEAM}
	}

	tablero.cells[1][0] = &BoardPieceVacant{RED_TEAM}
	tablero.cells[1][7] = &BoardPieceVacant{RED_TEAM}
	tablero.cells[8][0] = &BoardPieceVacant{BLUE_TEAM}
	tablero.cells[8][7] = &BoardPieceVacant{BLUE_TEAM}

	//Pruebas sobre tablero "ACE"

	//Poner reyes
	tablero.cells[5][0] = &BoardPieceKing{BLUE_TEAM, NONE}
	tablero.cells[4][7] = &BoardPieceKing{RED_TEAM, NONE}

	//Poner laseres
	tablero.blueTeamLaser = &BoardPieceLaser{BLUE_TEAM, DOWN}
	tablero.cells[0][0] = tablero.blueTeamLaser
	tablero.redTeamLaser = &BoardPieceLaser{RED_TEAM, UP}
	tablero.cells[9][7] = tablero.redTeamLaser

	//Poner escudos
	tablero.cells[4][0] = &BoardPieceShield{BLUE_TEAM, NONE, DOWN}
	tablero.cells[6][0] = &BoardPieceShield{BLUE_TEAM, NONE, DOWN}
	tablero.cells[3][7] = &BoardPieceShield{RED_TEAM, NONE, UP}
	tablero.cells[5][7] = &BoardPieceShield{RED_TEAM, NONE, UP}

	//poner switches
	tablero.cells[4][3] = &BoardPieceSwitch{BLUE_TEAM, NONE, DOWN}
	tablero.cells[5][3] = &BoardPieceSwitch{BLUE_TEAM, NONE, LEFT}
	tablero.cells[4][4] = &BoardPieceSwitch{RED_TEAM, NONE, LEFT}
	tablero.cells[5][4] = &BoardPieceSwitch{RED_TEAM, NONE, DOWN}

	//poner deflectores
	tablero.cells[7][0] = &BoardPieceDeflector{BLUE_TEAM, NONE, LEFT}
	tablero.cells[7][3] = &BoardPieceDeflector{BLUE_TEAM, NONE, LEFT}
	tablero.cells[7][4] = &BoardPieceDeflector{BLUE_TEAM, NONE, DOWN}
	tablero.cells[6][5] = &BoardPieceDeflector{BLUE_TEAM, NONE, LEFT}
	tablero.cells[0][3] = &BoardPieceDeflector{BLUE_TEAM, BLUE_TEAM, DOWN}
	tablero.cells[0][4] = &BoardPieceDeflector{BLUE_TEAM, BLUE_TEAM, LEFT}
	tablero.cells[2][1] = &BoardPieceDeflector{BLUE_TEAM, NONE, UP}
	tablero.cells[2][7] = &BoardPieceDeflector{RED_TEAM, NONE, RIGHT}
	tablero.cells[2][4] = &BoardPieceDeflector{RED_TEAM, NONE, RIGHT}
	tablero.cells[2][3] = &BoardPieceDeflector{RED_TEAM, NONE, UP}
	tablero.cells[3][2] = &BoardPieceDeflector{RED_TEAM, NONE, RIGHT}
	tablero.cells[9][4] = &BoardPieceDeflector{RED_TEAM, RED_TEAM, UP}
	tablero.cells[9][3] = &BoardPieceDeflector{RED_TEAM, RED_TEAM, RIGHT}
	tablero.cells[7][6] = &BoardPieceDeflector{RED_TEAM, NONE, DOWN}
}

func TestProbatTipoDeDato(t *testing.T) {
	tablero := Board{}

	//Iniciar tablero
	reiniciarTablero(&tablero)

	tablero.print()

	// === Rey === //
	t.Log("MOVIMIENTO")

	t.Log("== Rey ==")

	if tablero.ProcessTurn("Te8:e6") != false {
		t.Errorf("X - Se ha aceptado un movimiento fuera del alcance de la pieza")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	if tablero.ProcessTurn("Te8:e9") != false {
		t.Errorf("X - Se ha aceptado un movimiento fuera del tablero")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	if tablero.ProcessTurn("Te8:d8") != false {
		t.Errorf("X - Se ha aceptado un movimiento a una casilla ocupada")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	tablero.cells[4][6] = &BoardPieceVacant{BLUE_TEAM}
	if tablero.ProcessTurn("Te8:e7") != false {
		t.Errorf("X - Se ha aceptado un movimiento a una casilla del equipo opuesto")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}
	tablero.cells[4][6] = &BoardPieceVacant{NONE}

	if tablero.ProcessTurn("Te8:e7") != true {
		t.Errorf("X - Se ha rechazado un movimiento válido")
	} else {
		reiniciarTablero(&tablero)
		t.Log("OK")
	}

	tablero.cells[4][6] = &BoardPieceVacant{RED_TEAM}
	if tablero.ProcessTurn("Te8:e7") != true {
		t.Errorf("X - Se ha rechazado un movimiento válido")
	} else {
		reiniciarTablero(&tablero)
		t.Log("OK")
	}
	tablero.cells[4][6] = &BoardPieceVacant{NONE}

	t.Log("== Escudo ==")

	if tablero.ProcessTurn("Td8:d6") != false {
		t.Errorf("X - Se ha aceptado un movimiento fuera del alcance de la pieza")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	if tablero.ProcessTurn("Td8:d9") != false {
		t.Errorf("X - Se ha aceptado un movimiento fuera del tablero")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	if tablero.ProcessTurn("Tf8:e8") != false {
		t.Errorf("X - Se ha aceptado un movimiento a una casilla ocupada")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	tablero.cells[3][6] = &BoardPieceVacant{BLUE_TEAM}
	if tablero.ProcessTurn("Td8:d7") != false {
		t.Errorf("X - Se ha aceptado un movimiento a una casilla del equipo opuesto")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}
	tablero.cells[3][6] = &BoardPieceVacant{NONE}

	if tablero.ProcessTurn("Td8:c7") != true {
		t.Errorf("X - Se ha rechazado un movimiento válido")
	} else {
		reiniciarTablero(&tablero)
		t.Log("OK")
	}

	tablero.cells[3][6] = &BoardPieceVacant{RED_TEAM}
	if tablero.ProcessTurn("Td8:c7") != true {
		t.Errorf("X - Se ha rechazado un movimiento válido")
	} else {
		reiniciarTablero(&tablero)
		t.Log("OK")
	}
	tablero.cells[3][6] = &BoardPieceVacant{NONE}

	t.Log("== Deflector ==")

	if tablero.ProcessTurn("Tc8:c6") != false {
		t.Errorf("X - Se ha aceptado un movimiento fuera del alcance de la pieza")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	if tablero.ProcessTurn("Tc8:c9") != false {
		t.Errorf("X - Se ha aceptado un movimiento fuera del tablero")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	if tablero.ProcessTurn("Th4:h5") != false {
		t.Errorf("X - Se ha aceptado un movimiento a una casilla ocupada")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	if tablero.ProcessTurn("Tc2:b1") != false {
		t.Errorf("X - Se ha aceptado un movimiento a una casilla del equipo opuesto")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	if tablero.ProcessTurn("Ta4:b3") != true {
		t.Errorf("X - Se ha rechazado un movimiento válido")
	} else {
		reiniciarTablero(&tablero)
		t.Log("OK")
	}

	tablero.cells[3][6] = &BoardPieceVacant{RED_TEAM}
	if tablero.ProcessTurn("Ta4:a3") != true {
		t.Errorf("X - Se ha rechazado un movimiento válido")
	} else {
		reiniciarTablero(&tablero)
		t.Log("OK")
	}
	tablero.cells[3][6] = &BoardPieceVacant{NONE}

	t.Log("== Laser ==")

	if tablero.ProcessTurn("Ta1:a2") != false {
		t.Errorf("X - Esta pieza siempre debería devolver false, (rarete)")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	t.Log("== Switch ==")

	if tablero.ProcessTurn("Te5:g7") != false {
		t.Errorf("X - Se ha aceptado un movimiento fuera del alcance de la pieza")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	tablero.cells[6][7] = &BoardPieceSwitch{BLUE_TEAM, NONE, DOWN}
	if tablero.ProcessTurn("Tg8:g7") != false {
		t.Errorf("X - Se ha aceptado un movimiento fuera del tablero")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}
	tablero.cells[6][7] = &BoardPieceVacant{NONE}

	tablero.cells[6][7] = &BoardPieceSwitch{BLUE_TEAM, NONE, DOWN}
	if tablero.ProcessTurn("Tg8:h7") != false {
		t.Errorf("X - Se ha aceptado un movimiento a una casilla ocupada (permutación no posible)")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}
	tablero.cells[6][7] = &BoardPieceVacant{NONE}

	tablero.cells[7][7] = &BoardPieceSwitch{BLUE_TEAM, NONE, DOWN}
	if tablero.ProcessTurn("Th8:i8") != false {
		t.Errorf("X - Se ha aceptado un movimiento a una casilla del equipo opuesto")
		reiniciarTablero(&tablero)
	} else {
		t.Log("OK")
	}

	if tablero.ProcessTurn("Te5:f4") != true {
		t.Errorf("X - Se ha rechazado un movimiento válido (posible error permutación)")
	} else {
		reiniciarTablero(&tablero)
		t.Log("OK")
	}

	tablero.cells[6][7] = &BoardPieceSwitch{BLUE_TEAM, NONE, DOWN}
	if tablero.ProcessTurn("Tg8:f8") != true {
		t.Errorf("X - Se ha rechazado un movimiento válido (posible error permutación)")

	} else {
		t.Log("OK")
	}
	tablero.cells[6][7] = &BoardPieceVacant{NONE}

	// Test del recorrido del laser
	t.Log("Solo dispara el laser y debería finalizar por out of bounds")
	positions, terminationReason := tablero.blueTeamLaser.shootLaser(0, 0, &tablero)
	t.Log(positions)
	if terminationReason != OUT_OF_BOUNDS {
		t.Errorf("X - No ha terminado por la razón correcta, se esperaba OUT_OF_BOUNDS, recibido %d", terminationReason)
	} else {
		t.Log("OK")
	}

}
