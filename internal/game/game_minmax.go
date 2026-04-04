package game

import (
	"fmt"
	"math/rand/v2"
)

//Struct para acelerar los cáclulos
type AIGameState struct {
	rKings []int
	bKings []int
	rShields []int
	bShields []int
	rSwitches	[]int
	bSwitches	[]int
	rDeflectors	[]int
	bDeflectors []int

}

type AIMove struct {
	mType 		 rune
	fromX, fromY int
	toX, toY     int
	team 		 team_T
	score        float64 
}

const (
	MAX_SCORE  = 9999999
	MIN_SCORE = -9999999
)

//Iniciar el estado para acelerar cálculos
func initGameState(b *Board, a *AIGameState){

	for y := 0; y < YDIM; y++ {
		for x := 0; x < XDIM; x++ {
			switch v := b.cells[x][y].(type) {
			case *BoardPieceShield:
				switch v.team {
				case RED_TEAM:	
					a.rShields = append(a.rShields, x+(y*XDIM))
				case BLUE_TEAM:
					a.bShields = append(a.bShields, x+(y*XDIM))
				}
			case *BoardPieceDeflector:
				switch v.team {
				case RED_TEAM:	
					a.rDeflectors = append(a.rDeflectors, x+(y*XDIM))
				case BLUE_TEAM:
					a.bDeflectors= append(a.bDeflectors, x+(y*XDIM))
				}
			case *BoardPieceSwitch:
				switch v.team {
				case RED_TEAM:	
					a.rSwitches = append(a.rSwitches, x+(y*XDIM))
				case BLUE_TEAM:
					a.bSwitches = append(a.bSwitches, x+(y*XDIM))
				}
			case *BoardPieceKing:
				switch v.team {
				case RED_TEAM:	
					a.rKings = append(a.rKings, x+(y*XDIM))
				case BLUE_TEAM:
					a.bKings = append(a.bKings, x+(y*XDIM))
				}
			}
		}
	}
}

// ----------------------------- AUXILIARES --------------------------------- //
// TEMPORALES!
func cloneBoard(b *Board) *Board {
	newBoard := &Board{}
	for x := 0; x < XDIM; x++ {
		for y := 0; y < YDIM; y++ {
			newBoard.cells[x][y] = clonePiece(b.cells[x][y])
			if laser, ok := newBoard.cells[x][y].(*BoardPieceLaser); ok {
				if laser.team == BLUE_TEAM {
					newBoard.blueTeamLaser = laser
				} else {
					newBoard.redTeamLaser = laser
				}
			}
		}
	}
	return newBoard
}

func clonePiece(p BoardPiece) BoardPiece {
	switch v := p.(type) {
	case *BoardPieceDeflector:
		return &BoardPieceDeflector{v.team, v.pointing}
	case *BoardPieceShield:
		return &BoardPieceShield{v.team, v.pointing}
	case *BoardPieceSwitch:
		return &BoardPieceSwitch{v.team, v.pointing}
	case *BoardPieceLaser:
		return &BoardPieceLaser{v.team, v.pointing}
	case *BoardPieceKing:
		return &BoardPieceKing{v.team}
	default:
		return &BoardPieceVacant{}
	}
}

func (m AIMove) String() string {
	if m.mType == 'T' || m.mType == 'P' {
		return fmt.Sprintf("T%c%d:%c%d", m.fromX+'a', 8-m.fromY, m.toX+'a', 8-m.toY)
	}
		return fmt.Sprintf("%c%c%d", m.mType, m.fromX+'a', 8-m.fromY)
}

// ----------------------------- NECESARIAS --------------------------------- //

func (a *AIGameState)get8DirectionsAux(b* Board, px int, py int, team team_T) []AIMove {
	moves := make([]AIMove, 0, 8)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := px+dx, py+dy
			if b.isInbound(nx, ny) && b.cells[px][py].canMoveTo(nx, ny, b, team) == nil {
				moves = append(moves, AIMove{mType: 'T', fromX: px, fromY: py, toX: nx, toY: ny, team:team})
			}
		}
	}
	return moves
}

