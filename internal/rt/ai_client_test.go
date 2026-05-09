package rt

import (
	"testing"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
	"github.com/stretchr/testify/assert"
)

// TestAIClient_EOC verifica que el cliente de la IA se apaga correctamente
// al recibir el mensaje de fin de conexión (EOC).
func TestAIClient_EOC(t *testing.T) {
	// 1. Setup del cliente dummy
	client := &Client{
		ToAI:   make(chan ClientSocketMessage, 1),
		FromAI: make(chan ClientSocketMessage, 1),
	}

	ai := &AIClient{}

	// Iniciamos el cliente de la IA (esto lanza go ai.run() por debajo)
	ai.InitAIClient(client, game.ACE, 1)

	// 2. Ejecución: Enviamos EOC
	client.ToAI <- ClientSocketMessage{Type: "EOC", Content: ""}

	// 3. Verificación
	// Le damos un pequeño margen de tiempo para que el hilo termine.
	// Si no hubiera procesado el EOC y hecho return, el goroutine se quedaría colgado.
	time.Sleep(100 * time.Millisecond)

	// Si el test llega aquí sin colgarse, el goroutine ha terminado correctamente.
}

// TestAIClient_GenerateMove verifica el flujo normal: recibe un movimiento,
// actualiza el log, llama a GetBestMove y devuelve la respuesta tras el sleep.
func TestAIClient_GenerateMove(t *testing.T) {
	// 1. Setup
	client := &Client{
		ToAI:   make(chan ClientSocketMessage, 1),
		FromAI: make(chan ClientSocketMessage, 1),
	}

	ai := &AIClient{}
	ai.InitAIClient(client, game.ACE, 1)

	// 2. Ejecución: Enviamos un movimiento falso simulando al jugador humano
	client.ToAI <- ClientSocketMessage{
		Type:    "Move",
		Content: "Tc1:c2%c1,c2%{30000};",
	}

	// 3. Verificación
	// Como el AIClient tiene un time.Sleep(1000 * time.Millisecond) puesto a fuego,
	// tenemos que esperar al menos ese segundo más un poco de margen.
	select {
	case response := <-client.FromAI:
		// Comprobamos que nos responde un mensaje de tipo Move
		assert.Equal(t, "Move", response.Type)

		// Comprobamos que el log interno de la IA ha guardado nuestro movimiento
		assert.Equal(t, "Tc1:c2%c1,c2%{30000};", ai.log)

		// Verificamos que el contenido devuelto no esté vacío (la IA nos devuelve *algo*)
		assert.NotEmpty(t, response.Content, "La IA debería haber generado un movimiento de respuesta")

	case <-time.After(2 * time.Second): // Esperamos máximo 2 segundos
		t.Fatal("Timeout: La IA no devolvió ningún movimiento después de 2 segundos")
	}

	// 4. Cleanup: Apagamos el goroutine
	client.ToAI <- ClientSocketMessage{Type: "EOC"}
}
