package rt

import (
	"fmt"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

// Clase intermediaria para la comunicacion entre el cliente y la IA para las partidas
// contra bots

type AIClient struct {
	client *Client
	board  game.Board_T
	lvl    int
	log    string
}

func (ai *AIClient) InitAIClient(client *Client, board game.Board_T, lvl int) {
	ai.client = client
	ai.board = board
	ai.lvl = lvl
	ai.log = ""

	go ai.run()
}

// lee mensajes del socket y los manda a la Room
func (ai *AIClient) run() {
	for {
		fmt.Println("Esperando movimento")
		message := <-ai.client.ToAI
		ai.log += message.Content
		fmt.Println("Movimiento recibido")
		if message.Type == "EOC" {
			return
		}
		fmt.Println("MESSAGE TYPE RECIBIDO DE IA: ", message.Type)
		fmt.Println("Llamada a GetBestMove con board: ", ai.board, ", log: ", ai.log, ", y level: ", ai.lvl)
		move := game.GetBestMove(ai.board, ai.log, ai.lvl)

		// Tiempo para no saturar al usuario
		time.Sleep(1000 * time.Millisecond)

		ai.client.FromAI <- ClientSocketMessage{
			Type:    "Move",
			Content: move,
		}

		log := <-ai.client.ToAI
		ai.log += log.Content

	}
}
