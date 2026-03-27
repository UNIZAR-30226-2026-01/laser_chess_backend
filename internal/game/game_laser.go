package game

import "fmt"

// ============== LASER ============== //

type BoardPieceLaser struct {
	team     team_T     //temporal
	pointing pointing_T //temporal
}

func (c *BoardPieceLaser) getTeamTile() team_T {
	return c.team //NO HACE NADA SOLO CUMPLE CON INTERFAZ
}

func (c *BoardPieceLaser) setTeamTile(t team_T) {
	//NO HACE NADA, SOLO CUMPLE CON INTERFAZ
}

/*
* Desc: Esta funcion realiza el procesamiento del recorrido del haz laser en el tablero
*
* --- Parametros ---
* x int - Es la coordenada x en la que se encuentra la pieza laser.
* y int - Es la coordenada y en la que se encuentra la pieza laser.
* board *Board - Es un puntero al tablero de la partida, para que
* 				 el laser pueda tener acceso a el resto de piezas de la partida.
* --- Resultados ---
* []vector2_T - Contendrá todas las posiciones del tablero que ha recorrido el laser en orden, de la primera a la última.
*				MODIFICACIÓN - Ahora devuelve las esquinas del recorrido, los cambios de sentido
* laserInteractionResult_T - Indica la razón por la que se ha detenido el avance del laser:
* 							 STOP: ha alcanzado un límite inamovible
*							 HIT: ha matado a una pieza
 */
func (c *BoardPieceLaser) shootLaser(x int, y int, board *Board) ([]vector2_T, laserInteractionResult_T) {
	var traveledPositions []vector2_T  //vector que guarda las posiciones que ha recorrido el laser
	currentPosition := vector2_T{x, y} //posicion del laser en una iteracion
	//apilamos la primera posición
	traveledPositions = append(traveledPositions, currentPosition)

	// laserPointingDirection: direccion a la que apunta el laser,
	// se utiliza para procesar la direccion de la siguiente iteracion e indexar el vector de movimientos
	laserPointingDirection := c.pointing

	// laserMovementDirectionVector: Vector de dirección del laser,
	// este vector permite al laser avanzar en cada iteracion a la casilla que le corresponde
	// sumandoselo a la posicion actual del laser
	laserMovementDirectionVector := laserMovementVector[laserPointingDirection]
	interactionRes := CONTINUE //resultado de procesar la interaccion del haz laser con una pieza/casilla

	for interactionRes == CONTINUE {
		// apilamos la posicion actual del laser
		// YA NO, AHORA SOLO APILAMOS SI CAMBIA LA DIRECCIÓN

		// avanzamos la posicion en la direccion del vector
		currentPosition.x += laserMovementDirectionVector.x
		currentPosition.y += laserMovementDirectionVector.y

		if board.isInbound(currentPosition.x, currentPosition.y) {
			// llamamos a la función processLaser de la pieza que coincide con la posicion del laser en este momento
			// esta nos devolvera la nueva direccion a la que se dirige el haz y si este se detiene o continua
			laserPointingDirection, interactionRes = board.cells[currentPosition.x][currentPosition.y].processLaser(laserPointingDirection)

			//	si el vector de movimiento que tenemos ahora, va a cambiar en el siguiente movimiento
			//	esta casilla es una esquina, y la guardamos, mientras no sea un HIT porque en caso de hit la dirección devuelta
			// 	es irrelevante y la última casilla ya la guardamos al final
			if laserMovementDirectionVector != laserMovementVector[laserPointingDirection] && interactionRes != HIT {
				traveledPositions = append(traveledPositions, currentPosition)
			}

			//usamos esta la direccion para obtener el nuevo vector de movimiento
			laserMovementDirectionVector = laserMovementVector[laserPointingDirection]
		} else {
			interactionRes = OUT_OF_BOUNDS
		}
	}
	// apilamos la posicion en la que se ha detenido el laser
	traveledPositions = append(traveledPositions, currentPosition)
	return traveledPositions, interactionRes
}

func (c *BoardPieceLaser) canMoveTo(x int, y int, board *Board, team team_T) bool {

	if team != c.team {
		return false
	}

	fmt.Printf("Laser - canMoveTo\n")
	return false
}

func (c *BoardPieceLaser) canRotate(d rune, team team_T) bool {
	fmt.Printf("Laser - canRotate\n")

	if team != c.team {
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

func (c *BoardPieceLaser) frontSpaceAfterRotating(x int, y int, rot rune) (int, int) {
	var pointAux pointing_T

	switch rot {
	case 'L': // -1 Counterclockwise
		pointAux = (c.pointing + 3) % 4
	case 'R': // +1 Clockwise
		pointAux = (c.pointing + 1) % 4
	}

	laserMovementDirectionVector := laserMovementVector[pointAux]

	return (x + laserMovementDirectionVector.x), (y + laserMovementDirectionVector.y)
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

func (c *BoardPieceLaser) processLaser(dir pointing_T) (pointing_T, laserInteractionResult_T) {
	return 0, STOP
}

//DEBUG
func (c *BoardPieceLaser) printLaserInteractionResult(e laserInteractionResult_T) string {
	switch e {
	case 0:
		return "CONTINUE"
	case 1:
		return "HIT"
	case 2:
		return "STOP"
	default:
		return "OUT_OF_BOUNDS"
	}
}
