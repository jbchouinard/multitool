package errored

import (
	"fmt"
	"os"
)

func Check(err error, msg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Checkf(err error, msg string, args ...any) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg+"\n", args...)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func Fatalf(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
