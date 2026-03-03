package friendship

// DTOs para tratar con frienships

// Para crear una petición de amistad
type FriendshipDTO struct {
	SenderID   *int64 `json:"user1_id"`
	ReceiverID int64  `json:"user2_id" binding:"required"`
}

// Para recibir una petición de amistad
type FriendshipStatusDTO struct {
	SenderID       int64 `json:"user1_id" binding:"required"`
	ReceiverID     int64 `json:"user2_id" binding:"required"`
	SenderAccept   bool
	ReceiverAccept bool
}

// Para devolver los datos de los usuarios amigos de una consulta
type FriendshipReturnDTO struct {
	UserID   int64  `json:"account_id" binding:"required"`
	Username string `json:"username"`
	Level    int32  `json:"level"`
	Xp       int32  `json:"xp"`
}
