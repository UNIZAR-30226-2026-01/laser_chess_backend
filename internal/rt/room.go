package rt

// fichero que gestiona las rooms
// una room se encarga de gestionar los mensajes de los dos
// jugadores de una partida
// es el intermediario entre el front y el juego

type Room struct {
	Player1 *Client
	Player2 *Client

	// tendra mas cosas para hablar con el juego
	// para mandar movimientos y recibir info que
	// mandar a los clientes
	// probablemente unos canales para ida y otros para vuelta o algo

}

//TODO: todo
