package errored

import (
	"fmt"
	"os"
)

func Check(err error, msg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg)
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func Checkf(err error, msg string, args ...any) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg+"\n", args...)
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
