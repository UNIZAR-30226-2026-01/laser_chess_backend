package rating

import (
	"context"
	"errors"
	"strings"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type RatingService struct {
	store *db.Store
}

func NewService(s *db.Store) *RatingService {
	return &RatingService{store: s}
}

// Devuelve que tipo de elo corresponde a un tiempo
// Sirve tanto como para ranked, como para saber
// en una privada cual es el elo más cercano
func GetEloTypeFromBaseTime(baseTime int32) db.EloType {
	if baseTime < 600000 {
		// Blitz ~ 300 s
		return db.EloTypeBLITZ

	} else if baseTime < 1350000 {
		// Rapid ~ 900 s
		return db.EloTypeRAPID

	} else if baseTime < 2700000 {
		// Classic ~ 1800 s
		return db.EloTypeCLASSIC

	} else {
		// Extended ~ 3600 s
		return db.EloTypeEXTENDED
	}
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
func sqlcParamToDTO(res []db.Rating) *AllRatingsDTO {
	if len(res) == 0 {
		return &AllRatingsDTO{}
	}
	var blitz int32
	var rapid int32
	var extended int32
	var classic int32

	for _, r := range res {
		switch strings.ToUpper(string(r.EloType)) {
		case "BLITZ":
			blitz = r.Value
		case "RAPID":
			rapid = r.Value
		case "EXTENDED":
			extended = r.Value
		case "CLASSIC":
			classic = r.Value
		}
	}

	return &AllRatingsDTO{
		UserID:   res[0].UserID,
		Blitz:    blitz,
		Rapid:    rapid,
		Extended: extended,
		Classic:  classic,
	}
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
func (s *RatingService) GetAllElosByID(ctx context.Context, userID int64) (*AllRatingsDTO, error) {
	res, err := s.store.GetAllElos(ctx, userID)
	if err != nil {
		return nil, err
	}
	return sqlcParamToDTO(res), nil
}

func (s *RatingService) GetEloByID(ctx context.Context, userID int64, baseTime int32) (*RatingDTO, error) {
	eloType := GetEloTypeFromBaseTime(baseTime)

	switch eloType {
	case db.EloTypeBLITZ:
		return s.GetBlitzEloByID(ctx, userID)
	case db.EloTypeRAPID:
		return s.GetRapidEloByID(ctx, userID)
	case db.EloTypeCLASSIC:
		return s.GetClassicEloByID(ctx, userID)
	case db.EloTypeEXTENDED:
		return s.GetExtendedEloByID(ctx, userID)
	default:
		return nil, errors.New("tipo de elo inexistente")
	}
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
func (s *RatingService) GetBlitzEloByID(ctx context.Context, userID int64) (*RatingDTO, error) {
	res, err := s.store.GetBlitzElo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &RatingDTO{
		UserID:        res.UserID,
		EloType:       res.EloType,
		Value:         res.Value,
		Deviation:     res.Deviation,
		Volatility:    res.Volatility,
		LastUpdatedAt: res.LastUpdatedAt,
	}, nil
}

func (s *RatingService) GetExtendedEloByID(ctx context.Context, userID int64) (*RatingDTO, error) {
	res, err := s.store.GetExtendedElo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &RatingDTO{
		UserID:        res.UserID,
		EloType:       res.EloType,
		Value:         res.Value,
		Deviation:     res.Deviation,
		Volatility:    res.Volatility,
		LastUpdatedAt: res.LastUpdatedAt,
	}, nil
}

func (s *RatingService) GetRapidEloByID(ctx context.Context, userID int64) (*RatingDTO, error) {
	res, err := s.store.GetRapidElo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &RatingDTO{
		UserID:        res.UserID,
		EloType:       res.EloType,
		Value:         res.Value,
		Deviation:     res.Deviation,
		Volatility:    res.Volatility,
		LastUpdatedAt: res.LastUpdatedAt,
	}, nil
}

func (s *RatingService) GetClassicEloByID(ctx context.Context, userID int64) (*RatingDTO, error) {
	res, err := s.store.GetClassicElo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &RatingDTO{
		UserID:        res.UserID,
		EloType:       res.EloType,
		Value:         res.Value,
		Deviation:     res.Deviation,
		Volatility:    res.Volatility,
		LastUpdatedAt: res.LastUpdatedAt,
	}, nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que
* actualiza el valor del rating de un jugador en una categoría
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* rating, RatingDTO - Contiene los datos del jugador, de la categoría,
* y los nuevos valores que tiene el rating del jugador en esa categoría.
* --- Resultados ---
* error - Es el error que se haya provocado en la consulta, o nil en caso
* contrario.
*
 */
func (s *RatingService) UpdateEloByID(ctx context.Context, rating *RatingDTO) error {
	err := s.store.UpdateRating(ctx, db.UpdateRatingParams{
		Value:      rating.Value,
		Deviation:  rating.Deviation,
		Volatility: rating.Volatility,
		UserID:     rating.UserID,
		EloType:    rating.EloType,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *RatingService) GetTopRankUsers(ctx context.Context,
	eloType string) ([]RankUserDTO, error) {
	eloType = strings.ToUpper(eloType)
	res, err := s.store.GetTopRankUsers(ctx, db.EloType(eloType))
	if err != nil {
		return nil, err
	}
	return ParseRankingRow(res), nil
}

func (s *RatingService) GetRankById(ctx context.Context,
	eloType string, userID int64) (int64, error) {
	eloType = strings.ToUpper(eloType)
	res, err := s.store.GetRankById(ctx, db.GetRankByIdParams{
		EloType: db.EloType(eloType),
		UserID:  userID,
	})
	if err != nil {
		return 0, err
	}
	return res, nil
}

func ParseRankingRow(
	data []db.GetTopRankUsersRow,
) []RankUserDTO {

	var res []RankUserDTO

	for _, value := range data {
		res = append(res, RankUserDTO{
			UserID:   value.AccountID,
			Username: value.Username,
			Avatar:   value.Avatar,
			Rating:   value.Value,
		})
	}

	return res
}