/* Desc: Función que devuelve todos los movimientos de la frontera a explorar
  --- Parametros ---
	b *Board 	- Tablero
	team team_T - Turno del equipo team
  --- Resultados ---
	[]AIMove 	- Devuelve una lista con todos los movimientos válidos de cada pieza
*/
func (a *AIGameState)getFrontier(b *Board, team team_T) []AIMove {
	moves := make([]AIMove, 0, 50)
	switch team {
	case RED_TEAM:
		// LASER
		if(b.redTeamLaser.pointing == UP){
			moves = append(moves, AIMove{mType: 'L', fromX: XDIM-1, fromY:YDIM-1, team:team})
		}else{
			moves = append(moves, AIMove{mType: 'R', fromX: XDIM-1, fromY:YDIM-1, team:team})
		}
		// SWITCHES
		for _, p := range a.rSwitches {
			px := p % XDIM
			py := p / XDIM
			//ROTAR - si
			moves = append(moves, AIMove{mType: 'L', fromX: px, fromY: py, team:team})
			moves = append(moves, AIMove{mType: 'R', fromX: px, fromY: py, team:team})
			//MOVER - si
			moves = append(moves, a.get8DirectionsAux(b, px, py, team)...)
		}
		// ESCUDOS
			for _, p := range a.rShields {
			px := p % XDIM
			py := p / XDIM
			//ROTAR - si
			moves = append(moves, AIMove{mType: 'L', fromX: px, fromY: py, team:team})
			moves = append(moves, AIMove{mType: 'R', fromX: px, fromY: py, team:team})
			//MOVER - si
			moves = append(moves, a.get8DirectionsAux(b, px, py, team)...)
		}
		// DEFLECTORES
			for _, p := range a.rDeflectors {
			px := p % XDIM
			py := p / XDIM
			//ROTAR - si
			moves = append(moves, AIMove{mType: 'L', fromX: px, fromY: py, team:team})
			moves = append(moves, AIMove{mType: 'R', fromX: px, fromY: py, team:team})
			//MOVER - si
			moves = append(moves, a.get8DirectionsAux(b, px, py, team)...)
		}
		// REYES
		for _, p := range a.rKings {
			px := p % XDIM
			py := p / XDIM
			//ROTAR - no
			//MOVER - si
			moves = append(moves, a.get8DirectionsAux(b, px, py, team)...)
		}

		case BLUE_TEAM:
		// LASER
		if(b.blueTeamLaser.pointing == DOWN){
			moves = append(moves, AIMove{mType: 'L', fromX: 0, fromY:0, team:team})
		}else{
			moves = append(moves, AIMove{mType: 'R', fromX: 0, fromY:0, team:team})
		}
		// SWITCHES
		for _, p := range a.bSwitches {
			px := p % XDIM
			py := p / XDIM
			//ROTAR - si
			moves = append(moves, AIMove{mType: 'L', fromX: px, fromY: py, team:team})
			moves = append(moves, AIMove{mType: 'R', fromX: px, fromY: py, team:team})
			//MOVER - si
			moves = append(moves, a.get8DirectionsAux(b, px, py, team)...)
		}
		// ESCUDOS
			for _, p := range a.bShields {
			px := p % XDIM
			py := p / XDIM
			//ROTAR - si
			moves = append(moves, AIMove{mType: 'L', fromX: px, fromY: py, team:team})
			moves = append(moves, AIMove{mType: 'R', fromX: px, fromY: py, team:team})
			//MOVER - si
			moves = append(moves, a.get8DirectionsAux(b, px, py, team)...)
		}
		// DEFLECTORES
			for _, p := range a.bDeflectors {
			px := p % XDIM
			py := p / XDIM
			//ROTAR - si
			moves = append(moves, AIMove{mType: 'L', fromX: px, fromY: py, team:team})
			moves = append(moves, AIMove{mType: 'R', fromX: px, fromY: py, team:team})
			//MOVER - si
			moves = append(moves, a.get8DirectionsAux(b, px, py, team)...)
		}
		// REYES
		for _, p := range a.bKings {
			px := p % XDIM
			py := p / XDIM
			//ROTAR - no
			//MOVER - si
			moves = append(moves, a.get8DirectionsAux(b, px, py, team)...)
		}
		// LASER
			// TODO
	}
	return moves
}

