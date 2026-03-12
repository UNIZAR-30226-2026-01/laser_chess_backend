package rt

import (
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

// fichero que gestiona las rooms
// una room se encarga de gestionar los mensajes de los dos
// jugadores de una partida
// es el intermediario entre el front y el juego

type Room struct {
	Player1 *Client
	Player2 *Client
	Game *game.Board

	ConP1 chan interface{}
	ConP2 chan interface{}
	Broadcast chan interface{}
	
	// tendra mas cosas para hablar con el juego
	// para mandar movimientos y recibir info que
	// mandar a los clientes
	// probablemente unos canales para ida y otros para vuelta o algo

}

//TODO: todo

func (r *Room) InitRoom (Player1 *Client, Player2 *Client) {
	r.Player1 = Player1
	r.Player2 = Player2
	r.ConP1 = make(chan interface{})
	r.ConP2 = make(chan interface{})
	r.Broadcast = make(chan interface{})

	// r.Game lo que sea pero no se puede iniciar aun
}

func (r *Room) GetGameState () {
	// Falta funcion del juego para
	// devolver estado, o lo guardamos en room
}

func (r *Room) Run() {
	for {
		select {
		case message := <-r.ConP1:
			r.Player2.Send <- message
			
		case message := <-r.ConP2:
			r.Player2.Send <- message

		case message := <-r.Broadcast:
			r.Player1.Send <- message
			r.Player2.Send <- message
		} 	// AÑADIR OTRO CASE PARA MODELAR LA COMUNICACION CON EL JUEGO
			// POR CANALES
	}
}

func (r *Room) SendMoveToGame(move string) (string, error) {
	// resul, err := r.Game.ProcessTurn(move)
	// return resul, err
}

//HACER QUE DEVUELVA UN ERROR EN VEZ DE UN BOOLEANO
func (r *Room) MakeMove(AccountID int64, move string) bool {
	
	if AccountID != r.Player1.AccountID && AccountID != r.Player2.AccountID {
		return false
	}

	state, err := r.SendMoveToGame(move)
	if err != nil {
		return false
	}
	r.Broadcast <- state

}