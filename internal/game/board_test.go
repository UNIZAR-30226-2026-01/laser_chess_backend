package game

import (
	"fmt"
	"testing"

	boardtemplates "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game/boardTemplates"
)

func TestAllBoards(t *testing.T) {
	tablero, err := InitBoard(boardtemplates.ACE)
	if err != nil {
		t.Error(err)
	}

	fmt.Print("== ACE ==\n")

	tablero.print()

	fmt.Print("== CURIOSITY ==\n")

	tablero, err = InitBoard(boardtemplates.CURIOSITY)
	if err != nil {
		t.Error(err)
	}
	tablero.print()

	fmt.Print("== GRAIL ==\n")
	tablero, err = InitBoard(boardtemplates.GRAIL)
	if err != nil {
		t.Error(err)
	}
	tablero.print()

	fmt.Print("== MERCURY ==\n")
	// tablero.InitBoard("boardTemplates/mercury.csv") TODO
	// tablero.print()

	fmt.Print("== SOPHIE ==\n")
	// tablero.InitBoard("boardTemplates/sophie.csv") TODO
	// tablero.print()
}

// Test final de una partida y su resultado esperado
func TestMovements(t *testing.T) {

	fmt.Print("== TEST TRANSFORMACIONES ==\n")

	//Iniciar tablero
	tablero, err := InitBoard(boardtemplates.CURIOSITY)
	if err != nil {
		t.Error(err)
	}
	tablero.print()

	tablero.print()

	var log string
	var logPiece string
	var laser_end laserInteractionResult_T
	var path []vector2_T

	//Ejemplo de procesamiento
	movimiento_front := "La8"

	fmt.Println("LLEGA:\t" + movimiento_front)

	logPiece, path, laser_end, err = tablero.ProcessTurn(movimiento_front, BLUE_TEAM)
	logPiece = logPiece + "%{0}"

	//Error movimiento inválido
	if err != nil {
		t.Error(err)
	} else {
		//Camino del laser
		tablero.printlaser(path)
		fmt.Println("Camino del laser:", formatLaserPath(path))
		fmt.Println("Terminación del laser:", tablero.blueTeamLaser.printLaserInteractionResult(laser_end))
		fmt.Print("RESP:\t" + logPiece + "\n")
		log = log + logPiece + ";"
	}

	//Ejemplo de procesamiento
	movimiento_front = "Lj5"

	fmt.Println("LLEGA:\t" + movimiento_front)

	logPiece, path, laser_end, err = tablero.ProcessTurn(movimiento_front, RED_TEAM)
	logPiece = logPiece + "%{0}"

	//Error movimiento inválido
	if err != nil {
		t.Error(err)
	} else {
		//Camino del laser
		tablero.printlaser(path)
		fmt.Println("Camino del laser:", formatLaserPath(path))
		fmt.Println("Terminación del laser:", tablero.redTeamLaser.printLaserInteractionResult(laser_end))
		fmt.Print("RESP:\t" + logPiece + "\n")
		log = log + logPiece + ";"
	}

	//Ejemplo de procesamiento
	movimiento_front = "Tc7:c8"

	fmt.Println("LLEGA:\t" + movimiento_front)

	logPiece, path, laser_end, err = tablero.ProcessTurn(movimiento_front, BLUE_TEAM)
	logPiece = logPiece + "%{0}"

	//Error movimiento inválido
	if err != nil {
		t.Error(err)
	} else {
		//Camino del laser
		tablero.printlaser(path)
		fmt.Println("Camino del laser:", formatLaserPath(path))
		fmt.Println("Terminación del laser:", tablero.blueTeamLaser.printLaserInteractionResult(laser_end))
		fmt.Print("RESP:\t" + logPiece + "\n")
		log = log + logPiece + ";"
	}

	//Ejemplo de procesamiento
	movimiento_front = "Ra8"

	fmt.Println("LLEGA:\t" + movimiento_front)

	logPiece, path, laser_end, err = tablero.ProcessTurn(movimiento_front, BLUE_TEAM)
	logPiece = logPiece + "%{0}"

	//Error movimiento inválido
	if err != nil {
		t.Error(err)
	} else {
		//Camino del laser
		tablero.printlaser(path)
		fmt.Println("Camino del laser:", formatLaserPath(path))
		fmt.Println("Terminación del laser:", tablero.blueTeamLaser.printLaserInteractionResult(laser_end))
		fmt.Print("RESP:\t" + logPiece + "\n")
		log = log + logPiece + ";"
	}

	//Ejemplo de procesamiento
	movimiento_front = "Ra4"

	fmt.Println("LLEGA:\t" + movimiento_front)

	logPiece, path, laser_end, err = tablero.ProcessTurn(movimiento_front, BLUE_TEAM)
	logPiece = logPiece + "%{0}"

	//Error movimiento inválido
	if err != nil {
		t.Error(err)
	} else {
		//Camino del laser
		tablero.printlaser(path)
		fmt.Println("Camino del laser:", formatLaserPath(path))
		fmt.Println("Terminación del laser:", tablero.blueTeamLaser.printLaserInteractionResult(laser_end))
		fmt.Print("RESP:\t" + logPiece + "\n")
		log = log + logPiece + ";"
	}

	//Ejemplo de procesamiento
	movimiento_front = "Rf4"

	fmt.Println("LLEGA:\t" + movimiento_front)

	logPiece, path, laser_end, err = tablero.ProcessTurn(movimiento_front, BLUE_TEAM)
	logPiece = logPiece + "%{0}"

	//Error movimiento inválido
	if err != nil {
		t.Error(err)
	} else {
		//Camino del laser
		tablero.printlaser(path)
		fmt.Println("Camino del laser:", formatLaserPath(path))
		fmt.Println("Terminación del laser:", tablero.blueTeamLaser.printLaserInteractionResult(laser_end))
		fmt.Print("RESP:\t" + logPiece + "\n")
		log = log + logPiece + ";"
	}

}

