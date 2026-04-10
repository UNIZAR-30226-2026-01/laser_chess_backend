package game

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

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

// --- Auxiliares --- //

func (g *GameEngine) getInitialState() string {
	switch g.boardType {
	case ACE:
		return boardtemplates.ACE
	case CURIOSITY:
		return boardtemplates.CURIOSITY
	case GRAIL:
		return boardtemplates.GRAIL
	case SOPHIE:
		return boardtemplates.SOPHIE
	case MERCURY:
		return boardtemplates.MERCURY
	default:
		return ""
	}
}

func formatLaserPath(laserPath []vector2_T) string {

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

func formatTimeLeft(timeLeft time.Duration) string {
	timeStr := fmt.Sprintf("%.2f", timeLeft.Seconds())
	return "%{" + timeStr + "}" + ";"
}

func (g *GameEngine) InitEngine(boardType Board_T) {
	var err error
	g.boardType = boardType
	g.gameBoard, err = InitBoard(g.getInitialState())
	if err != nil {
		fmt.Println("error al inicializar el tablero", err)
	}
}

func (g *GameEngine) ProcessTurn(instruction string, team team_T, timeLeft time.Duration) (string, []vector2_T, laserInteractionResult_T, error) {
	result, laser, laserEnd, err := g.gameBoard.ProcessTurn(instruction, team)
	if err != nil {
		return result, laser, laserEnd, err
	}

	timeLeftStr := formatTimeLeft(timeLeft)
	g.gameLog += result + "%" + formatLaserPath(laser) + timeLeftStr

	result += timeLeftStr
	return result, laser, laserEnd, err
}

func (g *GameEngine) GetState() string {
	return g.gameLog
}

// SE PRESUPONE UN LOG QUE NO CAUSA ERRORES
func (g *GameEngine) ApplyLogToBoard(timeBase int32) (nextTeam team_T, redTimeLeft float64, blueTimeLeft float64) {
	//dividimos el log en cachitos
	logChunks := strings.Split(strings.TrimSuffix(g.gameLog, ";"), ";")
	nextTeam = RED_TEAM                                           //Equipo que empieza
	re := regexp.MustCompile(`^([^%]+)%(?:[^%]+)%\{([0-9.]+)\}$`) //regla expresión regular
	redTimeLeft = float64(timeBase)
	blueTimeLeft = float64(timeBase)

	//aplicamos cada cachito usando processTurn y los procesamos.
	for _, logChunk := range logChunks {
		//Tokenizamos usando la expresión regular
		tokens := re.FindStringSubmatch(logChunk)
		move := tokens[1]

		time, _ := strconv.ParseFloat(tokens[2], 64)

		//Aplicamos el movimiento
		g.gameBoard.ProcessTurn(move, nextTeam)

		// Permutamos entre los equipos y actualizamos el tiempo restante
		switch nextTeam {
		case BLUE_TEAM:
			blueTimeLeft = time
			nextTeam = RED_TEAM
		case RED_TEAM:
			redTimeLeft = time
			nextTeam = BLUE_TEAM
		}
	}

	return nextTeam, redTimeLeft, blueTimeLeft
}
