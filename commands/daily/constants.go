package daily

const (
	// Main API endpoints.
	Hk4eEndpoint     = "https://sg-hk4e-api.hoyolab.com"
	SgPublicEndpoint = "https://sg-public-api.hoyolab.com"

	// Daily reward endpoint parameters.
	DailyRewardHomeParam = "home"
	DailyRewardInfoParam = "info"
	DailyRewardSignParam = "sign"

	// Genshin endpoint parameters.
	GenshinEventId = "sol"
	GenshinActId   = "e202102251931481"

	// Star Rail endpoint parameters.
	StarRailEventId = "luna/os"
	StarRailActId   = "e202303301540311"

	// Zenless endpoint parameters.
	ZenlessEventId = "luna/zzz/os"
	ZenlessActId   = "e202406031448091"

	// Header values for x-rpc-signgame.
	GenshinSignGame  = "genshin"
	StarRailSignGame = "hsr"
	ZenlessSignGame  = "zzz"

	// Language.
	LangEnglish = "en-us"
)
