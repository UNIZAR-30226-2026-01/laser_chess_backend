package game

// Archivo que implementa los timers que se van a usar en partida

import (
	"sync"
	"time"
)

type GameTimer struct {
	Remaining time.Duration
	Increment time.Duration
	Expired   chan bool

	lastStartedAt time.Time
	isRunning     bool

	// (canal para cancelar la espera si se para
	// el timer antes de que acabe)
	stop chan struct{}

	mu sync.Mutex
}

func NewGameTimer(initial, inc time.Duration) *GameTimer {
	return &GameTimer{
		Remaining: initial,
		Increment: inc,
		Expired:   make(chan bool, 1),
	}
}

// gorutina que mira si el timer acaba mientras está
// activo, y avisa por el canal Expired
func (t *GameTimer) timerSurveillance(remaining time.Duration, stop chan struct{}) {
	timer := time.NewTimer(remaining)
	defer timer.Stop()

	select {
	case <-timer.C:
		// Si acaba el timer se avisa
		t.Expired <- true
	case <-stop:
		// Si para antes no se hace nada
	}
}

// Activa el timer
func (t *GameTimer) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.isRunning {
		t.lastStartedAt = time.Now()
		t.isRunning = true

		t.stop = make(chan struct{})
		go t.timerSurveillance(t.Remaining, t.stop)
	}
}

// Para el timer y añade incremento
func (t *GameTimer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.isRunning {
		close(t.stop)

		t.Remaining -= time.Since(t.lastStartedAt)
		t.Remaining += t.Increment
		t.isRunning = false

		// Evitar race condition si expira justo antes de stop
		select {
		case <-t.Expired:
		default:
		}
	}
}
