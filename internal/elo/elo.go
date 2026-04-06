package elo

// Paquete que se encarga de realizar los cálculos del elo
// Implementa el algoritmo de Glicko 2
// Crea una función por cada función que utiliza el algoritmo
// Para más información, ver el paper original

// ApplyDeviation - Aumenta la Deviation de un player si ha pasado tiempo sin jugar
// ProcessMatch - calcula los nuevos ratings para dos players a la vez

import (
	"math"
)

const (
	DefaultRating     = 1500.0
	DefaultDeviation  = 350
	DefaultVolatility = 0.06
	DefaultTau        = 0.5
	Multiplier        = 173.7178
	Tolerance         = 0.000001

	PeriodsPerDay = 1.0 / 7.0 // Un periodo es una semana

	PiSquared = math.Pi * math.Pi
)

type Rating struct {
	Value      float64
	Deviation  float64
	Volatility float64
}

// Paso 1 - Aumenta la Deviation de un player si ha pasado tiempo sin jugar
func ApplyInactivity(playerRating Rating, inactiveDays float64) Rating {
	if inactiveDays <= 0 {
		return playerRating
	}

	pastPeriods := inactiveDays * PeriodsPerDay
	phi := getPhi(playerRating.Deviation)

	newPhi := math.Sqrt(math.Pow(phi, 2) + pastPeriods*math.Pow(playerRating.Volatility, 2))
	newDeviation := newPhi * Multiplier

	// RD no puede ser mayor que la de un jugador nuevo
	if newDeviation > DefaultDeviation {
		newDeviation = DefaultDeviation
	}

	playerRating.Deviation = newDeviation
	return playerRating
}

// Paso 2
func getMu(rating float64) float64 {
	return (rating - 1500.0) / Multiplier
}

// Paso 2
func getPhi(deviation float64) float64 {
	return deviation / Multiplier
}

func getG(phi float64) float64 {
	return 1.0 / math.Sqrt(1.0+3.0*math.Pow(phi, 2)/PiSquared)
}

func getE(mu, muJ, phiJ float64) float64 {
	return 1.0 / (1.0 + math.Exp(-getG(phiJ)*(mu-muJ)))
}

// Paso 3
func getVariance(mu, muJ, phiJ float64) float64 {
	gPhiJ := getG(phiJ)
	e := getE(mu, muJ, phiJ)

	inverse := math.Pow(gPhiJ, 2) * e * (1.0 - e)
	return 1.0 / inverse
}

// Paso 4
func getDelta(v, mu, muJ, phiJ, score float64) float64 {
	gPhiJ := getG(phiJ)
	e := getE(mu, muJ, phiJ)

	return v * (gPhiJ * (score - e))
}

func getF(x, delta, phi, v, a, tau float64) float64 {
	expX := math.Exp(x)
	numerator := expX * (math.Pow(delta, 2) - math.Pow(phi, 2) - v - expX)
	denominator := 2.0 * math.Pow(math.Pow(phi, 2)+v+expX, 2)

	term1 := numerator / denominator
	term2 := (x - a) / math.Pow(tau, 2)

	return term1 - term2
}

// Paso 5
func getNewVolatility(delta, phi, v, volatility, tau float64) float64 {
	a := math.Log(math.Pow(volatility, 2))
	A := a

	var B float64
	if math.Pow(delta, 2) > math.Pow(phi, 2)+v {
		B = math.Log(math.Pow(delta, 2) - math.Pow(phi, 2) - v)
	} else {
		k := 1.0
		for getF(a-k*tau, delta, phi, v, a, tau) < 0 {
			k += 1.0
		}
		B = a - k*tau
	}

	fA := getF(A, delta, phi, v, a, tau)
	fB := getF(B, delta, phi, v, a, tau)

	for math.Abs(B-A) > Tolerance {
		C := A + (A-B)*fA/(fB-fA)
		fC := getF(C, delta, phi, v, a, tau)

		if fC*fB <= 0 {
			A = B
			fA = fB
		} else {
			fA = fA / 2.0
		}
		B = C
		fB = fC
	}

	return math.Exp(A / 2.0)
}

// Pasos 6 y 7
func getNewMuPhi(mu, phi, v, newVolatility, muJ, phiJ, score float64) (float64, float64) {
	// Paso 6
	phiStar := math.Sqrt(math.Pow(phi, 2) + math.Pow(newVolatility, 2))

	// Paso 7
	newPhi := 1.0 / math.Sqrt((1.0/math.Pow(phiStar, 2))+(1.0/v))

	gPhiJ := getG(phiJ)
	e := getE(mu, muJ, phiJ)

	newMu := mu + math.Pow(newPhi, 2)*(gPhiJ*(score-e))

	return newMu, newPhi
}

// Evalua a un jugador contra su oponente y devuelve su nuevo struct Rating
func updatePlayerRating(playerRating Rating, oponentRating Rating, score float64) Rating {
	// Paso 2
	mu := getMu(playerRating.Value)
	phi := getPhi(playerRating.Deviation)
	sigma := playerRating.Volatility

	muJ := getMu(oponentRating.Value)
	phiJ := getPhi(oponentRating.Deviation)

	// Paso 3 y 4
	v := getVariance(mu, muJ, phiJ)
	delta := getDelta(v, mu, muJ, phiJ, score)

	// Paso 5
	newSigma := getNewVolatility(delta, phi, v, sigma, DefaultTau)

	// Paso 6 y 7
	newMu, newPhi := getNewMuPhi(mu, phi, v, newSigma, muJ, phiJ, score)

	// Paso 8
	return Rating{
		Value:      newMu*Multiplier + 1500.0,
		Deviation:  newPhi * Multiplier,
		Volatility: newSigma,
	}
}

// Calcula los nuevos ratings para dos players a la vez
func ProcessMatch(p1Rating Rating, p2Rating Rating, scoreP1 float64) (Rating, Rating) {
	scoreP2 := 1.0 - scoreP1

	newP1Rating := updatePlayerRating(p1Rating, p2Rating, scoreP1)
	newP2Rating := updatePlayerRating(p2Rating, p1Rating, scoreP2)

	return newP1Rating, newP2Rating
}
