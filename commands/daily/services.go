package daily

import (
	"fmt"
	"net/http"

	"github.com/sparkeexd/mimo/internal/network"
)

// Daily reward endpoints are shared across different games with only minor differences to the URL.
// This struct consolidates the common differences between each game.
type DailyReward struct {
	BaseURL  string
	EventID  string
	ActID    string
	SignGame string
}

// Constructor.
func NewDailyReward(baseUrl string, eventId string, actId string, signGame string) DailyReward {
	return DailyReward{baseUrl, eventId, actId, signGame}
}

// Claim daily reward.
// e.g. Genshin daily sign in endpoint: https://sg-hk4e-api.hoyolab.com/event/sol/sign?act_id=e202102251931481
func (daily DailyReward) Claim(cookie network.Cookie) (DailyRewardClaimResponse, error) {
	var res DailyRewardClaimResponse

	handler := network.NewHandler()
	endpoint := fmt.Sprintf("%s/event/%s/%s?act_id=%s", daily.BaseURL, daily.EventID, DailyRewardSignParam, daily.ActID)

	request := network.NewRequest(endpoint, http.MethodPost, cookie).
		AddParam("lang", LangEnglish).
		AddHeader("x-rpc-signgame", daily.SignGame).
		Build()

	err := handler.Send(request, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
