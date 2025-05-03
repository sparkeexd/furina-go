package entity

// Daily reward claim model from HoYoLAB API.
type DailyClaim struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
}
