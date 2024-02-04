package iostream

import "github.com/muesli/termenv"

var profile = termenv.ColorProfile()

var (
	DarkThemeBlack       = profile.Color("#282828")
	DarkThemeBlackDeep   = profile.Color("#1d2021")
	DarkThemeRedLight    = profile.Color("#fb4934")
	DarkThemeRedDeep     = profile.Color("#cc241d")
	DarkThemeYellow      = profile.Color("#d79921")
	DarkThemeYellowLight = profile.Color("#fabd2f")
	DarkThemeGreen       = profile.Color("#98971a")
	DarkThemeGreenLight  = profile.Color("#b8bb26")
	DarkThemeOrange      = profile.Color("#d65d0e")
	DarkThemeOrangeLight = profile.Color("#fe8019")
	DarkThemeBlue        = profile.Color("#458588")
	DarkThemeBlueLight   = profile.Color("#83a598")
	DarkThemePurple      = profile.Color("#b16286")
	DarkThemePurpleLight = profile.Color("#d3869b")
	DarkThemeAqua        = profile.Color("#689d6a")
	DarkThemeAquaLight   = profile.Color("#8ec07c")
)
