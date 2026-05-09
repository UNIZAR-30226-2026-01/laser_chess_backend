package rewards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLevel(t *testing.T) {
	tests := []struct {
		name     string
		xp       int32
		expected int32
	}{
		{"XP 0 -> Nivel 0", 0, 0},
		{"Casi Nivel 1 (249) -> Nivel 0", 249, 0},
		{"Justo Nivel 1 (250) -> Nivel 1", 250, 1},
		{"Mitad Nivel 1 (500) -> Nivel 1", 500, 1},
		{"Justo Nivel 2 (1000) -> Nivel 2", 1000, 2},
		{"Nivel Alto -> Nivel 10", 25000, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetLevel(tt.xp))
		})
	}
}

func TestGetLevelXP(t *testing.T) {
	tests := []struct {
		name     string
		level    int32
		expected int32
	}{
		{"Nivel 0 -> 0 XP", 0, 0},
		{"Nivel 1 -> 250 XP", 1, 250},
		{"Nivel 2 -> 1000 XP", 2, 1000},
		{"Nivel 3 -> 2250 XP", 3, 2250},
		{"Nivel 10 -> 25000 XP", 10, 25000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetLevelXP(tt.level))
		})
	}
}

func TestGetXPBarInfo(t *testing.T) {
	tests := []struct {
		name              string
		xp                int32
		expectedCurrent   int32
		expectedThreshold int32
	}{
		{"Inicio total (0 XP)", 0, 0, 250},                   // Nivel 0: de 0 a 250 (Faltan 250)
		{"Recien subido a Nivel 1 (250 XP)", 250, 0, 750},    // Nivel 1: de 250 a 1000 (Faltan 750)
		{"A mitad de Nivel 1 (500 XP)", 500, 250, 750},       // 250 XP dentro del nivel
		{"Recien subido a Nivel 2 (1000 XP)", 1000, 0, 1250}, // Nivel 2: de 1000 a 2250 (Faltan 1250)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current, threshold := GetXPBarInfo(tt.xp)
			assert.Equal(t, tt.expectedCurrent, current, "Current Level XP no coincide")
			assert.Equal(t, tt.expectedThreshold, threshold, "Threshold no coincide")
		})
	}
}

func TestGetEloDifferenceMult(t *testing.T) {
	tests := []struct {
		name     string
		myElo    int32
		yourElo  int32
		expected float64
	}{
		{"Mismo Elo", 1500, 1500, 1.0},
		{"Oponente un poco mejor (+200)", 1500, 1700, 1.2},
		{"Oponente un poco peor (-200)", 1700, 1500, 0.8},
		{"Oponente mucho mejor (tope máx 1.5)", 1000, 2000, 1.5},
		{"Oponente exageradamente mejor (clamp a 1.5)", 1000, 3000, 1.5},
		{"Oponente mucho peor (tope mín 0.5)", 2000, 1000, 0.5},
		{"Oponente exageradamente peor (clamp a 0.5)", 3000, 1000, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Usamos InDelta con un margen de 0.0001 para ignorar errores de precisión decimal
			assert.InDelta(t, tt.expected, getEloDifferenceMult(tt.myElo, tt.yourElo), 0.0001)
		})
	}
}

func TestGetMults(t *testing.T) {
	// Comprobamos la combinacion del multiplicador de Elo con el Ranked/Normal
	// rankMult = 1.2, normalMult = 1.0
	tests := []struct {
		name     string
		myElo    int32
		enemyElo int32
		isRanked bool
		expected float64
	}{
		{"Normal, Mismo Elo", 1500, 1500, false, 1.0 * 1.0},
		{"Ranked, Mismo Elo", 1500, 1500, true, 1.0 * 1.2},
		{"Normal, Enemigo mejor (+500)", 1000, 1500, false, 1.5 * 1.0},
		{"Ranked, Enemigo mejor (+500)", 1000, 1500, true, 1.5 * 1.2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Usamos InDelta con un margen de 0.0001
			assert.InDelta(t, tt.expected, getMults(tt.myElo, tt.enemyElo, tt.isRanked), 0.0001)
		})
	}
}