// Test que comprueba el correcto funcionamiento del cargado de un log
func TestApplyLogToBoard(t *testing.T) {
	var gameEngine GameEngine

	// Cargamos una partida "estado inicial" y "log"//
	gameEngine.InitEngine(CURIOSITY)
	gameEngine.gameLog = `Rf1%j1,j4,i4,i5,j5,j9%{300};Tg6:f6%a8,a5,b5,b4,a4,a0%{250};Rb4%j1,j4,i4,i5,j5,j9%{200};Ri5xf6%a8,a5,b5,b4,e4,e5,f5,f6%{150};Re4xf8%j1,j4,i4,i5,f5,f4,e4,e5,f5,f8%{100};`

	// Aplicamos el log al estado inicial
	team, redTimeLeft, blueTimeLeft := gameEngine.ApplyLogToBoard(400)

	// Correcto el turno siguiente?
	if team != BLUE_TEAM {
		t.Error("No se gestionan bien los turnos")
	}

	if redTimeLeft != 100 {
		t.Log(redTimeLeft)
		t.Error("No se recuperan bien los tiempos del rojo")
	}

	if blueTimeLeft != 150 {
		t.Log(blueTimeLeft)
		t.Error("No se recuperan bien los tiempos del azul")
	}

	// Correcto la muerte del rey
	switch gameEngine.gameBoard.cells[7][5].(type) {
	case *BoardPieceKing:
		t.Error("Partida mal cargada - rey vivo")
	case *BoardPieceVacant:
		//Resultado esperado
	default:
		t.Error("Partida mal cargada - pieza no esperada")

	}

}

func TestAImove(t *testing.T) {
	tablero, _ := InitBoard(boardtemplates.CURIOSITY)

	for i := 0; i < 150; i++ {
		switch i % 2 {
		case 0: //ROJO
			move := GetBestMove(tablero, RED_TEAM, 3)
			logPiece, laserPath, interactionResult, err := tablero.ProcessTurn(move, RED_TEAM)
			if err != nil {
				t.Error("la IA ha hecho un movimiento malo en RED_TEAM", err)
			}
			logPiece, laserPath, _, _ = tablero.calculateReturnValues(logPiece, laserPath, interactionResult)

			fmt.Println(logPiece)
			tablero.printlaser(laserPath)	
			
		case 1: //AZUL
			move := GetBestMove(tablero, BLUE_TEAM, 3)
			logPiece, laserPath, interactionResult, err := tablero.ProcessTurn(move, BLUE_TEAM)
			if err != nil {
				t.Error("la IA ha hecho un movimiento malo en BLUE_TEAM", err)
			}
			logPiece, laserPath, _, _ = tablero.calculateReturnValues(logPiece, laserPath, interactionResult)

			fmt.Println(logPiece)
			tablero.printlaser(laserPath)		
		}
	}
}