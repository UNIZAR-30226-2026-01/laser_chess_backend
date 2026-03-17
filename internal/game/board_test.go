package game

import (
	"fmt"
	"testing"
)

// func TestProbatTipoDeDato(t *testing.T) {
// 	tablero := Board{}

// 	//Iniciar tablero
// 	initACE(&tablero)

// 	// === Rey === //
// 	t.Log("MOVIMIENTO")

// 	t.Log("== Rey ==")

// 	if tablero.ProcessTurn("Te8:e6", BLUE_TEAM) != false {
// 		t.Errorf("X - Se ha aceptado un movimiento fuera del alcance de la pieza")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	if tablero.ProcessTurn("Te8:e9") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento fuera del tablero")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	if tablero.ProcessTurn("Te8:d8") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento a una casilla ocupada")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	tablero.cells[4][6] = &BoardPieceVacant{BLUE_TEAM}
// 	if tablero.ProcessTurn("Te8:e7") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento a una casilla del equipo opuesto")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}
// 	tablero.cells[4][6] = &BoardPieceVacant{NONE}

// 	if tablero.ProcessTurn("Te8:e7") != true {
// 		t.Errorf("X - Se ha rechazado un movimiento válido")
// 	} else {
// 		initACE(&tablero)
// 		t.Log("OK")
// 	}

// 	tablero.cells[4][6] = &BoardPieceVacant{RED_TEAM}
// 	if tablero.ProcessTurn("Te8:e7") != true {
// 		t.Errorf("X - Se ha rechazado un movimiento válido")
// 	} else {
// 		initACE(&tablero)
// 		t.Log("OK")
// 	}
// 	tablero.cells[4][6] = &BoardPieceVacant{NONE}

// 	t.Log("== Escudo ==")

// 	if tablero.ProcessTurn("Td8:d6") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento fuera del alcance de la pieza")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	if tablero.ProcessTurn("Td8:d9") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento fuera del tablero")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	if tablero.ProcessTurn("Tf8:e8") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento a una casilla ocupada")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	tablero.cells[3][6] = &BoardPieceVacant{BLUE_TEAM}
// 	if tablero.ProcessTurn("Td8:d7") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento a una casilla del equipo opuesto")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}
// 	tablero.cells[3][6] = &BoardPieceVacant{NONE}

// 	if tablero.ProcessTurn("Td8:c7") != true {
// 		t.Errorf("X - Se ha rechazado un movimiento válido")
// 	} else {
// 		initACE(&tablero)
// 		t.Log("OK")
// 	}

// 	tablero.cells[3][6] = &BoardPieceVacant{RED_TEAM}
// 	if tablero.ProcessTurn("Td8:c7") != true {
// 		t.Errorf("X - Se ha rechazado un movimiento válido")
// 	} else {
// 		initACE(&tablero)
// 		t.Log("OK")
// 	}
// 	tablero.cells[3][6] = &BoardPieceVacant{NONE}

// 	t.Log("== Deflector ==")

// 	if tablero.ProcessTurn("Tc8:c6") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento fuera del alcance de la pieza")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	if tablero.ProcessTurn("Tc8:c9") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento fuera del tablero")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	if tablero.ProcessTurn("Th4:h5") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento a una casilla ocupada")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	if tablero.ProcessTurn("Tc2:b1") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento a una casilla del equipo opuesto")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	if tablero.ProcessTurn("Ta4:b3") != true {
// 		t.Errorf("X - Se ha rechazado un movimiento válido")
// 	} else {
// 		initACE(&tablero)
// 		t.Log("OK")
// 	}

// 	tablero.cells[3][6] = &BoardPieceVacant{RED_TEAM}
// 	if tablero.ProcessTurn("Ta4:a3") != true {
// 		t.Errorf("X - Se ha rechazado un movimiento válido")
// 	} else {
// 		initACE(&tablero)
// 		t.Log("OK")
// 	}
// 	tablero.cells[3][6] = &BoardPieceVacant{NONE}

// 	t.Log("== Laser ==")

// 	if tablero.ProcessTurn("Ta1:a2") != false {
// 		t.Errorf("X - Esta pieza siempre debería devolver false, (rarete)")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	t.Log("== Switch ==")

// 	if tablero.ProcessTurn("Te5:g7") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento fuera del alcance de la pieza")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	tablero.cells[6][7] = &BoardPieceSwitch{BLUE_TEAM, NONE, DOWN}
// 	if tablero.ProcessTurn("Tg8:g9") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento fuera del tablero")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}
// 	tablero.cells[6][7] = &BoardPieceVacant{NONE}

// 	tablero.cells[6][7] = &BoardPieceSwitch{BLUE_TEAM, NONE, DOWN}
// 	if tablero.ProcessTurn("Te5:f4") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento a una casilla ocupada (permutación no posible)")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}
// 	tablero.cells[6][7] = &BoardPieceVacant{NONE}

