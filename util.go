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
	fc, fl, ln := locString(3) // 3 is a magic number: not this func, or the one above ...
	w := ansiterm.NewTabWriter(os.Stdout, 6, 2, 1, ' ', 0)
	fmt.Fprintf(w, df("%s\t%s()\t%s:%d\n", e, fc, fl, ln))
	w.Flush()
}

func locString(d int) (fnc, file string, line int) {
	if pc, fl, l, ok := runtime.Caller(d); ok {
		f := runtime.FuncForPC(pc)
		fnc = filepath.Base(f.Name())
		file = filepath.Base(fl)
		line = l
	}
	return fnc, file, line
}
