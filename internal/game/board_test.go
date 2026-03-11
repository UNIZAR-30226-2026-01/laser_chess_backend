package game

import "testing"

func TestProbatTipoDeDato(t *testing.T){
	tablero := Board{}

	//Iniciar tablero

	for y := 0; y < 8; y++ {
		for x := 0; x < 10; x++ {
			tablero.cells[x][y] = &BoardPieceVacant{'n'}
		}
	}

	for y := 0; y < 8; y++ {
		tablero.cells[0][y] = &BoardPieceVacant{'b'}
		tablero.cells[9][y] = &BoardPieceVacant{'r'}
	}

	tablero.cells[1][0] = &BoardPieceVacant{'r'}
	tablero.cells[1][7] = &BoardPieceVacant{'r'}
	tablero.cells[8][0] = &BoardPieceVacant{'b'}
	tablero.cells[8][7] = &BoardPieceVacant{'b'}

	//Poner un rey en una casilla
	tablero.cells[6][6] = &BoardPieceKing{'r', 'n'}
	
	
	tablero.movePiece(6,6,6,7)
	tablero.movePiece(6,7,6,6)

	tablero.rotatePiece(6,6,'L')

	tablero.print()
}