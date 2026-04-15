package public

type MatchmakingRequestDTO struct {
	StartingTime  int32 `form:"time_base"`
	TimeIncrement int32 `form:"time_increment"`
	Board 		  int 	`form:"board"`
	Ranked        int   `form:"ranked"`
}
