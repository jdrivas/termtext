package termtext

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jdrivas/vconfig"
	"github.com/spf13/viper"
)

const (
	ScreenProfileKey        = "Termtext.screenProfile" // Viper Key for profile name - string
	ScreenNoColorDefaultKey = "termtextNoColor"        // No Color profile
	ScreenDarkDefaultKey    = "termtextDarkDefault"    // Dark Profile
	ScreenLightDefaultKey   = "termtextLightDefault"   // Light Profile
)

// JSONDisplayKey controls output.
// If set only JSON is displayed
// for renders which have response objects.
// For example, in a yaml configuration file you can set:
// Termtext:
//    JSONDisplay: True
//
const JSONDisplayKey = "Termtext.JSONDisplay" // bool

// Use this with github.com/juju/ansi term to get a TabWriter that works with color.
type ColorSprintfFunc func(string, ...interface{}) string

type Profile struct {
	Name, Key                  string           // Name for display, Key for dictionary and viper
	Title, SubTitle, Text      ColorSprintfFunc // Basic text formating
	Info, Highlight            ColorSprintfFunc // Some semantic formats
	Success, Warn, Fail, Alert ColorSprintfFunc // Alert style formats.
}

type ProfileList []ProfileList

type profileMap map[string]*Profile

var noColorDefaultProfile = &Profile{
	Name:      "Termtext No Color",
	Key:       ScreenDarkDefaultKey,
	Title:     fmt.Sprintf,
	SubTitle:  fmt.Sprintf,
	Text:      fmt.Sprintf,
	Info:      fmt.Sprintf,
	Highlight: fmt.Sprintf,
	Success:   fmt.Sprintf,
	Warn:      fmt.Sprintf,
	Fail:      fmt.Sprintf,
	Alert:     fmt.Sprintf,
}

var profiles = profileMap{
	ScreenNoColorDefaultKey: noColorDefaultProfile,
	ScreenDarkDefaultKey: &Profile{
		Name:      "Termtext Default Dark",
		Key:       ScreenDarkDefaultKey,
		Title:     color.New(color.FgHiWhite).SprintfFunc(),
		SubTitle:  color.New(color.FgWhite).SprintfFunc(),
		Text:      color.New(color.FgWhite).SprintfFunc(),
		Info:      color.New(color.FgWhite).SprintfFunc(),
		Highlight: color.New(color.FgGreen).SprintfFunc(),
		Success:   color.New(color.FgGreen).SprintfFunc(),
		Warn:      color.New(color.FgYellow).SprintfFunc(),
		Fail:      color.New(color.FgRed).SprintfFunc(),
		Alert:     color.New(color.FgRed).SprintfFunc(),
	},
	ScreenLightDefaultKey: &Profile{
		Name:     "Termtext Default Light",
		Key:      ScreenLightDefaultKey,
		Title:    color.New(color.FgBlack).SprintfFunc(),
		SubTitle: color.New(color.FgHiBlack).SprintfFunc(),
		Text:     color.New(color.FgHiBlack).SprintfFunc(),

		// Semantic Formatting
		Info:      color.New(color.FgBlack).SprintfFunc(),
		Highlight: color.New(color.FgGreen).SprintfFunc(),
		Success:   color.New(color.FgGreen).SprintfFunc(),
		Warn:      color.New(color.FgYellow).SprintfFunc(),
		Fail:      color.New(color.FgRed).SprintfFunc(),
		Alert:     color.New(color.FgRed).SprintfFunc(),
	},
}

// ScreenProfile is the current screen profile.
var screenProfile *Profile = noColorDefaultProfile

func InitTerm() {
	if vconfig.Debug() {
		fmt.Printf("Initing Term.\n")
	}

	// See if we have a profile default
	if profile, ok := profiles[viper.GetString(ScreenProfileKey)]; ok {
		setScreenProfile(profile)
	} else {
		setScreenProfile(noColorDefaultProfile)
	}

	if vconfig.Debug() {
		fmt.Printf("Term inited: %s %s %s %s\n", Title("Title"), SubTitle("SubTitle"), Text("Text"), Highlight("Highlight"))
	}
}

// Error formats an error string.
func Error(err error) string {
	return (fmt.Sprintf("%s %s", Title("Error: "), Fail("%v", err.Error())))
}

// I think I like this better than setting up function calls that
// map to a wrapper around the function call off of ScreenProfile
// eg.
// func Title(s string, args ...interface{}) string {
//   ScreenProfile.Title( ...... )
// }
//
// These are used to surround strings with color.
// e.g. fmt.Printf(t.Title("This is a title"))
// or
// fmt.Printf("%s: %s.", t.Alert("Houston we have a problem"), t.Text("is not what was actually said"))
var (
	Title    ColorSprintfFunc
	SubTitle ColorSprintfFunc
	Text     ColorSprintfFunc

	Info      ColorSprintfFunc
	Highlight ColorSprintfFunc
	Success   ColorSprintfFunc
	Warn      ColorSprintfFunc
	Fail      ColorSprintfFunc
	Alert     ColorSprintfFunc
)

func setScreenProfile(p *Profile) {
	screenProfile = p
	Title = p.Title
	SubTitle = p.SubTitle
	Text = p.Text

	// Semantic Formatting
	Info = p.Info
	Highlight = p.Highlight
	Success = p.Success
	Warn = p.Warn
	Fail = p.Fail
	Alert = p.Alert

}
