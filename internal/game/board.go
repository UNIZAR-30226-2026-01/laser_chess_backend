package game

// DEFINICIÓN PENDIENTE DE REVISIÓN
// Fichero que se encarga de definir las estructuras
// de datos relacionadas con la lógica de una partida
// de laser chess

import (
	"fmt"
	"strings"
)

// --- Constantes --- //
const XDIM int = 10
const YDIM int = 8

type pointing_T uint8

const (
	DOWN  pointing_T = 0
	LEFT  pointing_T = 1
	UP    pointing_T = 2
	RIGHT pointing_T = 3
)

type team_T uint8

const (
	NONE      team_T = 0
	BLUE_TEAM team_T = 1
	RED_TEAM  team_T = 2
)

// ===================================== //
//	PIECES								 //
// ===================================== //

// Interfaz general, toda pieza de tablero debe tener esto definido
type BoardPiece interface {
	canMoveTo(x int, y int) bool
	canRotate(d rune) bool //temporal
	//---Depuración---//
	VisualRep() string
}

// ===================================== //
//	BOARD								 //
// ===================================== //

type Board struct {
	cells [XDIM][YDIM]BoardPiece
}

func (b *Board) movePiece(x_from int, y_from int, x_to int, y_to int) bool {
	return b.cells[x_from][y_from].canMoveTo(x_to, y_to)
}

func (b *Board) rotatePiece(x_at int, y_at int, rot rune) bool {
	return b.cells[x_at][y_at].canRotate(rot)
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

func (b *Board) processTurn(instruction string) bool {

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
