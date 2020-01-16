package termtext

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/juju/ansiterm"
)

// Pef is an entry tracing printout to stdout. It will print the word Enter, a funciton name, file name and line number.
// It is designed to be used as the first line of a function, perhaps bracketed by a debug check.
func Pef() {
	printLString("Enter", Success)
}

// Pxf is an exit tracing pritout to stdout. It will print out the word exit, a function name, file anda line number.
// It is designed to be used just after a call to Pef() with a defer.
// Pef()
// defer Pxf()
func Pxf() {
	printLString("Exit", Alert)
}

func printLString(e string, df ColorSprintfFunc) {
	fc, fl, ln := loc(2) // not this function, but or the caller, but the callers caller.
	w := ansiterm.NewTabWriter(os.Stdout, 6, 2, 1, ' ', 0)
	fmt.Fprintf(w, df("%s\t%s()\t%s:%d\n", e, fc, fl, ln))
	w.Flush()
}

// LocString string for function and file location of depth d.
// 0 is depth of the function calling LocString, 1 is the callers caller etc.
func LocString(d int) string {
	fc, fl, ln := loc(d + 1)
	return fmt.Sprintf("%s() %s:%d", fc, fl, ln)
}

func loc(d int) (fnc, file string, line int) {
	if pc, fl, l, ok := runtime.Caller(d + 1); ok { // d is relative to the calling function.
		f := runtime.FuncForPC(pc)
		fnc = filepath.Base(f.Name())
		file = filepath.Base(fl)
		line = l
	}
	return fnc, file, line
}