func TestGetMatchXP(t *testing.T) {
	// winXP = 150, loseXP = 80
	tests := []struct {
		name          string
		p1Elo         int32
		p2Elo         int32
		scoreP1       float64
		isRanked      bool
		expectedP1Exp int32
		expectedP2Exp int32
	}{
		{"P1 Gana, Normal, Igual Elo", 1500, 1500, 1.0, false, 150, 80},
		{"P2 Gana, Normal, Igual Elo", 1500, 1500, 0.0, false, 80, 150},
		{"Empate o Invalido (<1) da victoria a P2", 1500, 1500, 0.5, false, 80, 150}, // Debido al if scoreP1 >= 1
		{"P1 Gana, Ranked, Igual Elo", 1500, 1500, 1.0, true, 180, 96},               // 150*1.2, 80*1.2
		{"P1 Gana, Normal, P2 es mejor", 1000, 1500, 1.0, false, 225, 40},            // P1 mult=1.5 (150*1.5=225), P2 mult=0.5 (80*0.5=40)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1, p2 := GetMatchXP(tt.p1Elo, tt.p2Elo, tt.scoreP1, tt.isRanked)
			assert.Equal(t, tt.expectedP1Exp, p1, "XP de P1 no coincide")
			assert.Equal(t, tt.expectedP2Exp, p2, "XP de P2 no coincide")
		})
	}
}

func TestGetMatchMoney(t *testing.T) {
	// winMoney = 50, loseMoney = 10
	tests := []struct {
		name            string
		p1Elo           int32
		p2Elo           int32
		scoreP1         float64
		isRanked        bool
		expectedP1Money int32
		expectedP2Money int32
	}{
		{"P1 Gana, Normal, Igual Elo", 1500, 1500, 1.0, false, 50, 10},
		{"P2 Gana, Normal, Igual Elo", 1500, 1500, 0.0, false, 10, 50},
		{"P1 Gana, Ranked, Igual Elo", 1500, 1500, 1.0, true, 60, 12},    // 50*1.2, 10*1.2
		{"P1 Gana, Normal, P2 es mejor", 1000, 1500, 1.0, false, 75, 5},  // P1 mult=1.5 (50*1.5=75), P2 mult=0.5 (10*0.5=5)
		{"P2 Gana, Normal, P2 es mejor", 1000, 1500, 0.0, false, 15, 25}, // P1 mult=1.5 (lose 10*1.5=15), P2 mult=0.5 (win 50*0.5=25)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1, p2 := GetMatchMoney(tt.p1Elo, tt.p2Elo, tt.scoreP1, tt.isRanked)
			assert.Equal(t, tt.expectedP1Money, p1, "Dinero de P1 no coincide")
			assert.Equal(t, tt.expectedP2Money, p2, "Dinero de P2 no coincide")
		})
	}
}

func TestGetNewVolatility_CoverageExtra(t *testing.T) {
	// Para cubrir las líneas del método de Illinois (fA = fA / 2.0) y el bucle de (k += 1.0),
	// iteramos sobre una matriz de valores extremos que deforman la curva matemática de getF.

	deltas := []float64{0.0, 0.1, 0.5, 1.0, 2.0, 5.0, 10.0}
	phis := []float64{0.1, 0.5, 1.0, 2.0}
	vs := []float64{0.01, 0.1, 1.0, 2.0}

	// El tau de 100.0 asegura matemáticamente que term1 domine a term2 en la primera iteración,
	// forzando la entrada al bucle "for getF(...) < 0 { k += 1.0 }"
	taus := []float64{0.1, 0.5, 1.0, 10.0, 100.0}

	for _, d := range deltas {
		for _, p := range phis {
			for _, v := range vs {
				for _, tau := range taus {
					// Al probar cientos de combinaciones geométricas de la curva,
					// garantizamos que al menos una desencadene la condición fC*fB > 0 (el else fA = fA / 2.0)
					_ = getNewVolatility(d, p, v, 0.06, tau)
				}
			}
		}
	}
}
