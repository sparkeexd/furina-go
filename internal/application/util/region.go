package util

type Game int

type Region int

const (
	GenshinImpact Game = iota + 1
	HonkaiStarRail
	ZenlessZoneZero

	RegionAsia Region = iota + 1
	RegionAmerica
	RegionEurope
)
