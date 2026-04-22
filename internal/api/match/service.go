package match

import (
	"context"
	"fmt"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/rating"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/boards"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/elo"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/rewards"
)

type MatchService struct {
	store          *db.Store
	accountService *account.AccountService
	ratingService  *rating.RatingService
}

func NewService(s *db.Store, as *account.AccountService, rs *rating.RatingService) *MatchService {
	return &MatchService{
		store:          s,
		accountService: as,
		ratingService:  rs,
	}
}

func toCreateMatchParamsFromSaveDTO(data MatchSaveDTO) db.CreateMatchParams {
	return db.CreateMatchParams{
		P1ID:            data.P1ID,
		P2ID:            data.P2ID,
		P1Elo:           data.P1Elo,
		P2Elo:           data.P2Elo,
		Date:            data.Date,
		Winner:          db.Winner(data.GameInfo.Winner),
		Termination:     db.Termination(data.GameInfo.Termination),
		MatchType:       db.MatchType(data.GameInfo.MatchType),
		Board:           boards.IntToBoard[data.GameInfo.BoardType],
		MovementHistory: data.GameInfo.Log,
		TimeBase:        int32(data.GameInfo.TimeBase),
		TimeIncrement:   int32(data.GameInfo.TimeIncrement),
	}
}

func toUpdateMatchParamsFromSaveDTO(data MatchSaveDTO) db.UpdateMatchParams {
	return db.UpdateMatchParams{
		MatchID:         data.GameInfo.MatchID,
		P1Elo:           data.P1Elo,
		P2Elo:           data.P2Elo,
		Date:            data.Date,
		Winner:          db.Winner(data.GameInfo.Winner),
		Termination:     db.Termination(data.GameInfo.Termination),
		MatchType:       db.MatchType(data.GameInfo.MatchType),
		Board:           db.BoardType(data.GameInfo.BoardType),
		MovementHistory: data.GameInfo.Log,
		TimeBase:        int32(data.GameInfo.TimeBase),
		TimeIncrement:   int32(data.GameInfo.TimeIncrement),
	}
}

// Funcion auxiliar: pasar de db.Match a MatchDTO
func parseMatches(data []db.Match) []MatchDTO {
	var res []MatchDTO

	for _, value := range data {
		res = append(res, MatchDTO{
			P1ID:            value.P1ID,
			P2ID:            value.P2ID,
			P1Elo:           value.P1Elo,
			P2Elo:           value.P2Elo,
			Date:            value.Date,
			Winner:          value.Winner,
			Termination:     value.Termination,
			MatchType:       value.MatchType,
			Board:           value.Board,
			MovementHistory: value.MovementHistory,
			TimeBase:        value.TimeBase,
			TimeIncrement:   value.TimeIncrement,
		})
	}

	return res
}

// Funcion auxiliar: pasar de db.Match a PausedMatchDTO
func parsePausedMatches(data []db.Match) []PausedMatchDTO {
	var res []PausedMatchDTO

	for _, value := range data {
		res = append(res, PausedMatchDTO{
			MatchID:         value.MatchID,
			P1ID:            value.P1ID,
			P2ID:            value.P2ID,
			P1Elo:           value.P1Elo,
			P2Elo:           value.P2Elo,
			Date:            value.Date,
			Winner:          value.Winner,
			Termination:     value.Termination,
			MatchType:       value.MatchType,
			Board:           value.Board,
			MovementHistory: value.MovementHistory,
			TimeBase:        value.TimeBase,
			TimeIncrement:   value.TimeIncrement,
		})
	}

	return res
}

func (s *MatchService) FinalizeMatch(ctx context.Context, summary MatchSummaryDTO) error {
	isRanked := summary.GameInfo.MatchType == "RANKED"

	// Obtener elos actuales de los players
	p1RatingData, err := s.ratingService.GetEloByID(ctx, summary.P1ID, summary.GameInfo.TimeBase)
	if err != nil {
		return fmt.Errorf("error obteniendo elo p1: %w", err)
	}

	p2RatingData, err := s.ratingService.GetEloByID(ctx, summary.P2ID, summary.GameInfo.TimeBase)
	if err != nil {
		return fmt.Errorf("error obteniendo elo p2: %w", err)
	}

	// Obtener xp, level y money actuales de los players
	p1StatsData, err := s.accountService.GetStats(ctx, summary.P1ID)
	if err != nil {
		return fmt.Errorf("error obteniendo stats p1: %w", err)
	}

	p2StatsData, err := s.accountService.GetStats(ctx, summary.P2ID)
	if err != nil {
		return fmt.Errorf("error obteniendo stats p2: %w", err)
	}

	p1Elo := elo.Rating{
		Value:      float64(p1RatingData.Value),
		Deviation:  float64(p1RatingData.Deviation),
		Volatility: p1RatingData.Volatility,
	}
	p2Elo := elo.Rating{
		Value:      float64(p2RatingData.Value),
		Deviation:  float64(p2RatingData.Deviation),
		Volatility: p2RatingData.Volatility,
	}

	var newP1Rating, newP2Rating elo.Rating = p1Elo, p2Elo
	var p1GainedXP, p2GainedXP, p1GainedMoney, p2GainedMoney int32 = 0, 0, 0, 0

	// Si hay resultado, calcular nuevos elos y recompensas
	if summary.GameInfo.Winner != "NONE" {
		var scoreP1 float64
		if summary.GameInfo.Winner == "P1_WINS" {
			scoreP1 = 1.0
		} else {
			scoreP1 = 0.0
		}

		if isRanked {
			// Aplicar inactividad
			p1Elo = elo.ApplyInactivity(p1Elo, p1RatingData.LastUpdatedAt)
			p2Elo = elo.ApplyInactivity(p2Elo, p2RatingData.LastUpdatedAt)

			// Procesar Glicko2
			newP1Rating, newP2Rating = elo.ProcessMatch(p1Elo, p2Elo, scoreP1)
		}

		// Calcular rewards
		p1GainedXP, p2GainedXP = rewards.GetMatchXP(p1RatingData.Value, p2RatingData.Value, scoreP1, isRanked)
		p1GainedMoney, p2GainedMoney = rewards.GetMatchMoney(p1RatingData.Value, p2RatingData.Value, scoreP1, isRanked)
	}

	// Preparar datos para query
	matchData := MatchSaveDTO{
		IsNewMatch: summary.IsNewMatch,
		GameInfo:   summary.GameInfo,
		P1ID:       summary.P1ID,
		P2ID:       summary.P2ID,
		P1Elo:      int32(newP1Rating.Value),
		P2Elo:      int32(newP2Rating.Value),
		Date:       summary.Date,
	}

	newP1XP := p1StatsData.Xp + p1GainedXP
	newP2XP := p2StatsData.Xp + p2GainedXP

	newP1Money := p1StatsData.Money + p1GainedMoney
	newP2Money := p2StatsData.Money + p2GainedMoney

	// Guardar y actualizar elos y recompensas
	return s.SaveMatchResultTx(ctx, matchData, newP1Rating, newP2Rating,
		newP1XP, newP2XP, newP1Money, newP2Money)

}

