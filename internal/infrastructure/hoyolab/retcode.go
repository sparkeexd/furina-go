package hoyolab

const (
	OK                  = 0     // Request performed successfully.
	InvalidCookie       = -100  // Cookie is invalid/has expired, or user is not logged in.
	DailyAlreadyClaimed = -5003 // Daily check-in rewards are already claimed.
)
