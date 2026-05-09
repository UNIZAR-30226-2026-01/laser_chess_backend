package elo

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplyInactivity(t *testing.T) {
	baseRating := Rating{
		Value:      1500.0,
		Deviation:  200.0,
		Volatility: 0.06,
	}

	tests := []struct {
		name              string
		daysInactive      float64
		expectedDeviation float64
	}{
		{
			name:              "Sin inactividad (0 días) - RD no cambia",
			daysInactive:      0,
			expectedDeviation: 200.0,
		},
		{
			name:              "Inactividad negativa (futuro) - RD no cambia",
			daysInactive:      -5,
			expectedDeviation: 200.0,
		},
		{
			name:              "Poca inactividad (7 días / 1 periodo) - RD aumenta ligeramente",
			daysInactive:      7,
			expectedDeviation: 200.2714, // Calculado aproximado para 1 periodo con Vol=0.06
		},
		{
			name:              "Inactividad exagerada (10000 días) - RD se capa al Default (350)",
			daysInactive:      10000, // Con 10000 días nos aseguramos de superar el límite de 350
			expectedDeviation: DefaultDeviation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastUpdatedAt := time.Now().Add(-time.Duration(tt.daysInactive*24) * time.Hour)

			newRating := ApplyInactivity(baseRating, lastUpdatedAt)

			assert.Equal(t, baseRating.Value, newRating.Value, "El valor del ELO no debe cambiar por inactividad")
			assert.Equal(t, baseRating.Volatility, newRating.Volatility, "La volatilidad no debe cambiar por inactividad")
			assert.InDelta(t, tt.expectedDeviation, newRating.Deviation, 0.1, "La desviación no coincide con lo esperado")
		})
	}
}

func TestProcessMatch(t *testing.T) {
	p1 := Rating{Value: 1500, Deviation: 200, Volatility: 0.06}
	p2 := Rating{Value: 1400, Deviation: 150, Volatility: 0.06} // Subimos la RD inicial a un valor realista (150)

	t.Run("P1 Gana (Score 1.0)", func(t *testing.T) {
		newP1, newP2 := ProcessMatch(p1, p2, 1.0)

		// P1 debería subir de ELO y P2 bajar
		assert.Greater(t, newP1.Value, p1.Value, "P1 gana, su ELO debe subir")
		assert.Less(t, newP2.Value, p2.Value, "P2 pierde, su ELO debe bajar")

		// La desviación de ambos debería disminuir al tener valores iniciales altos/normales
		assert.Less(t, newP1.Deviation, p1.Deviation, "La RD de P1 debe bajar tras jugar")
		assert.Less(t, newP2.Deviation, p2.Deviation, "La RD de P2 debe bajar tras jugar")
	})

	t.Run("P1 Pierde (Score 0.0)", func(t *testing.T) {
		newP1, newP2 := ProcessMatch(p1, p2, 0.0)

		assert.Less(t, newP1.Value, p1.Value, "P1 pierde, su ELO debe bajar")
		assert.Greater(t, newP2.Value, p2.Value, "P2 gana, su ELO debe subir")
	})

	t.Run("Empate (Score 0.5)", func(t *testing.T) {
		newP1, newP2 := ProcessMatch(p1, p2, 0.5)

		// Al empatar contra alguien de menor ELO, P1 debería perder un poco de ELO y P2 ganar
		assert.Less(t, newP1.Value, p1.Value, "P1 empata contra ELO inferior, debe bajar")
		assert.Greater(t, newP2.Value, p2.Value, "P2 empata contra ELO superior, debe subir")
	})
}

func TestGetNewVolatility_Branches(t *testing.T) {
	t.Run("Sorpresa Masiva (Fuerza Rama B > a)", func(t *testing.T) {
		p1 := Rating{Value: 2500, Deviation: 50, Volatility: 0.06}
		p2 := Rating{Value: 500, Deviation: 50, Volatility: 0.06}

		newP1, _ := ProcessMatch(p1, p2, 0.0)
		assert.NotZero(t, newP1.Volatility)
	})

	t.Run("Resultado Esperado (Fuerza Rama Else en B)", func(t *testing.T) {
		p1 := Rating{Value: 1500, Deviation: 100, Volatility: 0.06}
		p2 := Rating{Value: 1500, Deviation: 100, Volatility: 0.06}

		newP1, _ := ProcessMatch(p1, p2, 0.5)
		assert.NotZero(t, newP1.Volatility)
	})
}

func TestMathHelpers(t *testing.T) {
	mu := getMu(1500.0)
	assert.Equal(t, 0.0, mu, "Un rating de 1500 debe tener un mu de 0.0")

	phi := getPhi(DefaultDeviation)
	assert.InDelta(t, 2.01476, phi, 0.0001, "El phi por defecto debería ser cercano a 2.014")

	gValue := getG(phi)
	// Valor corregido según la fórmula matemática
	assert.InDelta(t, 0.669069, gValue, 0.0001, "Cálculo matemático de G no coincide")

	delta := 0.5
	v := 0.2
	a := math.Log(math.Pow(0.06, 2))
	tau := 0.5
	fVal := getF(a, delta, phi, v, a, tau)
	assert.NotPanics(t, func() { getF(a, delta, phi, v, a, tau) })
	assert.NotZero(t, fVal)

	newMu, newPhi := getNewMuPhi(mu, phi, v, 0.06, mu, phi, 1.0)
	assert.NotZero(t, newMu)
	assert.NotZero(t, newPhi)
}
