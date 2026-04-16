package friendship

import (
	"context"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/sse"
)

type FriendshipService struct {
	store          *db.Store
	eventSystem    *sse.EventSystem
	accountService *account.AccountService
}

func NewService(s *db.Store, events *sse.EventSystem,
	accounts *account.AccountService) *FriendshipService {
	return &FriendshipService{store: s, eventSystem: events,
		accountService: accounts}
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
func (s FriendshipService) Create(ctx context.Context, data *FriendshipDTO) error {

	var AuxUser1ID int64
	var AuxUser2ID int64
	var AuxAccepted1 bool
	var AuxAccepted2 bool

	// Ordenamos para la inserción en BDD
	if *data.SenderID < data.ReceiverID {
		AuxUser1ID = *data.SenderID
		AuxUser2ID = data.ReceiverID
		AuxAccepted1 = true
		AuxAccepted2 = false
	} else if *data.SenderID > data.ReceiverID {
		AuxUser1ID = data.ReceiverID
		AuxUser2ID = *data.SenderID
		AuxAccepted1 = false
		AuxAccepted2 = true
	} else {
		//Bad_request Sender == receiver
		return apierror.ErrBadRequest
	}

	_, err := s.GetFriendshipStatus(ctx, data)
	if err == nil {
		return apierror.ErrAlreadyFriends
	}

	err = s.store.CreateFriendship(ctx, db.CreateFriendshipParams{
		User1ID:   AuxUser1ID,
		User2ID:   AuxUser2ID,
		Accepted1: AuxAccepted1,
		Accepted2: AuxAccepted2,
	})

	if err != nil {
		return err
	}

	senderUsername, err := s.accountService.GetUsernameByID(ctx, *data.SenderID)
	if err != nil {
		return err
	}
	s.eventSystem.SendEvent(data.ReceiverID, &sse.Event{
		EventType: "FriendRequest",
		Data:      senderUsername,
	}, true)

	return nil
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
	ctx context.Context,
	data *FriendshipDTO,
) (*FriendshipStatusDTO, error) {

	if *data.SenderID == data.ReceiverID {
		return nil, apierror.ErrBadRequest
	}

	res, err := s.store.GetFriendship(ctx, db.GetFriendshipParams{
		User1ID: *data.SenderID,
		User2ID: data.ReceiverID,
	})

	if err != nil {
		return nil, err
	}

	if *data.SenderID == res.User1ID {
		return &FriendshipStatusDTO{
			SenderID:       res.User1ID,
			ReceiverID:     res.User2ID,
			SenderAccept:   res.Accepted1,
			ReceiverAccept: res.Accepted2,
		}, err
	} else {
		return &FriendshipStatusDTO{
			SenderID:       res.User2ID,
			ReceiverID:     res.User1ID,
			SenderAccept:   res.Accepted2,
			ReceiverAccept: res.Accepted1,
		}, err
	}
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que devuelve un listado
de las amistades de un jugador dado su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* accountID, int64 - Es el id del usuario.
* --- Resultados ---
* []FriendshipReturnDTO - Lista con la información de las amistades del usuario,
incluyendo los datos relevante del otro usuario.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) GetUserFriendships(
	ctx context.Context,
	accountID int64,
) ([]FriendshipReturnDTO, error) {

	res, err := s.store.GetUserFriendships(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return ParseFriendshipRow(res), nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que devuelve un listado
de las amistades enviadas pendientes de un jugador dado su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* accountID, int64 - Es el id del usuario.
* --- Resultados ---
* []FriendshipReturnDTO - Lista con la información de las amistades enviadas
pendientes del usuario, incluyendo los datos relevante del otro usuario.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) GetUserPendingSentFriendships(
	ctx context.Context,
	accountID int64,
) ([]FriendshipReturnDTO, error) {

	res, err := s.store.GetUserPendingSentFriendships(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return ParsePendingSentRow(res), nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que devuelve un listado
de las amistades recibidas pendientes de un jugador dado su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* accountID, int64 - Es el id del usuario.
* --- Resultados ---
* []FriendshipReturnDTO - Lista con la información de las amistades recibidas
pendientes del usuario, incluyendo los datos relevante del otro usuario.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) GetUserPendingReceivedFriendships(
	ctx context.Context,
	accountID int64,
) ([]FriendshipReturnDTO, error) {

	res, err := s.store.GetUserPendingReceivedFriendships(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return ParsePendingReceivedRow(res), nil
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que devuelve un listado
de las amistades recibidas pendientes de un jugador dado su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* accountID, int64 - Es el id del usuario.
* --- Resultados ---
* []FriendshipReturnDTO - Lista con la información de las amistades recibidas
pendientes del usuario, incluyendo los datos relevante del otro usuario.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s FriendshipService) GetUserPendingReceivedFriendshipsCount(
	ctx context.Context,
	accountID int64,
) (RequestCount, error) {

	res, err := s.store.GetUserPendingReceivedFriendshipsCount(ctx, accountID)
	if err != nil {
		return RequestCount{Count: 0}, err
	}
	if len(res) == 0 {
		return RequestCount{Count: 0}, err
	}
	return RequestCount{Count: int32(res[0])}, nil

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
func (s FriendshipService) AcceptFriendship(
	ctx context.Context,
	data *FriendshipDTO,
) error {

	friendship, err := s.GetFriendshipStatus(ctx, data)
	if err != nil {
		return err
	}

	if friendship.SenderAccept && friendship.ReceiverAccept {
		return apierror.ErrAlreadyFriends
	}

	err = s.store.AcceptFriendship(ctx, db.AcceptFriendshipParams{
		User1ID: *data.SenderID,
		User2ID: data.ReceiverID,
	})
	if err != nil {
		return err
	}

	senderUsername, err := s.accountService.GetUsernameByID(ctx, *data.SenderID)
	if err != nil {
		return err
	}

	receiverUsername, err := s.accountService.GetUsernameByID(ctx, data.ReceiverID)
	if err != nil {
		return err
	}

	s.eventSystem.SendEvent(data.ReceiverID, &sse.Event{
		EventType: "NewFriend",
		Data:      senderUsername,
	}, true)

	s.eventSystem.SendEvent(*data.SenderID, &sse.Event{
		EventType: "NewFriend",
		Data:      receiverUsername,
	}, true)

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
func (s FriendshipService) DeleteFriendship(
	ctx context.Context,
	data *FriendshipDTO,
) error {

	err := s.store.DeleteFriendship(ctx, db.DeleteFriendshipParams{
		User1ID: *data.SenderID,
		User2ID: data.ReceiverID,
	})
	if err != nil {
		return err
	}

	return nil
}

// FUNCIONES AUXILIARES

// Funcion auxiliar: pasar de db.GetUserFriendshipsRow a FriendshipReturnDTO
func ParseFriendshipRow(data []db.GetUserFriendshipsRow) []FriendshipReturnDTO {
	var res []FriendshipReturnDTO

	for _, value := range data {
		res = append(res, FriendshipReturnDTO{
			UserID:   value.UserID,
			Username: value.Username,
			Level:    value.Level,
			Avatar:   value.Avatar,
		})
	}

	return res
}

// Funcion auxiliar: pasar de db.GetUserPendingSentFriendshipsRow
// a FriendshipReturnDTO
func ParsePendingSentRow(
	data []db.GetUserPendingSentFriendshipsRow,
) []FriendshipReturnDTO {
	var res []FriendshipReturnDTO

	for _, value := range data {
		res = append(res, FriendshipReturnDTO{
			UserID:   value.UserID,
			Username: value.Username,
			Level:    value.Level,
			Avatar:   value.Avatar,
		})
	}

	return res
}

// Funcion auxiliar: pasar de db.GetUserPendingReceivedFriendshipsRow
// a FriendshipReturnDTO
func ParsePendingReceivedRow(
	data []db.GetUserPendingReceivedFriendshipsRow,
) []FriendshipReturnDTO {

	var res []FriendshipReturnDTO

	for _, value := range data {
		res = append(res, FriendshipReturnDTO{
			UserID:   value.UserID,
			Username: value.Username,
			Level:    value.Level,
			Avatar:   value.Avatar,
		})
	}

	return res
}
