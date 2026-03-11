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
	//---Depuración---//
	VisualRep() string;
}

// ============== VACANT ============== //

type BoardPieceVacant struct{
	team rune //temporal (Casillas que no permiten )
}

func (c *BoardPieceVacant) canMoveTo(x int, y int) bool {
	fmt.Printf("Empty - canMoveTo\n")
	return false;
}

func (c *BoardPieceVacant) canRotate(d rune) bool {
	fmt.Printf("Empty - canRotate\n")
	return false;
}

//---Depuración---//
func (c *BoardPieceVacant) VisualRep() string {
	retval := "·"
	if (c.team == 'r') { retval = "\033[31;1m" + retval + "\033[0m"}
	if (c.team == 'b') { retval = "\033[34;1m" + retval + "\033[0m"}
	return retval
}

// ============== KING ============== //

type BoardPieceKing struct{
	team rune //Equipo al que pertenezco
	tile rune //Baldosa sobre la que estoy situado
}

func (c *BoardPieceKing) canMoveTo(x int, y int) bool {
	fmt.Printf("king - canMoveTo\n")
	return true; //TODO
}

func (c *BoardPieceKing) canRotate(d rune) bool {
	fmt.Printf("king - canRotate\n")
	return false; //TODO
}

//---Depuración---//
func (c *BoardPieceKing) VisualRep() string {
	retval := "K"
	if (c.team == 'r') { retval = "\033[31;1m" + retval + "\033[0m"}
	if (c.team == 'b') { retval = "\033[34;1m" + retval + "\033[0m"}
	return retval
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

//---Depuración---//
func (c *BoardPieceShield) VisualRep() string {
	retval := "S"
	if (c.team == 'r') { retval = "\033[31;1m" + retval + "\033[0m"}
	if (c.team == 'b') { retval = "\033[34;1m" + retval + "\033[0m"}
	return retval
}

// ===================================== //
//	BOARD								 //
// ===================================== //

type Board struct {
	cells [10][8]BoardPiece
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

func (b *Board) print(){
	for y := 0; y < 8; y++ {
		for x := 0; x < 10; x++ {
			cell := b.cells[x][y].VisualRep()
			fmt.Print(cell)
			fmt.Printf(" ")
		}
		fmt.Printf("\n")
	}
}