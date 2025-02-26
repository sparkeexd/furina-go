package daily

// Daily reward endpoints are shared across different games with only minor differences to the URL.
// This struct consolidates the common differences between each game.
type DailyReward struct {
	BaseUrl  string
	EventId  string
	ActId    string
	SignGame string
}

// Daily reward claim response structure from HoYoLab API.
type DailyRewardClaimResponse struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
}