/* Desc: Dado un estado, devuelve la valoración de la posición mediante una serie de heuristicas

  --- Parametros ---
	b *Board 	- Puntero al tablero de la partida, estado del juego
	team team_T - evaluación respecto al

  --- Resultados ---

*/
func (a *AIGameState)evaluateBoard(b *Board) (score int){
	score = 0
	//Comer piezas,  quién tenga más piezas está en mejor posición
	score += (len(a.rDeflectors) - len(a.bDeflectors)) * 30
	score += (len(a.rShields) - len(a.bShields)) * 50
	score += (len(a.rKings) - len(a.bKings)) * 10000
	//Control del laser, quién controle más puntos de inflexion está en mejor posicion
	laserPath, _ := b.redTeamLaser.shootLaser(XDIM-1, YDIM-1, b)
	for _, pair := range laserPath[1 : len(laserPath)-1] {
		switch p := b.cells[pair.x][pair.y].(type){
		case *BoardPieceDeflector:
			switch p.team {
			case RED_TEAM: score += 2
			case BLUE_TEAM: score -= 1
			}
		case *BoardPieceSwitch:
			switch p.team {
			case RED_TEAM: score += 2
			case BLUE_TEAM: score -= 1
			}
		case *BoardPieceShield:
			switch p.team {
			case RED_TEAM: score += 2
			case BLUE_TEAM: score -= 1
			}
		default:
		}
	}
	//Control del laser, quién controle más puntos de inflexion está en mejor posicion
	laserPath, _ = b.blueTeamLaser.shootLaser(0, 0, b)
	for _, pair := range  laserPath[1 : len(laserPath)-1] {
		switch p := b.cells[pair.x][pair.y].(type){
		case *BoardPieceDeflector:
			switch p.team {
			case RED_TEAM: score += 1
			case BLUE_TEAM: score -= 2
			}
		case *BoardPieceSwitch:
			switch p.team {
			case RED_TEAM: score += 1
			case BLUE_TEAM: score -= 2
			}
		case *BoardPieceShield:
			switch p.team {
			case RED_TEAM: score += 1
			case BLUE_TEAM: score -= 2
			}
		default:
		}
	}

	score += evaluateKingDefense(b, a.rKings[0], RED_TEAM)
	score -= evaluateKingDefense(b, a.bKings[0], BLUE_TEAM)
	score += rand.IntN(6) - 3

	return score
	
}

func evaluateKingDefense(b *Board, king int, team team_T) int {
	bonus := 0
	kingPos := vector2_T{x:king%XDIM,y:king/XDIM}
	// Rodear al rey
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := kingPos.x+dx, kingPos.y+dy
			if b.isInbound(nx, ny) {
				switch v := b.cells[nx][ny].(type) {
				case *BoardPieceShield:
					if v.team == team {
						bonus += 10
					}
				case *BoardPieceDeflector:
					if v.team == team {
						bonus += 5
					}else{
						bonus -= 10
					}
				case *BoardPieceSwitch:
					if v.team == team {
						bonus += 3
					}else{
						bonus -= 7
					}
				}
			}
		}
	}

	return bonus
}

