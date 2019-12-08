package termtext

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jdrivas/vconfig"
	"github.com/spf13/viper"
)

// const ScreenProfileKey = "screenProfile"
// const ScreenDarkValue = "dark"

// Use this with github.com/juju/ansi term to get a TabWriter that works with color.
type ColorSprintfFunc func(string, ...interface{}) string

const (
	ScreenProfileKey = "screenProfileKey"
	ScreenDarkValue  = "dark"
)

var (
	// Text Formatting
	Title    = color.New(color.FgBlack).SprintfFunc()
	SubTitle = color.New(color.FgHiBlack).SprintfFunc()
	Text     = color.New(color.FgHiCyan).SprintfFunc()

	// Semantic Formatting
	Info      = color.New(color.FgBlack).SprintfFunc()
	Highlight = color.New(color.FgGreen).SprintfFunc()
	Success   = color.New(color.FgGreen).SprintfFunc()
	Warn      = color.New(color.FgYellow).SprintfFunc()
	Fail      = color.New(color.FgRed).SprintfFunc()
	Alert     = color.New(color.FgRed).SprintfFunc()
)

var screen = screenDarkValue

func InitTerm() {
	if vconfig.Debug() {
		fmt.Printf("Initing Term.\n")
	}
	screen = viper.GetString(ScreenProfileKey)
	if screen == ScreenDarkValue {
		fmt.Printf("Doing Dark collors.\n")
		// Text Formatting
		Title = color.New(color.FgHiWhite).SprintfFunc()
		SubTitle = color.New(color.FgWhite).SprintfFunc()
		Text = color.New(color.FgWhite).SprintfFunc()

		// Semantic Formatting
		Info = color.New(color.FgWhite).SprintfFunc()
		Highlight = color.New(color.FgGreen).SprintfFunc()
		Success = color.New(color.FgGreen).SprintfFunc()
		Warn = color.New(color.FgYellow).SprintfFunc()
		Fail = color.New(color.FgRed).SprintfFunc()
		Alert = color.New(color.FgRed).SprintfFunc()

	}

	if vconfig.Debug() {
		fmt.Printf("Term inited: %s %s %s %s\n", Title("Title"), SubTitle("SubTitle"), Text("Text"), Highlight("Highlight"))
	}
}

// Error formats an error string.
func Error(err error) string {
	return (fmt.Sprintf("%s %s", Title("Error: "), Fail("%v", err.Error())))
}
