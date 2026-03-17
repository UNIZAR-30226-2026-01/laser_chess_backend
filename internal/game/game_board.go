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

func InitBoard(boardType board_T) Board {
	var newBoard Board
	switch boardType {
	case ACE:
		initACE(&newBoard)
	case CURIOSITY:
		//TODO
		break
	case GRAIL:
		//TODO
		break
	case MERCURY:
		//TODO
		break
	case SOPHIE:
		//TODO
		break
	}
	return newBoard
}

func initACE(tablero *Board) {

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
	tablero.cells[4][0] = &BoardPieceShield{BLUE_TEAM, NONE, DOWN} //[][1]
	tablero.cells[6][0] = &BoardPieceShield{BLUE_TEAM, NONE, DOWN}
	tablero.cells[3][7] = &BoardPieceShield{RED_TEAM, NONE, UP}
	tablero.cells[5][7] = &BoardPieceShield{RED_TEAM, NONE, UP}

	//poner switches
	tablero.cells[4][3] = &BoardPieceSwitch{BLUE_TEAM, NONE, DOWN} //DOWN
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

// devuelve true si la posicion que se le pasa esta dentro del tablero
func (b *Board) isInbound(x int, y int) bool {
	return (0 <= x && x < XDIM) && (0 <= y && y < YDIM)
}

func (b *Board) movePiece(x_from int, y_from int, x_to int, y_to int) bool {
	canmove := false
	if (b.isInbound(x_from, y_from) && b.isInbound(x_to, y_to)) {
		if (x_from - x_to < -1 || x_from - x_to > 1 || y_from - y_to < -1 || y_from - y_to > 1 ) {
			return false
		} else {
			canmove = b.cells[x_from][y_from].canMoveTo(x_to, y_to, b)
		}
	}

	if !canmove {
		return false
	}

	// Realizamos el movimiento legal
	destinyTileType := b.cells[x_to][y_to].getTeamTile()
	originTileType := b.cells[x_from][y_from].getTeamTile()
	
	b.cells[x_from][y_from].setTeamTile(destinyTileType)
	b.cells[x_from][y_from].setTeamTile(originTileType)

	destinyPiece := b.cells[x_to][y_to]
	originPiece := b.cells[x_from][y_from]

	b.cells[x_to][y_to] = originPiece
	b.cells[x_from][y_from] = destinyPiece

	return true
}

func (b *Board) rotatePiece(x_at int, y_at int, rot rune) bool {
	if (b.isInbound(x_at, y_at) && (rot == 'R' || rot == 'L')){
		switch laser := b.cells[x_at][y_at].(type) {
		case *BoardPieceLaser: //evitar rotacion ilegal de laser (Caso límite)
			x_after, y_after := laser.frontSpaceAfterRotating(x_at, y_at, rot)
			if !b.isInbound(x_after, y_after) {
				fmt.Print("OUT OF BOUND ROTATION")
				return false
			}
		}
	
		return b.cells[x_at][y_at].canRotate(rot)
	} else {
		return false
	}
	
}

//---Depuración---//
func (b *Board) printlaser(laser []vector2_T){
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

//---Depuración---//
func (b *Board) print(){
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
	//MOVIMIENTO TRANSLACIÓN
	case 'T', 'P':
		var token1, token3 int
		var token2, token4 rune
		_, err := fmt.Fscanf(reader, "%c%d:%c%d", &token2, &token1, &token4, &token3)
		if err != nil { /*TODO*/
		}

		y_from := token1 - 1        // old x
		x_from := int(token2 - 'a') // old y
		y_to := token3 - 1        // new x
		x_to := int(token4 - 'a') // new y

		

		return b.movePiece(x_from, y_from, x_to, y_to)

	case 'R', 'L':
	//ROTACIÓN
		var token1 int
		var token2 rune
		_, err := fmt.Fscanf(reader, "%c%d", &token2, &token1)
		if err != nil { /*TODO*/
		}

		y_at := token1 - 1        // x
		x_at := int(token2 - 'a') // y
		rot := inst

		return b.rotatePiece(x_at, y_at, rot)
	}

	if err != nil {
		//Error (no cumple con el formato)
	}
	return false
}
