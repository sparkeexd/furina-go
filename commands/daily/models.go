package daily

// Daily reward claim response structure from HoYoLab API.
type DailyRewardClaimResponse struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
}
