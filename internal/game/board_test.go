package game

import "testing"

func TestProbatTipoDeDato(t *testing.T){
	tablero := Board{}

	for x := 0; x < 3; x++ {
		tablero.cells[x][0] = &BoardPieceKing{
			team: 'r',
		}
		tablero.cells[x][1] = &BoardPieceShield{
			team: 'r',
		}
		tablero.cells[x][2] = &BoardPieceShield{
			team: 'b',
		}
	}

	//SE PRUEBA LA LLAMADA A LO QUE DEBERÍA SER UN REY (TRUE ESPERADO)
	tablero.movePiece(0, 0, 1, 1)

	//SE PRUEBA LA LLAMADA A LO QUE DEBERÍA SER UN REY (FALSE ESPERADO)
	tablero.rotatePiece(0, 1, 'L')

	//SE PRUEBA LA LLAMADA A LO QUE DEBERÍA SER UN ESCUDO (TRUE ESPERADO)
	tablero.movePiece(0, 1, 1, 1)

	//SE PRUEBA LA LLAMADA A LO QUE DEBERÍA SER UN ESCUDO (TRUE ESPERADO)
	tablero.rotatePiece(0, 1, 'L')
}