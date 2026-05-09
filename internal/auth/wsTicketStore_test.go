package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWsTicketStore_StoreAndGet(t *testing.T) {
	store := NewWsTicketStore()
	ticket := "ticket-secreto-123"
	accountID := int64(456)

	// Guardamos un ticket válido por 1 minuto
	store.Store(ticket, accountID, 1*time.Minute)

	// 1. Primer intento: el ticket debe ser válido y devolver el accountID correcto
	id, valid := store.GetAndDelete(ticket)
	assert.True(t, valid, "El ticket debería ser válido")
	assert.Equal(t, accountID, id, "El AccountID no coincide")

	// 2. Segundo intento: el ticket ya no debe existir porque GetAndDelete lo borra
	_, valid2 := store.GetAndDelete(ticket)
	assert.False(t, valid2, "El ticket no debería poder usarse dos veces")
}

func TestWsTicketStore_GetNonExistent(t *testing.T) {
	store := NewWsTicketStore()
	_, valid := store.GetAndDelete("ticket-inventado")
	assert.False(t, valid, "Un ticket que no existe debe devolver false")
}

func TestWsTicketStore_ExpiredTicket(t *testing.T) {
	store := NewWsTicketStore()
	ticket := "ticket-caducado"

	// Lo guardamos con un TTL negativo para que nazca ya caducado
	store.Store(ticket, 789, -1*time.Second)

	// Intentamos usarlo
	_, valid := store.GetAndDelete(ticket)
	assert.False(t, valid, "El ticket no debería ser válido porque ha caducado")
}

func TestWsTicketStore_Cleanup(t *testing.T) {
	store := NewWsTicketStore()

	// Guardamos uno válido y uno caducado
	store.Store("ticket-bueno", 1, 1*time.Minute)
	store.Store("ticket-malo", 2, -1*time.Second)

	// Ejecutamos la función privada de limpieza que usa la goroutine internamente
	store.cleanup()

	// Comprobamos el estado interno del mapa
	store.mu.Lock()
	_, existBueno := store.tickets["ticket-bueno"]
	_, existMalo := store.tickets["ticket-malo"]
	store.mu.Unlock()

	assert.True(t, existBueno, "El ticket válido debe sobrevivir a la limpieza")
	assert.False(t, existMalo, "El ticket caducado debe ser eliminado en la limpieza")
}
