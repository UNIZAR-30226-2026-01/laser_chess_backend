package game

import "fmt"

// ============== SWITCH ============== //

type BoardPieceSwitch struct {
	team     team_T     //temporal
	tile     team_T     //Baldosa sobre la que estoy situado
	pointing pointing_T //temporal
}

func (c *BoardPieceSwitch) canMoveTo(x int, y int) bool {
	fmt.Printf("Switch - canMoveTo\n")
	return true //TODO
}

func (c *BoardPieceSwitch) canRotate(d rune) bool {
	fmt.Printf("Switch - canRotate\n")
	return true //TODO
}

//---Depuración---//
func (c *BoardPieceSwitch) VisualRep() string {
	var sprites = [4]string{"⧅", "⧄", "⧅", "⧄"}
	retval := sprites[c.pointing]
	if c.team == RED_TEAM {
		retval = "\033[31;1m" + retval + "\033[0m"
	}
	if c.team == BLUE_TEAM {
		retval = "\033[34;1m" + retval + "\033[0m"
	}
	return retval
}


func (c *BoardPieceSwitch)processLaser(dir pointing_T) (pointing_T, laserInteractionResult_T){
	switch (c.pointing - dir + 4) % 4 {
	case DOWN, UP:
		return (dir + 3) % 4, CONTINUE
	case RIGHT, LEFT:
		return (dir + 1) % 4, CONTINUE
	}

	return 0, 0
}