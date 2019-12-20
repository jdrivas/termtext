package termtext

// JSONDisplayKey controls output.
// If set only JSON is displayed
// for renders which have response objects.
// For example, in a yaml configuration file you can set:
// Termtext:
//    JSONDisplay: True
//
const JSONDisplayKey = "Termtext.JSONDisplay" // bool

// This is how you pick a screen profile.

// Variables
const (
	ScreenProfileKey = "Termtext.screenProfile" // Viper Key for profile name - string
)

// Values
const (
	ScreenNoColorDefaultKey = "termtextNoColor"      // No Color profile
	ScreenDarkDefaultKey    = "termtextDarkDefault"  // Dark Profile
	ScreenLightDefaultKey   = "termtextLightDefault" // Light Profile
)
