package hoyolab

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sparkeexd/furina/internal/domain/entity"
	"github.com/sparkeexd/furina/pkg/logger"
	"github.com/sparkeexd/furina/pkg/network"
)

const (
	// Main API endpoints.
	Hk4eEndpoint     = "https://sg-hk4e-api.hoyolab.com"
	SgPublicEndpoint = "https://sg-public-api.hoyolab.com"

	// Daily reward endpoint parameters.
	DailyRewardHomeParam = "home"
	DailyRewardInfoParam = "info"
	DailyRewardSignParam = "sign"

	// Genshin endpoint parameters.
	GenshinEventID = "sol"
	GenshinActID   = "e202102251931481"

	// Star Rail endpoint parameters.
	StarRailEventID = "luna/os"
	StarRailActID   = "e202303301540311"

	// Zenless endpoint parameters.
	ZenlessEventID = "luna/zzz/os"
	ZenlessActID   = "e202406031448091"

	// Header values for x-rpc-signgame.
	GenshinSignGame  = "genshin"
	StarRailSignGame = "hsr"
	ZenlessSignGame  = "zzz"

	// Language.
	LangEnglish = "en-us"
)

// Repository for handling daily reward claim.
type DailyRepository struct {
	Logger *logger.Logger
}

// Daily reward endpoints are shared across different games with only minor differences to the URL.
// This context consolidates the common differences between each game.
type DailyRewardContext struct {
	BaseURL  string
	EventID  string
	ActID    string
	SignGame string
}

// Create a new daily repository.
func NewDailyRepository(logger *logger.Logger) DailyRepository {
	return DailyRepository{Logger: logger}
}

// Create a new daily reward context.
func NewDailyRewardContext(baseURL, eventID, actID, signGame string) DailyRewardContext {
	return DailyRewardContext{
		BaseURL:  baseURL,
		EventID:  eventID,
		ActID:    actID,
		SignGame: signGame,
	}
}

// Claim daily reward.
// e.g. Genshin daily sign in endpoint: https://sg-hk4e-api.hoyolab.com/event/sol/sign?act_id=e202102251931481
func (daily *DailyRepository) Claim(cookie network.Cookie, context DailyRewardContext) (entity.DailyClaim, error) {
	var res entity.DailyClaim

	handler := network.NewHTTPHandler()
	endpoint := fmt.Sprintf("%s/event/%s/%s?act_id=%s", context.BaseURL, context.EventID, DailyRewardSignParam, context.ActID)
	daily.Logger.Info("Claiming daily reward", slog.String("endpoint", endpoint))

	request := network.NewRequest(endpoint, http.MethodPost).
		AddCookie(cookie).
		AddParam("lang", LangEnglish).
		AddHeader("X-Rpc-Signgame", context.SignGame).
		Build()

	err := handler.Send(request, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
