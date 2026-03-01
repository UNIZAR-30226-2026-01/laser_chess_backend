package rating

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

const INITIAL_RATING = 1500

type RatingService struct {
	store *db.Store
}

func NewService(s *db.Store) *RatingService {
	return &RatingService{store: s}
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que
* crea todos los ratings de un jugador
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* userID, int64 - Es el id del jugador del que se van a crear los ratings.
* --- Resultados ---
* AllRatingsDTO - Contiene el id del jugador y los nuevos valores
* de rating que se le han asignado.
* error - Es el error que se haya provocado en la consulta, o nil en caso
* contrario.
*
 */
func (s *RatingService) CreateRating(ctx context.Context, userID int64) (AllRatingsDTO, error) {
	newRating := db.CreateRatingsParams{
		UserID:  userID,
		Value:   INITIAL_RATING,
		Value_2: INITIAL_RATING,
		Value_3: INITIAL_RATING,
		Value_4: INITIAL_RATING,
	}

	res, err := s.store.CreateRatings(ctx, newRating)
	if err != nil {
		return AllRatingsDTO{}, err
	}
	return s.SqlcParamToDTO(res), err
}

/*
*
* Desc: Función intermedia que permite cambiar del restultado que devuelve
* sqlc al DTO que contiene todos los ratings de manera
* --- Parametros ---
* res, []db.Rating - Es un array de los Ratings que tiene un jugador.
* --- Resultados ---
* AllRatingsDTO - Contiene el id del jugador y los nuevos valores
* de rating que se le han asignado.
*
 */
func (s RatingService) SqlcParamToDTO(res []db.Rating) AllRatingsDTO {
	ratingsDTO := AllRatingsDTO{}
	ratingsDTO.UserID = res[0].UserID
	ratingsDTO.Elo1.EloType = res[0].EloType
	ratingsDTO.Elo1.Value = res[0].Value
	ratingsDTO.Elo2.EloType = res[1].EloType
	ratingsDTO.Elo2.Value = res[1].Value
	ratingsDTO.Elo3.EloType = res[2].EloType
	ratingsDTO.Elo3.Value = res[2].Value
	ratingsDTO.Elo4.EloType = res[3].EloType
	ratingsDTO.Elo4.Value = res[3].Value
	return ratingsDTO
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que
* obtiene todos los ratings de un jugador determinado
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* userID, int64 - Es el id del jugador del que se van a crear los ratings.
* --- Resultados ---
* AllRatingsDTO - Contiene el id del jugador y los nuevos valores
* de rating que se le han asignado.
* error - Es el error que se haya provocado en la consulta, o nil en caso
* contrario.
*
 */
func (s *RatingService) GetAllElosByID(ctx context.Context, userID int64) (AllRatingsDTO, error) {
	res, err := s.store.GetAllElos(ctx, userID)
	if err != nil {
		return AllRatingsDTO{}, err
	}
	return s.SqlcParamToDTO(res), nil
}

/*
*
* Desc: Esta funcion y todas las posteriores con nombre similar,
* llama a una query generada por sqlc que devuelve el valor del elo de un jugador
* en el modo de juego contenido en el nombre de la función
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* userID, int64 - Es el id del jugador del que se van a crear los ratings.
* --- Resultados ---
* RatingDTO - Contiene el id del jugador, la categoría del rating,
* y el valor que tiene el jugador en esa categoría
* error - Es el error que se haya provocado en la consulta, o nil en caso
* contrario.
*
 */
func (s *RatingService) GetBlitzEloByID(ctx context.Context, userID int64) (RatingDTO, error) {
	res, err := s.store.GetBlitzElo(ctx, userID)
	if err != nil {
		return RatingDTO{}, err
	}
	return RatingDTO{
		UserID:  res.UserID,
		EloType: res.EloType,
		Value:   res.Value,
	}, nil
}

func (s *RatingService) GetBulletEloByID(ctx context.Context, userID int64) (RatingDTO, error) {
	res, err := s.store.GetBulletElo(ctx, userID)
	if err != nil {
		return RatingDTO{}, err
	}
	return RatingDTO{
		UserID:  res.UserID,
		EloType: res.EloType,
		Value:   res.Value,
	}, nil
}

func (s *RatingService) GetRapidEloByID(ctx context.Context, userID int64) (RatingDTO, error) {
	res, err := s.store.GetRapidElo(ctx, userID)
	if err != nil {
		return RatingDTO{}, err
	}
	return RatingDTO{
		UserID:  res.UserID,
		EloType: res.EloType,
		Value:   res.Value,
	}, nil
}

func (s *RatingService) GetClassicEloByID(ctx context.Context, userID int64) (RatingDTO, error) {
	res, err := s.store.GetClassicElo(ctx, userID)
	if err != nil {
		return RatingDTO{}, err
	}
	return RatingDTO{
		UserID:  res.UserID,
		EloType: res.EloType,
		Value:   res.Value,
	}, nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que
* actualiza el valor del rating de un jugador en una categoría
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* rating, RatingDTO - Contiene los datos del jugador, de la categoría,
* y del nuevo valor que tiene el rating del jugador en esa categoría.
* --- Resultados ---
* RatingDTO - Contiene el id del jugador, la categoría del rating,
* y el valor que tiene el jugador en esa categoría
* error - Es el error que se haya provocado en la consulta, o nil en caso
* contrario.
*
 */
func (s *RatingService) UpdateEloByID(ctx context.Context, rating RatingDTO) (RatingDTO, error) {
	res, err := s.store.UpdateRating(ctx, db.UpdateRatingParams{
		UserID:  rating.UserID,
		EloType: rating.EloType,
		Value:   rating.Value,
	})
	if err != nil {
		return RatingDTO{}, err
	}
	return RatingDTO{
		UserID:  res.UserID,
		EloType: res.EloType,
		Value:   res.Value,
	}, nil
}
