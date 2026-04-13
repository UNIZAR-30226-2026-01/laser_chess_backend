package bot

type BotMatchRequestDTO struct {
	Board         *int   `form:"board"`
	StartingTime  *int32 `form:"starting_time"`
	TimeIncrement *int32 `form:"time_increment"`
	Level         *int   `form:"level"`
}
