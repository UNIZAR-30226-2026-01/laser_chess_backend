package game

// ===================================== //
//	BOARD								 //
// ===================================== //

import (
	"fmt"
	"strings"
)

type Board struct {
	cells         [XDIM][YDIM]BoardPiece
	blueTeamLaser *BoardPieceLaser
	redTeamLaser  *BoardPieceLaser
}

// devuelve true si la posicion que se le pasa esta dentro del tablero
func (b *Board) isInbound(x int, y int) bool {
	return (0 <= x && x < XDIM) && (0 <= y && y < YDIM)
}

func (b *Board) movePiece(x_from int, y_from int, x_to int, y_to int) bool {
	return b.cells[x_from][y_from].canMoveTo(x_to, y_to)
}

func (b *Board) rotatePiece(x_at int, y_at int, rot rune) bool {
	return b.cells[x_at][y_at].canRotate(rot)
}

// ---Depuración---//
func (b *Board) printlaser(laser []vector2_T) {
	for y := 0; y < YDIM; y++ {
		fmt.Printf("%d | ", y+1) // numero
		for x := 0; x < XDIM; x++ {

			cell := b.cells[x][y].VisualRep()

			for i := 0; i < len(laser); i++ {
				if laser[i].x == x && laser[i].y == y {
					cell = "\033[47m" + cell + "\033[0m"
				}
			}

			fmt.Print(cell)
			fmt.Printf(" ")
		}
		fmt.Printf("\n")
	}
	fmt.Println("  +---------------------") // letra
	fmt.Println("    A B C D E F G H I J ") // letra
}

// ---Depuración---//
func (b *Board) print() {
	for y := 0; y < YDIM; y++ {
		fmt.Printf("%d | ", y+1) // numero
		for x := 0; x < XDIM; x++ {

			cell := b.cells[x][y].VisualRep()

			fmt.Print(cell)
			fmt.Printf(" ")
		}
		fmt.Printf("\n")
	}
	fmt.Println("  +---------------------") // letra
	fmt.Println("    A B C D E F G H I J ") // letra
}

// --- INTERFAZ DE COMUNICACIÓN CON EL MÓDULO --- //

func (b *Board) ProcessTurn(instruction string) bool {

	//La versión de golang del string stream
	reader := strings.NewReader(instruction)

	var inst rune //Instrucción
	_, err := fmt.Fscanf(reader, "%c", &inst)
	if err != nil { /*TODO*/
	}

	switch inst {
	case 'T', 'P':
		var token1, token3 int
		var token2, token4 rune
		_, err := fmt.Fscanf(reader, "%c%d:%c%d", &token2, &token1, &token4, &token3)
		if err != nil { /*TODO*/
		}

		param1 := token1 - 1        // old x
		param2 := int(token2 - 'a') // old y
		param3 := token3 - 1        // new x
		param4 := int(token4 - 'a') // new y

		return b.movePiece(param2, param1, param4, param3)

	case 'R', 'L':
		var token1 int
		var token2 rune
		_, err := fmt.Fscanf(reader, "%c%d", &token2, &token1)
		if err != nil { /*TODO*/
		}

		param1 := token1 - 1        // x
		param2 := int(token2 - 'a') // y
		param3 := inst

		return b.rotatePiece(param2, param1, param3)
	}

	if err != nil {
		//Error (no cumple con el formato)
	}
	return false
}
