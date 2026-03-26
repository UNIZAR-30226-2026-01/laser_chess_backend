package game

import (
	"fmt"
	"testing"
)

func TestAllBoards(t *testing.T) {
	tablero, err := InitBoard("boardTemplates/ace.csv")
	if err != nil {
		t.Error(err)
	}

	fmt.Print("== ACE ==\n")

	tablero.print()

	fmt.Print("== CURIOSITY ==\n")

	tablero, err = InitBoard("boardTemplates/curiosity.csv")
	if err != nil {
		t.Error(err)
	}
	tablero.print()

	fmt.Print("== GRAIL ==\n")
	tablero, err = InitBoard("boardTemplates/curiosity.grail")
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
	tablero, err := InitBoard("boardTemplates/curiosity.csv")
	if err != nil {
		t.Error(err)
	}
	tablero.print()

	tablero.print()

	var log string
	var logPiece string
	var path []vector2_T

	//Ejemplo de procesamiento
	movimiento_front := "La8"

	fmt.Println("LLEGA:\t" + movimiento_front)

	logPiece, path, _, err = tablero.ProcessTurn(movimiento_front, BLUE_TEAM)
	logPiece = logPiece + "%{0}"

	//Error movimiento inválido
	if err != nil {
		t.Error(err)
	} else {
		//Camino del laser
		tablero.printlaser(path)
		fmt.Print("RESP:\t" + logPiece + "\n")
		log = log + logPiece + ";"
	}

	//Ejemplo de procesamiento
	movimiento_front = "Lh2"

	fmt.Println("LLEGA:\t" + movimiento_front)

	logPiece, path, _, err = tablero.ProcessTurn(movimiento_front, RED_TEAM)
	logPiece = logPiece + "%{0}"

	//Error movimiento inválido
	if err != nil {
		t.Error(err)
	} else {
		//Camino del laser
		tablero.printlaser(path)
		fmt.Print("RESP:\t" + logPiece + "\n")
		log = log + logPiece + ";"
	}

	//Ejemplo de procesamiento
	movimiento_front = "Tc7:c8"

	fmt.Println("LLEGA:\t" + movimiento_front)

	logPiece, path, _, err = tablero.ProcessTurn(movimiento_front, BLUE_TEAM)
	logPiece = logPiece + "%{0}"

	//Error movimiento inválido
	if err != nil {
		t.Error(err)
	} else {
		//Camino del laser
		tablero.printlaser(path)
		fmt.Print("RESP:\t" + logPiece + "\n")
		log = log + logPiece + ";"
	}
}
