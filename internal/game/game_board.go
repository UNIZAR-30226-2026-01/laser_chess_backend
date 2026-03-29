package game

// ===================================== //
//	BOARD								 //
// ===================================== //

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Board struct {
	cells         [XDIM][YDIM]BoardPiece
	blueTeamLaser *BoardPieceLaser
	redTeamLaser  *BoardPieceLaser
}

func getTeamTile(x int, y int) team_T {
	if x == 0 {
		return BLUE_TEAM
	} else if x == 9 {
		return RED_TEAM
	} else if x == 1 && (y == 0 || y == 7) {
		return RED_TEAM
	} else if x == 8 && (y == 0 || y == 7) {
		return BLUE_TEAM
	} else {
		return NONE
	}
}

func InitBoard(csvBoard string) (*Board, error) {
	b := &Board{}

	var kingsBlue, kingsRed int
	var lasersBlue, lasersRed int

	reader := csv.NewReader(strings.NewReader(csvBoard))
	records, err := reader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("error al leer como CSV: %w", err)
	}

	//Recorrer el csv todos los "records" grabados para crear el tablero
	for y, fila := range records {
		for x, celdaStr := range fila {
			//creamos una pieza y la colocamos en la baldosa correspondiente
			pieza := constructorPieza(celdaStr)
			b.cells[x][y] = pieza

			//NO ES NECESARIA LA PARTE DE CONTAR, PERO AÑADE FIABILIDAD

			//Guardar los punteros a los laseres y contar para no declarar un
			//estado inicial inválido
			if laser, ok := pieza.(*BoardPieceLaser); ok {
				switch laser.team {
				case BLUE_TEAM:
					lasersBlue++
					b.blueTeamLaser = laser
				case RED_TEAM:
					lasersRed++
					b.redTeamLaser = laser
				}
			}

			if king, ok := pieza.(*BoardPieceKing); ok {
				switch king.team {
				case BLUE_TEAM:
					kingsBlue++
				case RED_TEAM:
					kingsRed++
				}
			}
		}
	}

	if lasersBlue != 1 || lasersRed != 1 {
		return nil, fmt.Errorf("Numero incorrecto de laseres - AZUL:%d -ROJO:%d", lasersBlue, lasersRed)
	}

	if kingsBlue != 1 || kingsRed != 1 {
		return nil, fmt.Errorf("Numero incorrecto de Reyes - AZUL:%d -ROJO:%d", kingsBlue, kingsRed)
	}

	return b, nil

}

func constructorPieza(codigo string) BoardPiece {
	if codigo == "" {
		return &BoardPieceVacant{}
	}

	// Extraemos Equipo
	var team team_T
	if len(codigo) >= 2 {
		switch codigo[1] {
		case 'A':
			team = BLUE_TEAM
		case 'R':
			team = RED_TEAM
		}
	}

	// Extraemos Dirección
	var dir pointing_T
	if len(codigo) >= 3 {
		switch codigo[2] {
		case 'U':
			dir = UP
		case 'D':
			dir = DOWN
		case 'L':
			dir = LEFT
		case 'R':
			dir = RIGHT
		}
	}

	// Mapeamos al struct correcto según la primera letra
	switch codigo[0] {
	case 'K': // King
		return &BoardPieceKing{team}
	case 'L': // Laser
		return &BoardPieceLaser{team, dir}
	case 'E': // Escudo (Shield)
		return &BoardPieceShield{team, dir}
	case 'S': // Switch
		return &BoardPieceSwitch{team, dir}
	case 'D': // Deflector
		return &BoardPieceDeflector{team, dir}
	default:
		return nil
	}
}

// devuelve true si la posicion que se le pasa esta dentro del tablero
func (b *Board) isInbound(x int, y int) bool {
	return (0 <= x && x < XDIM) && (0 <= y && y < YDIM)
}

func (b *Board) movePiece(x_from int, y_from int, x_to int, y_to int, team team_T) error {
	// Check board bounds
	if !(b.isInbound(x_from, y_from) && b.isInbound(x_to, y_to)) {
		return fmt.Errorf("Error - movimiento fuera de límites")
	}

	// Check reach
	if x_from-x_to < -1 || x_from-x_to > 1 || y_from-y_to < -1 || y_from-y_to > 1 {
		return fmt.Errorf("Error - movimiento fuera de rango")
	}

	// Check if can move
	err := b.cells[x_from][y_from].canMoveTo(x_to, y_to, b, team)
	if err != nil {
		return fmt.Errorf("Error - movimiento ilegal %s", err)
	}

	// Realizamos el movimiento legal
	destinyPiece := b.cells[x_to][y_to]
	originPiece := b.cells[x_from][y_from]

	b.cells[x_to][y_to] = originPiece
	b.cells[x_from][y_from] = destinyPiece

	return nil
}