/*
* Guarda una partida y actualiza los Elos de ambos jugadores de
* forma atómica.
 */
func (s *MatchService) SaveMatchResultTx(ctx context.Context,
	match MatchSaveDTO, p1Rating elo.Rating, p2Rating elo.Rating,
	p1XP int32, p2XP int32, p1Money int32, p2Money int32) error {

	return s.store.ExecTx(ctx, func(q *db.Queries) error {

		// Guardar partida
		if match.IsNewMatch {
			dbParams := toCreateMatchParamsFromSaveDTO(match)
			_, err := q.CreateMatch(ctx, dbParams)
			if err != nil {
				return err
			}
		} else {
			dbParams := toUpdateMatchParamsFromSaveDTO(match)
			_, err := q.UpdateMatch(ctx, dbParams)
			if err != nil {
				return err
			}
		}

		// Actualizar elos
		if match.GameInfo.MatchType == "RANKED" {
			eloType := rating.GetEloTypeFromBaseTime(match.GameInfo.TimeBase)

			err := q.UpdateRating(ctx, db.UpdateRatingParams{
				UserID:     match.P1ID,
				EloType:    eloType,
				Value:      int32(p1Rating.Value),
				Deviation:  int32(p1Rating.Deviation),
				Volatility: p1Rating.Volatility,
			})
			if err != nil {
				return err
			}

			err = q.UpdateRating(ctx, db.UpdateRatingParams{
				UserID:     match.P2ID,
				EloType:    eloType,
				Value:      int32(p2Rating.Value),
				Deviation:  int32(p2Rating.Deviation),
				Volatility: p2Rating.Volatility,
			})
			if err != nil {
				return err
			}
		}

		// Actualizar XP y Dinero
		err := q.UpdateStats(ctx, db.UpdateStatsParams{
			AccountID: match.P1ID,
			Level:     rewards.GetLevel(p1XP),
			Xp:        p1XP,
			Money:     p1Money,
		})
		if err != nil {
			return err
		}

		err = q.UpdateStats(ctx, db.UpdateStatsParams{
			AccountID: match.P2ID,
			Level:     rewards.GetLevel(p2XP),
			Xp:        p2XP,
			Money:     p2Money,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que busca una partida
dado su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* matchID, int64 - Es el id de la partida.
* --- Resultados ---
* MatchDTO - Objeto que contiene los datos de la partida.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s *MatchService) GetByID(ctx context.Context, matchID int64) (*MatchDTO, error) {

	res, err := s.store.GetMatch(ctx, matchID)
	if err != nil {
		return nil, err
	}

	return &MatchDTO{
		P1ID:            res.P1ID,
		P2ID:            res.P2ID,
		P1Elo:           res.P1Elo,
		P2Elo:           res.P2Elo,
		Date:            res.Date,
		Winner:          res.Winner,
		Termination:     res.Termination,
		MatchType:       res.MatchType,
		Board:           res.Board,
		MovementHistory: res.MovementHistory,
		TimeBase:        res.TimeBase,
		TimeIncrement:   res.TimeIncrement,
	}, nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que busca el historial
de partidas de un jugador dado su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* userID, int64 - Es el id de la partida.
* --- Resultados ---
* []MatchDTO - Listado de las partidas del jugador.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s *MatchService) GetUserHistory(ctx context.Context, userID int64) ([]MatchDTO, error) {
	res, err := s.store.GetUserHistory(ctx, userID)
	println(len(res))
	if err != nil {
		return nil, err
	}

	return parseMatches(res), nil
}

func (s *MatchService) GetPausedMatches(ctx context.Context, userID int64) ([]PausedMatchDTO, error) {
	res, err := s.store.GetPausedMatches(ctx, userID)
	println(len(res))
	if err != nil {
		return nil, err
	}

	return parsePausedMatches(res), nil
}
