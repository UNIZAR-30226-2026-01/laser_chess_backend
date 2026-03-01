package match

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type MatchService struct {
	store *db.Store
}

func NewService(s *db.Store) *MatchService {
	return &MatchService{store: s}
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que crea una partida
dado su DTO.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* data, MatchDTO - Es el DTO con los datos de la partida.
* --- Resultados ---
* MatchDTO - Objeto que contiene los ids de los jugadores.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s *MatchService) Create(ctx context.Context, data MatchDTO) (MatchDTO, error) {

	res, err := s.store.CreateMatch(ctx, db.CreateMatchParams{
		P1ID:            data.P1ID,
		P2ID:            data.P2ID,
		P1Elo:           data.P1Elo,
		P2Elo:           data.P2Elo,
		Date:            data.Date,
		Winner:          data.Winner,
		Termination:     data.Termination,
		MatchType:       data.MatchType,
		Board:           data.Board,
		MovementHistory: data.MovementHistory,
		TimeBase:        data.TimeBase,
		TimeIncrement:   data.TimeIncrement,
	})

	if err != nil {
		return MatchDTO{}, err
	}

	return MatchDTO{P1ID: res.P1ID, P2ID: res.P2ID}, nil
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
func (s *MatchService) GetByID(ctx context.Context, matchID int64) (MatchDTO, error) {

	res, err := s.store.GetMatch(ctx, matchID)
	if err != nil {
		return MatchDTO{}, err
	}

	return MatchDTO{
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
