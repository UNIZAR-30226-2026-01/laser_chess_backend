package xp

import "math"

// La formula para los niveles es scale * nivel^2

const (
	scale  = 250.0
	winXP  = 150.0
	loseXP = 80.0
)

func GetLevel(XP int32) int32 {
	level := math.Sqrt(float64(XP) / scale)
	return int32(math.Floor(level))
}

func GetLevelXP(level int32) int32 {
	return int32(math.Floor(250.0 * float64(math.Sqrt(float64(level)))))
}

func GetMatchXP(p1Elo int32, p2Elo int32, scoreP1 int) (int32, int32) {
	p1Mult := getEloDifferenceMult(p1Elo, p2Elo)
	p2Mult := getEloDifferenceMult(p2Elo, p1Elo)
	if scoreP1 >= 1 {
		// Gana P1
		p1Exp := math.Floor(winXP * p1Mult)
		p2Exp := math.Floor(loseXP * p2Mult)

		return int32(p1Exp), int32(p2Exp)
	} else {
		// Gana P2
		p1Exp := math.Floor(loseXP * p1Mult)
		p2Exp := math.Floor(winXP * p2Mult)

		return int32(p1Exp), int32(p2Exp)
	}
}

func getEloDifferenceMult(myElo int32, yourElo int32) float64 {
	return 1.0 + ((float64(myElo) - float64(yourElo)) / 1000.0)
}

func GetXPBarInfo(XP int32) (int32, int32) {
	level := GetLevel(XP)

	threshold := GetLevelXP(level+1) - GetLevelXP(level)

	currentLevelXP := XP - GetLevelXP(level)

	return currentLevelXP, threshold
}
