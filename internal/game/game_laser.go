package game

import "fmt"

// ============== LASER ============== //

type BoardPieceLaser struct {
	team     team_T     //temporal
	pointing pointing_T //temporal
}

func (c *BoardPieceLaser) canMoveTo(x int, y int) bool {
	fmt.Printf("Laser - canMoveTo\n")
	return false
}

func (c *BoardPieceLaser) canRotate(d rune) bool {
	fmt.Printf("Laser - canRotate\n")
	return true //TODO
}

//---Depuración---//
func (c *BoardPieceLaser) VisualRep() string {
	var sprites = [4]string{"▼", "◀", "▲", "▶"}
	retval := sprites[c.pointing]
	if c.team == RED_TEAM {
		retval = "\033[31;1m" + retval + "\033[0m"
	}
	if c.team == BLUE_TEAM {
		retval = "\033[34;1m" + retval + "\033[0m"
	}
	return retval
}
