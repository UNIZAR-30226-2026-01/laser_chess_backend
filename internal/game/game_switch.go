package game

import "fmt"

// ============== SWITCH ============== //

type BoardPieceSwitch struct {
	team     team_T     //equipo
	pointing pointing_T //rotación
}

func (c *BoardPieceSwitch) canMoveTo(x int, y int, board *Board, team team_T) error {

	// Es ficha de tu equipo
	if (team != c.team){
		return fmt.Errorf("ficha del equipo opuesto")
	}

	// El movimiento termina en una casilla válida para tu ficha
	destinyTeamTile := getTeamTile(x, y)
	if !(destinyTeamTile == c.team || destinyTeamTile == NONE){
		return fmt.Errorf("casilla destino del equipo opuesto")
	}

	// Permutación válida en caso de switch son 3 tipos
	switch board.cells[x][y].(type) {
		case *BoardPieceVacant:
			return nil
		case *BoardPieceShield:
			return  nil
		case *BoardPieceDeflector:
			return  nil

	}

	// Casilla destino ocupada
	return fmt.Errorf("casilla destino ocupada")
}

func (c *BoardPieceSwitch) canRotate(d rune, team team_T) error {

	// Es ficha de tu equipo
	if (team != c.team){
		return fmt.Errorf("ficha del equipo opuesto")
	}
	
	switch d {
	case 'L': // -1 Counterclockwise
		c.pointing = (c.pointing + 3) % 4
		return nil
	case 'R': // +1 Clockwise
		c.pointing = (c.pointing + 1) % 4
		return nil
	default :
		//NO LLEGA NUNCA
		return fmt.Errorf("dirección mal especificada")
	} 
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