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

type vector2_T struct{
	x int
	y int
}

type laserInteractionResult_T uint8
const(
	CONTINUE laserInteractionResult_T = 0
	HIT laserInteractionResult_T = 1
	STOP laserInteractionResult_T = 2
)

var laserMovementVector = [...]vector2_T {{0,1},{-1,0},{0,-1},{1,0}}

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
	processLaser(pointing_T) (pointing_T)
	//---Depuración---//
	VisualRep() string;
}
