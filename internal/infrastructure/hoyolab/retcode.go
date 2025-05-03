package hoyolab

const (
	InvalidCookie       Retcode = -100  // Cookie is invalid/has expired, or user is not logged in.
	DailyAlreadyClaimed Retcode = -5003 // Daily check-in rewards are already claimed.
)

// HoYoLAB return codes that come from 200 OK responses.
type Retcode int
