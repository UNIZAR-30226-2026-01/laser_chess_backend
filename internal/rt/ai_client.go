package rt

import (
	"fmt"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

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
		fmt.Println("Llamada a GetBestMove con board: ", ai.board, ", log: ", ai.log, ", y level: ", ai.lvl)
		move := game.GetBestMove(ai.board, ai.log, ai.lvl)
		ai.log += move

		// Pausa
		time.Sleep(1 * time.Second)

		ai.client.FromAI <- ClientSocketMessage{
			Type:    "Move",
			Content: move,
		}
		fmt.Println("Fin de la llamada a GetBestMove")

	}
}
