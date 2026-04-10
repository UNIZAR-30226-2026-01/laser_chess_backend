package rewards

import "math"

// La formula para los niveles es scale * nivel^2

const (
	scale = 250.0

	rankMult   = 1.2
	normalMult = 1

	winXP  = 150.0
	loseXP = 80.0

	winMoney  = 50
	loseMoney = 10
)

func GetLevel(XP int32) int32 {
	level := math.Sqrt(float64(XP) / scale)
	return int32(math.Floor(level))
}

func GetLevelXP(level int32) int32 {
	return int32(math.Floor(250.0 * float64(math.Pow(float64(level), 2))))
}

func GetXPBarInfo(XP int32) (int32, int32) {
	level := GetLevel(XP)

	threshold := GetLevelXP(level+1) - GetLevelXP(level)

	currentLevelXP := XP - GetLevelXP(level)

	return currentLevelXP, threshold
}

func getEloDifferenceMult(myElo int32, yourElo int32) float64 {
	return 1.0 + ((float64(myElo) - float64(yourElo)) / 1000.0)
}

func getMults(p1Elo int32, p2Elo int32, isRanked bool) float64 {
	eloMult := getEloDifferenceMult(p1Elo, p2Elo)

	var modeMult float64
	if isRanked {
		modeMult = rankMult
	} else {
		modeMult = normalMult
	}

	return eloMult * modeMult
}

func GetMatchXP(p1Elo int32, p2Elo int32, scoreP1 int, isRanked bool) (int32, int32) {
	p1Mult := getMults(p1Elo, p2Elo, isRanked)
	p2Mult := getMults(p2Elo, p2Elo, isRanked)

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

func GetMatchMoney(p1Elo int32, p2Elo int32, scoreP1 int, isRanked bool) (int32, int32) {
	p1Mult := getMults(p1Elo, p2Elo, isRanked)
	p2Mult := getMults(p2Elo, p2Elo, isRanked)

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
