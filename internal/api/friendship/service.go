package friendship

import (
	"context"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type FriendshipService struct {
	store *db.Store
}

func NewService(s *db.Store) *FriendshipService {
	return &FriendshipService{store: s}
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que crea una amistad
dado su DTO.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* data, FriendshipDTO - Es el DTO con los datos de la amistad.
* --- Resultados ---
* FriendshipDTO - Objeto que contiene los ids de los jugadores.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) Create(
	ctx context.Context, data FriendshipDTO) (FriendshipDTO, error) {

	var AuxUser1ID int64
	var AuxUser2ID int64
	var AuxAccepted1 bool
	var AuxAccepted2 bool

	// Ordenamos para la inserción en BDD
	if data.SenderID < data.RecieverID {
		AuxUser1ID = data.SenderID
		AuxUser2ID = data.RecieverID
		AuxAccepted1 = true
		AuxAccepted2 = false
	} else if data.SenderID > data.RecieverID {
		AuxUser1ID = data.RecieverID
		AuxUser2ID = data.SenderID
		AuxAccepted1 = false
		AuxAccepted2 = true
	} else {
		//Bad_request Sender == reciever
		return FriendshipDTO{}, apierror.ErrBadRequest
	}

	res, err := s.store.CreateFriendship(ctx, db.CreateFriendshipParams{
		User1ID:   AuxUser1ID,
		User2ID:   AuxUser2ID,
		Accepted1: AuxAccepted1,
		Accepted2: AuxAccepted2,
	})

	if err != nil {
		return FriendshipDTO{}, err
	}

	return FriendshipDTO{SenderID: res.User1ID, RecieverID: res.User2ID}, nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que devuelve los datos
de una amistad dados los ids de los usuarios.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* data, FriendshipDTO - Es el DTO con los ids de los usuarios.
* --- Resultados ---
* FriendshipDTO - Objeto que contiene la informacion de la amistad.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) GetFriendshipStatus(
	ctx context.Context, data FriendshipDTO) (FriendshipStatusDTO, error) {

	if data.SenderID == data.RecieverID {
		return FriendshipStatusDTO{}, apierror.ErrBadRequest
	}

	res, err := s.store.GetFriendship(ctx, db.GetFriendshipParams{
		User1ID: data.SenderID,
		User2ID: data.RecieverID,
	})

	if data.SenderID == res.User1ID {
		return FriendshipStatusDTO{
			SenderID:       res.User1ID,
			RecieverID:     res.User2ID,
			SenderAccept:   res.Accepted1,
			RecieverAccept: res.Accepted2,
		}, err
	} else {
		return FriendshipStatusDTO{
			SenderID:       res.User2ID,
			RecieverID:     res.User1ID,
			SenderAccept:   res.Accepted2,
			RecieverAccept: res.Accepted1,
		}, err
	}

}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que devuelve un listado
de las amistades de un jugador dado su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* userID, int64 - Es el id del usuario.
* --- Resultados ---
* []FriendshipReturnDTO - Lista con la información de las amistades del usuario,
incluyendo los datos relevante del otro usuario.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) GetUserFriendships(
	ctx context.Context, userID int64) ([]FriendshipReturnDTO, error) {

	res, err := s.store.GetUserFriendships(ctx, userID)
	if err != nil {
		return []FriendshipReturnDTO{}, err
	}

	return parseFrienshipRow(res), nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que devuelve un listado
de las amistades enviadas pendientes de un jugador dado su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* userID, int64 - Es el id del usuario.
* --- Resultados ---
* []FriendshipReturnDTO - Lista con la información de las amistades enviadas
pendientes del usuario, incluyendo los datos relevante del otro usuario.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) GetUserPendingSentFriendships(
	ctx context.Context, userID int64) ([]FriendshipReturnDTO, error) {

	res, err := s.store.GetUserPendingSentFriendships(ctx, userID)
	if err != nil {
		return []FriendshipReturnDTO{}, err
	}

	return parsePendingSentRow(res), nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que devuelve un listado
de las amistades recibidas pendientes de un jugador dado su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* userID, int64 - Es el id del usuario.
* --- Resultados ---
* []FriendshipReturnDTO - Lista con la información de las amistades recibidas
pendientes del usuario, incluyendo los datos relevante del otro usuario.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) GetUserPendingRecievedFriendships(
	ctx context.Context, userID int64) ([]FriendshipReturnDTO, error) {

	res, err := s.store.GetUserPendingRecievedFriendships(ctx, userID)
	if err != nil {
		return []FriendshipReturnDTO{}, err
	}

	return parsePendingRecievedRow(res), nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que marca una amistad
como aceptada dados los ids de ambos usuarios.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* data, FriendshipDTO - Es DTO que los ids de los usuarios.
* --- Resultados ---
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) AcceptFrienship(
	ctx context.Context, data FriendshipDTO) error {

	err := s.store.SetFriendship(ctx, db.SetFriendshipParams{
		User1ID: data.SenderID,
		User2ID: data.RecieverID,
	})
	if err != nil {
		return err
	}

	return nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que elimina una amistad
dados los ids de ambos usuarios.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* data, FriendshipDTO - Es DTO que los ids de los usuarios.
* --- Resultados ---
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) DeleteFrienship(
	ctx context.Context, data FriendshipDTO) error {

	err := s.store.DeleteFriendship(ctx, db.DeleteFriendshipParams{
		User1ID: data.SenderID,
		User2ID: data.RecieverID,
	})
	if err != nil {
		return err
	}

	return nil
}

// FUNCIONES AUXILIARES

// Funcion auxiliar: pasar de db.GetUserFrienshipsRow a FriendshipReturnDTO
func parseFrienshipRow(data []db.GetUserFriendshipsRow) []FriendshipReturnDTO {
	var res []FriendshipReturnDTO

	for _, value := range data {
		res = append(res, FriendshipReturnDTO{
			UserID:   value.UserID,
			Username: value.Username,
			Level:    value.Level,
			Xp:       value.Xp,
		})
	}

	return res
}

// Funcion auxiliar: pasar de db.GetUserPendingSentFriendshipsRow
// a FriendshipReturnDTO
func parsePendingSentRow(
	data []db.GetUserPendingSentFriendshipsRow) []FriendshipReturnDTO {
	var res []FriendshipReturnDTO

	for _, value := range data {
		res = append(res, FriendshipReturnDTO{
			UserID:   value.UserID,
			Username: value.Username,
			Level:    value.Level,
			Xp:       value.Xp,
		})
	}

	return res
}

// Funcion auxiliar: pasar de db.GetUserPendingRecievedFriendshipsRow
// a FriendshipReturnDTO
func parsePendingRecievedRow(
	data []db.GetUserPendingRecievedFriendshipsRow) []FriendshipReturnDTO {
	var res []FriendshipReturnDTO

	for _, value := range data {
		res = append(res, FriendshipReturnDTO{
			UserID:   value.UserID,
			Username: value.Username,
			Level:    value.Level,
			Xp:       value.Xp,
		})
	}

	return res
}
