package match

import (
	"context"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/rating"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/elo"
)

type MatchService struct {
	store *db.Store
}

func NewService(s *db.Store) *MatchService {
	return &MatchService{store: s}
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
		Board:           db.BoardType(data.GameInfo.BoardType),
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

/*
* Guarda una partida sin modificar Elos. Crea o actualiza dependiendo
* de si es una partida nueva o retomada.
 */
func (s *MatchService) SaveMatch(ctx context.Context, data MatchSaveDTO) error {
	if data.IsNewMatch {
		dbParams := toCreateMatchParamsFromSaveDTO(data)
		_, err := s.store.CreateMatch(ctx, dbParams)
		return err
	}

	dbParams := toUpdateMatchParamsFromSaveDTO(data)
	_, err := s.store.UpdateMatch(ctx, dbParams)
	return err
}

/*
* Guarda una partida y actualiza los Elos de ambos jugadores de
* forma atómica.
 */
func (s *MatchService) SaveMatchResultTx(ctx context.Context,
	match MatchSaveDTO, p1Rating elo.Rating, p2Rating elo.Rating) error {

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
		eloType, eloErr := rating.GetEloTypeFromBaseTime(match.GameInfo.TimeBase)
		if eloErr != nil {
			return eloErr
		}

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

		// Actualizar elos
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
