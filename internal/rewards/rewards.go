package rewards

import "math"

// La formula para los niveles es scale * nivel^2

const (
	scale = 250.0

	rankMult   = 1.2
	normalMult = 1.0

	winXP  = 150.0
	loseXP = 80.0

	winMoney  = 50.0
	loseMoney = 10.0
)

// Devuelve el nivel dado xp global
func GetLevel(XP int32) int32 {
	level := math.Sqrt(float64(XP) / scale)
	return int32(math.Floor(level))
}

// Devuelve la XP necesaria para alcanzar un nivel especifico
func GetLevelXP(level int32) int32 {
	return int32(math.Floor(scale * float64(level) * float64(level)))
}

// Devuelve la XP de dentro del nivel, y la xp maxima del nivel
// Es decir si a nivel 10 llegas con xp 100, al 11 con 110, y tienes 105 xp
// 	devuelve 5, 10
func GetXPBarInfo(XP int32) (int32, int32) {
	level := GetLevel(XP)

	threshold := GetLevelXP(level+1) - GetLevelXP(level)
	currentLevelXP := XP - GetLevelXP(level)

	return currentLevelXP, threshold
}

func getEloDifferenceMult(myElo int32, yourElo int32) float64 {
	rawDiff := (float64(yourElo) - float64(myElo)) / 1000.0
	mult := 1.0 + rawDiff

	return math.Max(0.5, math.Min(1.5, mult))
}

func getMults(myElo int32, enemyElo int32, isRanked bool) float64 {
	eloMult := getEloDifferenceMult(myElo, enemyElo)

	var modeMult float64
	if isRanked {
		modeMult = rankMult
	} else {
		modeMult = normalMult
	}

	return eloMult * modeMult
}

func GetMatchXP(p1Elo int32, p2Elo int32, scoreP1 float64, isRanked bool) (int32, int32) {
	p1Mult := getMults(p1Elo, p2Elo, isRanked)
	p2Mult := getMults(p2Elo, p1Elo, isRanked)

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

func GetMatchMoney(p1Elo int32, p2Elo int32, scoreP1 float64, isRanked bool) (int32, int32) {
	p1Mult := getMults(p1Elo, p2Elo, isRanked)
	p2Mult := getMults(p2Elo, p1Elo, isRanked)

	if scoreP1 >= 1 {
		// Gana P1
		p1Money := math.Floor(winMoney * p1Mult)
		p2Money := math.Floor(loseMoney * p2Mult)

		return int32(p1Money), int32(p2Money)
	} else {
		// Gana P2
		p1Money := math.Floor(loseMoney * p1Mult)
		p2Money := math.Floor(winMoney * p2Mult)

		return int32(p1Money), int32(p2Money)
	}
}
