package game

// Definición de constantes e interfaz común entre piezas

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

// Interfaz general, toda pieza de tablero debe tener esto definido
type BoardPiece interface {
	canMoveTo(x int, y int) bool
	canRotate(d rune) bool //temporal
	//---Depuración---//
	VisualRep() string;
}