/* Desc: Dado un movimiento válido devuelve el estado subsecuente del juego 

  --- Parametros ---
	b *Board 	- Tablero
	team team_T - Turno del equipo team
  --- Resultados ---
	[]AIMove 	- Devuelve una lista con todos los movimientos válidos de cada pieza
*/
func transitionFunc(b *Board, move AIMove) *Board {
	nb := cloneBoard(b)
	var laserPath []vector2_T
	var termRes laserInteractionResult_T 

	//Aplicar transformación
	switch move.mType {
	case 'T':
		nb.movePiece(move.fromX, move.fromY, move.toX, move.toY, move.team)
	case 'R':
		nb.rotatePiece(move.fromX ,move.fromY, 'R', move.team)
	case 'L':
		nb.rotatePiece(move.fromX ,move.fromY, 'L', move.team)
	}
	//Disparar laser
	switch move.team {
	case RED_TEAM: 
		laserPath, termRes = nb.redTeamLaser.shootLaser(XDIM-1, YDIM-1, nb)
	case BLUE_TEAM:
		laserPath, termRes = nb.blueTeamLaser.shootLaser(0, 0, nb)
	}
	//Procesar disparo
	point := laserPath[len(laserPath)-1]
	if nb.isInbound(point.x, point.y) && termRes >= 3 {
		nb.cells[point.x][point.y] = &BoardPieceVacant{}
	}

	return nb
}

// ----------------------------- FUNDAMENTALES ------------------------------ //


/* TODO
  Desc: Algorimo MINMAX
  --- Parametros ---
	b *Board 	- Puntero al tablero de la partida, estado raíz del árbol de juego
	team team_T - Siguiente equipo que debe mover, maz

  --- Resultados ---

*/
func minmax(b *Board, depth int, alpha int, beta int, myTeam team_T) (score int, move AIMove) {
    bestMove := AIMove{}
    a := AIGameState{}
    initGameState(b, &a)
    
    // 0. Caso base, fin de la recursividad, estado final de partida "hojas"
    if (len(a.bKings) == 0) {
        return MAX_SCORE + depth * 1000 , bestMove
    } else if(len(a.rKings) == 0){
        return MIN_SCORE + depth * -1000, bestMove
    } else if(depth == 0) {
        score = a.evaluateBoard(b)
        return score, bestMove
    }

    // 1. Obtener la frontera de movimientos posibles desde este estado
    frontier := a.getFrontier(b, myTeam)

    // 2. Aplicar los movimientos y devolver el mejor
    switch myTeam {
    case RED_TEAM: //ROJO MAXIMIZA
        maximizedScore := MIN_SCORE
        for _, move := range frontier {
            nextState := transitionFunc(b, move)
            // Pasamos alpha y beta a la llamada recursiva
            score, _ := minmax(nextState, depth - 1, alpha, beta, BLUE_TEAM)
            
            if maximizedScore < score {
                maximizedScore = score
                bestMove = move
            }
            
            if alpha < score {
                alpha = score // Actualizamos el mejor valor que Rojo puede asegurar
            }
            if beta <= alpha {
                break // Poda: Azul nunca permitiría llegar a esta rama
            }
        }
        return maximizedScore + depth * 5, bestMove
        
    case BLUE_TEAM: // AZUL MINIMIZA
        minimizedScore := MAX_SCORE
        for _, move := range frontier {
            nextState := transitionFunc(b, move)
            // Pasamos alpha y beta a la llamada recursiva
            score, _ := minmax(nextState, depth - 1, alpha, beta, RED_TEAM)
            
            if minimizedScore > score {
                minimizedScore = score
                bestMove = move
            }
            
            if beta > score {
                beta = score // Actualizamos el mejor valor que Azul puede asegurar
            }
            if beta <= alpha {
                break // Poda: Rojo nunca permitiría llegar a esta rama
            }
        }
        return minimizedScore + depth * -5 , bestMove
    }
    return 0, bestMove
}

/*
  Desc: Interfaz de la IA - minmax, árbol de juego
  --- Parametros ---
	b *Board 	- Puntero al tablero de la partida, estado raíz del árbol de juego
	team team_T - Siguiente equipo que debe mover.
	lvl int 	- Profundidad de la búsqueda en el arbol de juego
  --- Resultados ---
	string - devuelve el siguiente mejor movimiento posible
*/
func GetBestMove(b *Board, team team_T, lvl int) string {
	//crear el estado inicial del motor de minmax para acelerar la evaluación
	score, move := minmax(b, lvl, MIN_SCORE, MAX_SCORE, team)
	fmt.Printf("El score es %d\n", score)
	fmt.Print(move.fromX, move.fromY, move.toX, move.toY)

	return move.String()
}

