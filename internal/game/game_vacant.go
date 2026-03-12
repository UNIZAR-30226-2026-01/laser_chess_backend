package game

import (
	"fmt"
)

// ============== VACANT ============== //

type BoardPieceVacant struct {
	tile team_T
}

func (c *BoardPieceVacant) canMoveTo(x int, y int) bool {
	fmt.Printf("Empty - canMoveTo\n")
	return false
}

func (c *BoardPieceVacant) canRotate(d rune) bool {
	fmt.Printf("Empty - canRotate\n")
	return false
}

// ---Depuración---//
func (c *BoardPieceVacant) VisualRep() string {
	retval := "·"
	if c.tile == RED_TEAM {
		retval = "\033[31;1m" + retval + "\033[0m"
	}
	if c.tile == BLUE_TEAM {
		retval = "\033[34;1m" + retval + "\033[0m"
	}
	return retval
}
