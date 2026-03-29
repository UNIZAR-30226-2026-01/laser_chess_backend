package game

import (
	"fmt"
)

// ============== KING ============== //

type BoardPieceKing struct {
	team team_T //Equipo al que pertenezco
}

func (c *BoardPieceKing) canMoveTo(x int, y int, board *Board, team team_T) error {

	// Es ficha de tu equipo
	if team != c.team {
		return fmt.Errorf("ficha del equipo opuesto")
	}

	// El movimiento termina en una casilla válida para tu ficha
	destinyTeamTile := getTeamTile(x, y)
	if !(destinyTeamTile == c.team || destinyTeamTile == NONE) {
		return fmt.Errorf("casilla destino del equipo opuesto")
	}

	// Permutación válida en caso de switch son 3 tipos
	switch board.cells[x][y].(type) {
	case *BoardPieceVacant:
		return nil
	case *BoardPieceShield:
		return nil
	case *BoardPieceDeflector:
		return nil

	}

	// Casilla destino ocupada
	return fmt.Errorf("casilla destino ocupada")
}

func (c *BoardPieceKing) canRotate(d rune, team team_T) error {
	return fmt.Errorf("rey no puede rotar")
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

func (c *BoardPieceKing) processLaser(dir pointing_T) (pointing_T, laserInteractionResult_T) {
	switch c.team {
	case RED_TEAM:
		return 0, HIT_RED_KING
	case BLUE_TEAM:
		return 0, HIT_BLUE_KING
	default:
		// Imposible
		return 0, HIT
	}
}
