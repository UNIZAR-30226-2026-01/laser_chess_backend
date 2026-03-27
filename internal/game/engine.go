package game

import (
	"fmt"
	"strconv"

	boardtemplates "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game/boardTemplates"
)

//=======================================================================
//				DISEÑO INICIAL DE LA LÓGICA DEL JUEGO
//=======================================================================
// fichero que se encarga de la logica del juego de laser chess
// tendra una rutina con la maquina de estados (o lo que sea)
// principal del juego

// ESPECIFICACIÓN

/* --- TABLERO --- */

/* - MoverPiezaDesde(old_x int, old_y int, new_x int, new_y int)
* Desc:
 */

/* - RotarPieza(x int, y int, rotation char)
* Desc:
 */

/* - PermutarPiezas(old_x int, old_y int, new_x int, new_y int)
* Desc:
 */

/* --- Lógica de juego ---
- Se mira si el movimiento es legal
- Se ejecuta el movimiento
- Se activa el laser
- Se transforma el tablero
- Se comprueba condición de victoria
- Pasar turno
*/
//=======================================================================

//Esta clase sirve de mediador entre el tablero y la clase game, procesa el log de la partida

type GameEngine struct {
	gameLog   string
	gameBoard *Board
	boardType Board_T
}

func (g *GameEngine) initEngine(boardType Board_T) {
	var err error
	g.gameBoard, err = InitBoard(boardtemplates.ACE) // TODO: NO poner el nombre del csv, poner un switch case
	g.boardType = boardType
	if err != nil {
		fmt.Println("ey tio hay un error", err)
	}
}

func formatearLaserPath(laserPath []vector2_T) string {

	retVal := ""
	for i, point := range laserPath {
		// Transformación de los enteros a coordenadas de tablero
		retVal += string(rune(point.x+'a')) + strconv.Itoa(8-point.y)
		if i != len(laserPath)-1 {
			retVal += ","
		}
	}
	return retVal
}

func (g *GameEngine) getInitialState() string {
	switch g.boardType {
	case ACE:
		return boardtemplates.ACE
	case CURIOSITY:
		return boardtemplates.CURIOSITY
	case GRAIL:
		return boardtemplates.GRAIL
		// poner el resto
	default:
		return ""
	}
}

func (g *GameEngine) ProcessTurn(instruction string, team team_T) (string, []vector2_T, laserInteractionResult_T, error) {
	resul, laser, laserEnd, err := g.gameBoard.ProcessTurn(instruction, team)
	if err != nil {
		return resul, laser, laserEnd, err
	}
	//TODO -- Crear el log en cada turno
	g.gameLog += resul + "%" + formatearLaserPath(laser) + "%{300}" + ";"
	return resul, laser, laserEnd, err
}

func (g *GameEngine) GetState() string {
	return g.gameLog
}
