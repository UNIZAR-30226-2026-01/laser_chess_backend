package game

import (
	"fmt"
)

// ============== VACANT ============== //

type BoardPieceVacant struct {
}

func (c *BoardPieceVacant) canMoveTo(x int, y int, board *Board, team team_T) error {
	return fmt.Errorf("No hay una pieza en esta casilla")
}

func (c *BoardPieceVacant) canRotate(d rune, team team_T) error {
	return fmt.Errorf("No hay una pieza en esta casilla")
}

// ---Depuración---//
func (c *BoardPieceVacant) VisualRep() string {
	retval := "·"
	// if c.tile == RED_TEAM {
	// 	retval = "\033[31;1m" + retval + "\033[0m"
	// }
	// if c.tile == BLUE_TEAM {
	// 	retval = "\033[34;1m" + retval + "\033[0m"
	// }
	return retval
}

func (c *BoardPieceVacant)processLaser(dir pointing_T) (pointing_T, laserInteractionResult_T){
	return dir, CONTINUE
}