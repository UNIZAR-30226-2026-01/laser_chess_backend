package game

import "fmt"

// ============== LASER ============== //

type BoardPieceLaser struct {
	team     team_T     //temporal
	pointing pointing_T //temporal
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
* laserInteractionResult_T - Indica la razón por la que se ha detenido el avance del laser:
* 							 STOP: ha alcanzado un límite inamovible
*							 HIT: ha matado a una pieza
 */
func (c *BoardPieceLaser) shootLaser(x int, y int, board *Board) ([]vector2_T, laserInteractionResult_T) {
	var traveledPositions []vector2_T  //vector que guarda las posiciones que ha recorrido el laser
	currentPosition := vector2_T{x, y} //posicion del laser en una iteracion

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
		traveledPositions = append(traveledPositions, currentPosition)

		// avanzamos la posicion en la direccion del vector
		currentPosition.x += laserMovementDirectionVector.x
		currentPosition.y += laserMovementDirectionVector.y

		// llamamos a la función processLaser de la pieza que coincide con la posicion del laser en este momento
		// esta nos devolvera la nueva direccion a la que se dirige el haz y si este se detiene o continua
		laserPointingDirection, interactionRes = board.cells[currentPosition.x][currentPosition.y].processLaser(laserPointingDirection)

		//usamos esta nueva direccion para obtener el nuevo vector de movimiento
		laserMovementDirectionVector = laserMovementVector[laserPointingDirection]
	}
	// apilamos la posicion en la que se ha detenido el laser
	traveledPositions = append(traveledPositions, currentPosition)
	return traveledPositions, interactionRes
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

func (c *BoardPieceLaser)processLaser(dir pointing_T) (pointing_T, laserInteractionResult_T){
	return 0, STOP
}
