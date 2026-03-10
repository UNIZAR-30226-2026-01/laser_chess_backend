package game

// DEFINICIÓN PENDIENTE DE REVISIÓN
// Fichero que se encarga de definir las estructuras
// de datos relacionadas con la lógica de una partida
// de laser chess

import "fmt"

// ===================================== //
//	PIECES								 //
// ===================================== //

// Interfaz general, toda pieza de tablero debe tener esto definido
type BoardPiece interface {
	canMoveTo(x int, y int) bool
	canRotate(d rune) bool //temporal
}

// ============== KING ============== //

type BoardPieceKing struct{
	team rune //temporal
}

func (c *BoardPieceKing) canMoveTo(x int, y int) bool {
	fmt.Printf("king - canMoveTo\n")
	return true; //TODO
}

func (c *BoardPieceKing) canRotate(d rune) bool {
	fmt.Printf("king - canRotate\n")
	return false; //TODO
}

// ============== SHIELD ============== //

type BoardPieceShield struct{
	team rune //temporal
	pointing int //temporal
}

func (c *BoardPieceShield) canMoveTo(x int, y int) bool {
	fmt.Printf("shield - canMoveTo\n")
	return true; //TODO
}

func (c *BoardPieceShield) canRotate(d rune) bool {
	fmt.Printf("shield - canRotate\n")
	return true; //TODO
}

// ===================================== //
//	BOARD								 //
// ===================================== //

type Board struct {
	cells [3][3]BoardPiece
}

func (b *Board) movePiece(x_from int, y_from int, x_to int, y_to int){
	if(b.cells[x_from][y_from].canMoveTo(x_to, y_to)){
		fmt.Printf("SAID TRUE\n")
	}else{
		fmt.Printf("SAID FALSE\n")
	}
}

func (b *Board) rotatePiece(x_at int, y_at int, rot rune){
	if(b.cells[x_at][y_at].canRotate(rot)){
		fmt.Printf("SAID TRUE\n")
	}else{
		fmt.Printf("SAID FALSE\n")
	}
}