// 	tablero.cells[7][7] = &BoardPieceSwitch{RED_TEAM, NONE, DOWN}
// 	if tablero.ProcessTurn("Th8:i8") != false {
// 		t.Errorf("X - Se ha aceptado un movimiento a una casilla del equipo opuesto")
// 		initACE(&tablero)
// 	} else {
// 		t.Log("OK")
// 	}

// 	tablero.cells[7][7] = &BoardPieceSwitch{RED_TEAM, NONE, DOWN}
// 	if tablero.ProcessTurn("Th8:h7") != true {
// 		t.Errorf("X - Se ha rechazado un movimiento válido")
// 	} else {
// 		initACE(&tablero)
// 		t.Log("OK")
// 	}

// 	tablero.cells[6][7] = &BoardPieceSwitch{BLUE_TEAM, NONE, DOWN}
// 	if tablero.ProcessTurn("Tg8:f8") != true {
// 		t.Errorf("X - Se ha rechazado un movimiento válido (posible error permutación)")

// 	} else {
// 		t.Log("OK")
// 	}
// 	tablero.cells[6][7] = &BoardPieceVacant{NONE}

// 	// Test del recorrido del laser
// 	t.Log("Solo dispara el laser y debería finalizar por out of bounds")
// 	positions, terminationReason := tablero.blueTeamLaser.shootLaser(0, 0, &tablero)
// 	tablero.printlaser(positions)
// 	if terminationReason != OUT_OF_BOUNDS {
// 		t.Errorf("X - No ha terminado por la razón correcta, se esperaba OUT_OF_BOUNDS, recibido %d", terminationReason)
// 	} else {
// 		t.Log("OK")
// 	}

// 	t.Log(tablero.ProcessTurn("Ra1"))

// 	t.Log(tablero.ProcessTurn("La1"))
// 	positions, terminationReason = tablero.blueTeamLaser.shootLaser(0, 0, &tablero)
// 	tablero.printlaser(positions)
// 	t.Log(tablero.ProcessTurn("Ra1"))
// 	positions, terminationReason = tablero.blueTeamLaser.shootLaser(0, 0, &tablero)
// 	tablero.printlaser(positions)

// }

func TestMovements(t *testing.T) {
	tablero := Board{}

	//Iniciar tablero
	initACE(&tablero)

	fmt.Print("== TEST TRANSFORMACIONES ==\n")
	tablero.print()

	var log string

	//Ejemplo de procesamiento
	logPiece , path, _ , err := tablero.ProcessTurn("La1", BLUE_TEAM)
	if err != nil{
		t.Error(err)
	}
	log = log + ";" + logPiece + "%T{...}"
	fmt.Println(log)
	tablero.printlaser(path)

	//Ejemplo de procesamiento
	logPiece , path, _ , err = tablero.ProcessTurn("Rf8", RED_TEAM)
	if err != nil{
		t.Error(err)
	}
	log = log + ";" + logPiece + "%T{...}"
	fmt.Println(log)
	tablero.printlaser(path)

	//Ejemplo de procesamiento
	logPiece , path, _ , err = tablero.ProcessTurn("Tc2:c1", BLUE_TEAM)
	if err != nil{
		t.Error(err)
	}
	log = log + ";" + logPiece + "%T{...}"
	fmt.Println(log)
	tablero.printlaser(path)

	//Ejemplo de procesamiento
	logPiece , path, _ , err = tablero.ProcessTurn("Rc5", RED_TEAM)
	if err != nil{
		t.Error(err)
	}
	log = log + ";" + logPiece + "%T{...}"
	fmt.Println(log)
	tablero.printlaser(path)

		//Ejemplo de procesamiento
	logPiece , path, _ , err = tablero.ProcessTurn("Tg1:f2", BLUE_TEAM)
	if err != nil{
		t.Error(err)
	}
	log = log + ";" + logPiece + "%T{...}"
	fmt.Println(log)
	tablero.printlaser(path)

	//Ejemplo de procesamiento
	logPiece , path, _ , err = tablero.ProcessTurn("Td3:d2", RED_TEAM)
	if err != nil{
		t.Error(err)
	}
	log = log + ";" + logPiece + "%T{...}"
	fmt.Println(log)
	tablero.printlaser(path)


	//Ejemplo de procesamiento
	logPiece , path, _ , err = tablero.ProcessTurn("Re4", BLUE_TEAM)
	if err != nil{
		t.Error(err)
	}
	log = log + ";" + logPiece + "%T{...}"
	fmt.Println(log)
	tablero.printlaser(path)

	//Ejemplo de procesamiento
	logPiece , path, _ , err = tablero.ProcessTurn("Td2:d1", RED_TEAM)
	if err != nil{
		t.Error(err)
	}
	log = log + ";" + logPiece + "%T{...}"
	fmt.Println(log)
	tablero.printlaser(path)

		//Ejemplo de procesamiento
	logPiece , path, _ , err = tablero.ProcessTurn("Th1:i1", BLUE_TEAM)
	if err != nil{
		t.Error(err)
	}
	log = log + ";" + logPiece + "%T{...}"
	fmt.Println(log)
	tablero.printlaser(path)


	
}