func (b *Board) rotatePiece(x_at int, y_at int, rot rune, team team_T) error {
	// Check board bounds
	if !b.isInbound(x_at, y_at) {
		return fmt.Errorf("Error - rotación fuera de límites")
	}

	// Check syntax
	if !(rot == 'R' || rot == 'L') {
		return fmt.Errorf("Error - rotación mal especificada")
	}

	// Check laser legal rotations dependent on x/y
	switch laser := b.cells[x_at][y_at].(type) {
	case *BoardPieceLaser: //evitar rotacion ilegal de laser (Caso límite)
		x_after, y_after := laser.frontSpaceAfterRotating(x_at, y_at, rot)
		if !b.isInbound(x_after, y_after) {
			return fmt.Errorf("Error - rotación ilegal del laser")
		}
	}

	// Check if can rotate
	return b.cells[x_at][y_at].canRotate(rot, team)
}

// ---Depuración---//
func (b *Board) printlaser(laser []vector2_T) {
	for y := 0; y < YDIM; y++ {
		fmt.Printf("%d | ", 8-y) // numero
		for x := 0; x < XDIM; x++ {

			cell := b.cells[x][y].VisualRep()

			for i := 0; i < len(laser); i++ {
				if laser[i].x == x && laser[i].y == y {
					cell = "\033[47m" + cell + "\033[0m"
				}
			}

			fmt.Print(cell)
			fmt.Printf(" ")
		}
		fmt.Printf("\n")
	}
	fmt.Println("  +---------------------") // letra
	fmt.Println("    A B C D E F G H I J ") // letra
}

// ---Depuración---//
func (b *Board) print() {
	for y := 0; y < YDIM; y++ {
		fmt.Printf("%d | ", 8-y) // numero
		for x := 0; x < XDIM; x++ {

			cell := b.cells[x][y].VisualRep()

			fmt.Print(cell)
			fmt.Printf(" ")
		}
		fmt.Printf("\n")
	}
	fmt.Println("  +---------------------") // letra
	fmt.Println("    A B C D E F G H I J ") // letra
}

// --- INTERFAZ DE COMUNICACIÓN CON EL MÓDULO --- //

// Función auxiliar de cálculo de valores de retorno
func (b *Board) calculateReturnValues(instruction string, laserPath []vector2_T, result laserInteractionResult_T) (string, []vector2_T, laserInteractionResult_T, error) {
	retVal := instruction
	// Si el resultado de la interacción del laser es algún tipo de HIT
	if result > 2 && result < 6 {
		point := laserPath[len(laserPath)-1]
		retVal += "x" + string(rune(point.x+'a')) + strconv.Itoa(8-point.y) // y
		if b.isInbound(point.x, point.y) {
			b.cells[point.x][point.y] = &BoardPieceVacant{}
		}
	}
	return retVal, laserPath, result, nil
}

func (b *Board) ProcessTurn(instruction string, team team_T) (string, []vector2_T, laserInteractionResult_T, error) {

	//La versión de golang del string stream
	reader := strings.NewReader(instruction)

	var inst rune //Instrucción
	_, err := fmt.Fscanf(reader, "%c", &inst)
	if err != nil { /*TODO*/
	}

	switch inst {
	//MOVIMIENTO TRANSLACIÓN
	case 'T', 'P':
		var token1, token3 int
		var token2, token4 rune
		_, err := fmt.Fscanf(reader, "%c%d:%c%d", &token2, &token1, &token4, &token3)
		if err != nil { /*TODO*/
		}

		y_from := 8 - token1        // old x
		x_from := int(token2 - 'a') // old y
		y_to := 8 - token3          // new x
		x_to := int(token4 - 'a')   // new y

		legalMove := b.movePiece(x_from, y_from, x_to, y_to, team)

		if legalMove != nil {
			return "", nil, 0, legalMove
		}

		//PROCESAR EL MOVIMIENTO DEL LASER DEPENDIENDO DEL EQUIPO
		switch team {
		case BLUE_TEAM:
			laserPath, result := b.blueTeamLaser.shootLaser(0, 0, b)
			return b.calculateReturnValues(instruction, laserPath, result)
		case RED_TEAM:
			laserPath, result := b.redTeamLaser.shootLaser(XDIM-1, YDIM-1, b)
			return b.calculateReturnValues(instruction, laserPath, result)
		}

	case 'R', 'L':
		//ROTACIÓN
		var token1 int
		var token2 rune
		_, err := fmt.Fscanf(reader, "%c%d", &token2, &token1)
		if err != nil { /*TODO*/
		}

		y_at := 8 - token1        // x
		x_at := int(token2 - 'a') // y
		rot := inst

		legalMove := b.rotatePiece(x_at, y_at, rot, team)

		if legalMove != nil {
			return "", nil, 0, legalMove
		}

		//PROCESAR EL MOVIMIENTO DEPENDIENDO DEL EQUIPO
		switch team {
		case BLUE_TEAM:
			laserPath, result := b.blueTeamLaser.shootLaser(0, 0, b)
			return b.calculateReturnValues(instruction, laserPath, result)
		case RED_TEAM:
			laserPath, result := b.redTeamLaser.shootLaser(XDIM-1, YDIM-1, b)
			return b.calculateReturnValues(instruction, laserPath, result)
		}
	}

	return "", nil, 0, errors.New("Formato inválido")
}
