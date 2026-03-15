package game

import (
	"fmt"
)

// ============== KING ============== //

type BoardPieceKing struct {
	team team_T //Equipo al que pertenezco
	tile team_T //Baldosa sobre la que estoy situado
}

func (c *BoardPieceKing) canMoveTo(x int, y int, board *Board) bool {
	fmt.Printf("king - canMoveTo\n")
	return true //TODO
}

func (c *BoardPieceKing) canRotate(d rune) bool {
	fmt.Printf("king - canRotate\n")
	return false 
}

// ---Depuración---//
func (c *BoardPieceKing) VisualRep() string {
	retval := "K"
	if c.team == RED_TEAM {
		retval = "\033[31;1m" + retval + "\033[0m"
	}
	if c.team == BLUE_TEAM {
		retval = "\033[34;1m" + retval + "\033[0m"
	}
	return retval
}

func (c *BoardPieceKing)processLaser(dir pointing_T) (pointing_T, laserInteractionResult_T){
	return 0, HIT
}
