package game

// DEFINICIÓN PENDIENTE DE REVISIÓN
// Fichero que se encarga de definir las estructuras
// de datos relacionadas con la lógica de una partida
// de laser chess

import "fmt"

// --- Constantes --- //
const XDIM int = 10
const YDIM int = 8

type pointing_T uint8
const (
	DOWN pointing_T = 0
	LEFT pointing_T = 1
	UP pointing_T = 2
	RIGHT pointing_T = 3

)

type team_T uint8
const (
	NONE team_T = 0
	BLUE_TEAM team_T = 1
	RED_TEAM team_T = 2
)

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
	tile team_T 
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
	if (c.tile == RED_TEAM) { retval = "\033[31;1m" + retval + "\033[0m"}
	if (c.tile == BLUE_TEAM) { retval = "\033[34;1m" + retval + "\033[0m"}
	return retval
}

// ============== KING ============== //

type BoardPieceKing struct{
	team team_T //Equipo al que pertenezco
	tile team_T //Baldosa sobre la que estoy situado
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
	if (c.team == RED_TEAM) { retval = "\033[31;1m" + retval + "\033[0m"}
	if (c.team == BLUE_TEAM) { retval = "\033[34;1m" + retval + "\033[0m"}
	return retval
}

// ============== SHIELD ============== //

type BoardPieceShield struct{
	team team_T //temporal
	tile team_T //Baldosa sobre la que estoy situado
	pointing pointing_T //temporal
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
	var sprites = [4]string{"⬓", "◧", "⬒", "◨"}
	retval := sprites[c.pointing]
	if (c.team == RED_TEAM) { retval = "\033[31;1m" + retval + "\033[0m"}
	if (c.team == BLUE_TEAM) { retval = "\033[34;1m" + retval + "\033[0m"}
	return retval
}

// ============== DEFLECTOR ============== //

type BoardPieceDeflector struct{
	team team_T //temporal
	tile team_T //Baldosa sobre la que estoy situado
	pointing pointing_T //temporal
}

func (c *BoardPieceDeflector) canMoveTo(x int, y int) bool {
	fmt.Printf("deflector - canMoveTo\n")
	return true; //TODO
}

func (c *BoardPieceDeflector) canRotate(d rune) bool {
	fmt.Printf("deflector - canRotate\n")
	return true; //TODO
}

//---Depuración---//
func (c *BoardPieceDeflector) VisualRep() string {
	var sprites = [4]string{"◣", "◤", "◥", "◢"}
	retval := sprites[c.pointing]
	if (c.team == RED_TEAM) { retval = "\033[31;1m" + retval + "\033[0m"}
	if (c.team == BLUE_TEAM) { retval = "\033[34;1m" + retval + "\033[0m"}
	return retval
}

// ============== SWITCH ============== //

type BoardPieceSwitch struct{
	team team_T //temporal
	tile team_T //Baldosa sobre la que estoy situado
	pointing pointing_T //temporal
}

func (c *BoardPieceSwitch) canMoveTo(x int, y int) bool {
	fmt.Printf("Switch - canMoveTo\n")
	return true; //TODO
}

func (c *BoardPieceSwitch) canRotate(d rune) bool {
	fmt.Printf("Switch - canRotate\n")
	return true; //TODO
}

//---Depuración---//
func (c *BoardPieceSwitch) VisualRep() string {
	var sprites = [4]string{"⧅", "⧄", "⧅", "⧄"}
	retval := sprites[c.pointing]
	if (c.team == RED_TEAM) { retval = "\033[31;1m" + retval + "\033[0m"}
	if (c.team == BLUE_TEAM) { retval = "\033[34;1m" + retval + "\033[0m"}
	return retval
}

// ============== LASER ============== //

type BoardPieceLaser struct{
	team team_T //temporal
	pointing pointing_T //temporal
}

func (c *BoardPieceLaser) canMoveTo(x int, y int) bool {
	fmt.Printf("Switch - canMoveTo\n")
	return true; //TODO
}

func (c *BoardPieceLaser) canRotate(d rune) bool {
	fmt.Printf("Switch - canRotate\n")
	return true; //TODO
}

//---Depuración---//
func (c *BoardPieceLaser) VisualRep() string {
	var sprites = [4]string{"▼", "◀", "▲", "▶"}
	retval := sprites[c.pointing]
	if (c.team == RED_TEAM) { retval = "\033[31;1m" + retval + "\033[0m"}
	if (c.team == BLUE_TEAM) { retval = "\033[34;1m" + retval + "\033[0m"}
	return retval
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

//---Depuración---//
func (b *Board) print(){
	for y := 0; y < YDIM; y++ {
		for x := 0; x < XDIM; x++ {
			cell := b.cells[x][y].VisualRep()
			fmt.Print(cell)
			fmt.Printf(" ")
		}
		fmt.Printf("\n")
	}
}