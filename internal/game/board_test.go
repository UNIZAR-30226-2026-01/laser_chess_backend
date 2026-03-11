package game

import "testing"

func TestProbatTipoDeDato(t *testing.T){
	tablero := Board{}

	//Iniciar tablero

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
	tablero.cells[0][0] = &BoardPieceLaser{BLUE_TEAM, DOWN}
	tablero.cells[9][7] = &BoardPieceLaser{RED_TEAM, UP}

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

	// === TESTS AQUI ===
	tablero.movePiece(6,6,6,7)

	tablero.movePiece(6,7,6,6)
	
	tablero.rotatePiece(6,6,'L')

	tablero.print()
}