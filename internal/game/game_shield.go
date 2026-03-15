package game

import (
	"fmt"
)

// ============== SHIELD ============== //

type BoardPieceShield struct {
	team     team_T     //temporal
	tile     team_T     //Baldosa sobre la que estoy situado
	pointing pointing_T //temporal
}

func (c *BoardPieceShield) canMoveTo(x int, y int, board *Board) bool {
	fmt.Printf("shield - canMoveTo\n")
	return true //TODO
}

func (c *BoardPieceShield) canRotate(d rune) bool {
	fmt.Printf("shield - canRotate\n")
	return true //TODO
}

// ---Depuración---//
func (c *BoardPieceShield) VisualRep() string {
	var sprites = [4]string{"⬓", "◧", "⬒", "◨"}
	retval := sprites[c.pointing]
	if c.team == RED_TEAM {
		retval = "\033[31;1m" + retval + "\033[0m"
	}
	if c.team == BLUE_TEAM {
		retval = "\033[34;1m" + retval + "\033[0m"
	}
	return retval
}

func (c *BoardPieceShield)processLaser(dir pointing_T) (pointing_T, laserInteractionResult_T){
	switch ((c.pointing - dir)%4){
	case UP:
		return 0, STOP
	case DOWN, LEFT, RIGHT:
		return 0, HIT
	}

	return 0,0 //Nunca llega es para que no se queje eboard *Boardl compilador
}