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

/** CREATE
 * 	Recibe: dos usuarios, «senderID» y «recieverID», los ordena para que no haya
 *  repetidos en base de datos, pone la petición aceptada del lado del remitente
 * 	e inserta en orden de tamaño User1ID - el menor y User2ID - el mayor
 *  Devuelve: error si no se puede crear la amistad porque ya existe o por otro...
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

	// NOTA - Puede estar bien devolver algo significativo cuando ya existe en
	// base de datos
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

/** GET FRIENDSHIP
 *	Recibe: dos usuarios, «senderID» y «recieverID», donde sender es el usuario
 *  principal, y revieverID es el usuario secundario y devuelve el estado de la
 *	amistad
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

	// var AuxUser1ID int64
	// var AuxUser2ID int64
	// var AuxAccepted1 bool
	// var AuxAccepted2 bool

	// // Ordenamos para la inserción en BDD
	// if data.SenderID < data.RecieverID {
	// 	AuxUser1ID = data.SenderID
	// 	AuxUser2ID = data.RecieverID
	// } else if data.SenderID > data.RecieverID {
	// 	AuxUser1ID = data.RecieverID
	// 	AuxUser2ID = data.SenderID
	// } else {
	// 	//Bad_request Sender == reciever
	// 	return FriendshipStatusDTO{}, apierror.ErrBadRequest
	// }

	// if err != nil {
	// 	return FriendshipStatusDTO{}, err
	// }
	// if data.SenderID < data.RecieverID {
	// 	AuxAccepted1 = res.Accepted1
	// 	AuxAccepted2 = res.Accepted2
	// } else { //		(data.SenderID > data.RecieverID) {
	// 	AuxAccepted1 = res.Accepted2
	// 	AuxAccepted2 = res.Accepted2
	// }

	// return FriendshipStatusDTO{
	// 	SenderID:       AuxUser1ID,
	// 	RecieverID:     AuxUser2ID,
	// 	SenderAccept:   AuxAccepted1,
	// 	RecieverAccept: AuxAccepted2,
	// }, err
}

func (s FriendshipService) GetUserFriendships(
	ctx context.Context, userID int64) ([]FriendshipReturnDTO, error) {

	res, err := s.store.GetUserFriendships(ctx, userID)
	if err != nil {
		return []FriendshipReturnDTO{}, err
	}

	return parseFrienshipRow(res), nil
}

func (s FriendshipService) GetUserPendingSentFriendships(
	ctx context.Context, userID int64) ([]FriendshipReturnDTO, error) {

	res, err := s.store.GetUserPendingSentFriendships(ctx, userID)
	if err != nil {
		return []FriendshipReturnDTO{}, err
	}

	return parsePendingSentRow(res), nil
}

func (s FriendshipService) GetUserPendingRecievedFriendships(
	ctx context.Context, userID int64) ([]FriendshipReturnDTO, error) {

	res, err := s.store.GetUserPendingRecievedFriendships(ctx, userID)
	if err != nil {
		return []FriendshipReturnDTO{}, err
	}

	return parsePendingRecievedRow(res), nil
}

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
