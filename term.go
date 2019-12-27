package termtext

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jdrivas/vconfig"
	"github.com/spf13/viper"
)

// ColorSprintfFunc use this with github.com/juju/ansi term to get a TabWriter that works with color.
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

// InitTerm sets profile defaults from viper.
func InitTerm() {
	if vconfig.Debug() {
		Pef()
		defer Pxf()
	}

	// See if a profile has been set, othewise use the default.
	pn := viper.GetString(ScreenProfileKey)
	if profile, ok := profiles[pn]; ok {
		// fmt.Printf("Updating profile to: %s\n", pn)
		setScreenProfile(profile)
	} else {
		// fmt.Printf("Using default profile: %s\n", noColorDefaultProfile.Name)
		setScreenProfile(noColorDefaultProfile)
	}

	if vconfig.Debug() {
		fmt.Printf("Term inited with profile \"%s\": %s %s %s %s\n",
			Title(screenProfile.Name), Title("Title"), SubTitle("SubTitle"), Text("Text"), Highlight("Highlight"))
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
// It's true the set up is more expensive, but it does remove a
// function call from the execution path ......
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

func init() {
	setScreenProfile(noColorDefaultProfile)
}
