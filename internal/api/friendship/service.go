package friendship

import (
	"context"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type friendshipService struct {
	store *db.Store
}

func NewService(s *db.Store) *friendshipService {
	return &friendshipService{store: s}
}

/** CREATE
 * 	Recibe: dos usuarios, «senderID» y «recieverID», los ordena para que no haya 
 *  repetidos en base de datos, pone la petición aceptada del lado del remitente
 * 	e inserta en orden de tamaño User1ID - el menor y User2ID - el mayor
 *  Devuelve: error si no se puede crear la amistad porque ya existe o por otro...
 */
func (s *friendshipService) Create(
	ctx context.Context, data FriendshipDTO) (error) {
	
	var AuxUser1ID int64
	var AuxUser2ID int64
	var AuxAccepted1 bool
	var AuxAccepted2 bool

	// Ordenamos para la inserción en BDD
	if 			(data.SenderID < data.RecieverID) {
		AuxUser1ID = data.SenderID
		AuxUser2ID = data.RecieverID
		AuxAccepted1 = true
		AuxAccepted2 = false
	} else if	(data.SenderID > data.RecieverID){
		AuxUser1ID = data.RecieverID
		AuxUser2ID = data.SenderID
		AuxAccepted1 = false
		AuxAccepted2 = true
	} else {
		//Bad_request Sender == reciever
		return apierror.ErrBadRequest
	}

	// NOTA - Puede estar bien devolver algo significativo cuando ya existe en
	// base de datos
	_ , err := s.store.CreateFriendship(ctx, db.CreateFriendshipParams{
		User1ID: AuxUser1ID,
		User2ID: AuxUser2ID,
		Accepted1: AuxAccepted1,
		Accepted2: AuxAccepted2,
	})

	if err != nil {
		return err
	}

	return err
}

/** GET FRIENDSHIP
 *	Recibe: dos usuarios, «senderID» y «recieverID», donde sender es el usuario
 *  principal, y revieverID es el usuario secundario y devuelve el estado de la
 *	amistad
 */
func (s *friendshipService) GetFreindshipStatus(
	ctx context.Context, data FriendshipDTO) (FriendshipStatusDTO, error){

	var AuxUser1ID int64
	var AuxUser2ID int64
	var AuxAccepted1 bool
	var AuxAccepted2 bool

	// Ordenamos para la inserción en BDD
	if 			(data.SenderID < data.RecieverID) {
		AuxUser1ID = data.SenderID
		AuxUser2ID = data.RecieverID
	} else if	(data.SenderID > data.RecieverID){
		AuxUser1ID = data.RecieverID
		AuxUser2ID = data.SenderID
	} else {
		//Bad_request Sender == reciever
		return FriendshipStatusDTO{}, apierror.ErrBadRequest
	}

	res, err := s.store.GetFriendship(ctx, db.GetFriendshipParams{
		User1ID: AuxUser1ID,
		User2ID: AuxUser2ID,
	})

	if err != nil {
		return FriendshipStatusDTO{}, err
	}
		if 			(data.SenderID < data.RecieverID) {
		AuxAccepted1 = res.Accepted1
		AuxAccepted2 = res.Accepted2
	} else {//		(data.SenderID > data.RecieverID) {
		AuxAccepted1 = res.Accepted2
		AuxAccepted2 = res.Accepted2
	}

	return FriendshipStatusDTO{
			SenderID: 		AuxUser1ID,
			RecieverID: 	AuxUser2ID,
			SenderAccept: 	AuxAccepted1,
			RecieverAccept: AuxAccepted2,
		}, err
}