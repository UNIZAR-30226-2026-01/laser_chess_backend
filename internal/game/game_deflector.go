package game

import "fmt"

// ============== DEFLECTOR ============== //

type BoardPieceDeflector struct {
	team     team_T     //temporal
	tile     team_T     //Baldosa sobre la que estoy situado
	pointing pointing_T //temporal
}

func (c *BoardPieceDeflector) canMoveTo(x int, y int) bool {
	fmt.Printf("deflector - canMoveTo\n")
	return true //TODO
}

func (c *BoardPieceDeflector) canRotate(d rune) bool {
	fmt.Printf("deflector - canRotate\n")
	return true //TODO
}

//---Depuración---//
func (c *BoardPieceDeflector) VisualRep() string {
	var sprites = [4]string{"◣", "◤", "◥", "◢"}
	retval := sprites[c.pointing]
	if c.team == RED_TEAM {
		retval = "\033[31;1m" + retval + "\033[0m"
	}
	if c.team == BLUE_TEAM {
		retval = "\033[34;1m" + retval + "\033[0m"
	}
	return retval
}

func (c *BoardPieceDeflector) processLaser(dir pointing_T) (pointing_T, laserInteractionResult_T) {
	switch (c.pointing + dir) % 4 {
	case UP:
		return 0, HIT
	case LEFT:
		return UP, CONTINUE
	case DOWN:
		return RIGHT, CONTINUE
	case RIGHT:
		return 0, HIT
	}

	return 0, 0 //Nunca llega es para que no se queje el compilador
}
