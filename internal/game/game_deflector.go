package game

import "fmt"

// ============== DEFLECTOR ============== //

type BoardPieceDeflector struct {
	team     team_T     //temporal
	tile     team_T     //Baldosa sobre la que estoy situado
	pointing pointing_T //temporal
}

func (c *BoardPieceDeflector) getTeamTile() team_T {
	return c.tile
}

func (c *BoardPieceDeflector) setTeamTile(t team_T) {
	c.tile = t
}

func (c *BoardPieceDeflector) canMoveTo(x int, y int, board *Board, team team_T) bool {
	fmt.Printf("deflector - canMoveTo\n")

	if (team != c.team){
		return false
	}

	switch cell := board.cells[x][y].(type) {
		case *BoardPieceVacant:
			return c.team == cell.getTeamTile() || NONE == cell.getTeamTile() 
	}
	return false
}

func (c *BoardPieceDeflector) canRotate(d rune, team team_T) bool {
	fmt.Printf("deflector - canRotate\n")

	if (team != c.team){
		return false
	}
	
	switch d {
	case 'L': // -1 Counterclockwise
		c.pointing = (c.pointing + 3) % 4
		return true
	case 'R': // +1 Clockwise
		c.pointing = (c.pointing + 1) % 4
		return true
	}

	return false
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
	switch (c.pointing - dir + 4) % 4 {
	case DOWN:
		return (dir + 3) % 4, CONTINUE
	case RIGHT:
		return (dir + 1) % 4, CONTINUE
	case UP, LEFT:
		return 0, HIT
	}

	return 0, 0
}
