package errored

import (
	"fmt"
	"os"
)

func Check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		os.Exit(1)
	}
}

func Checkf(err error, msg string, args ...any) {
	if err != nil {
		fmt.Printf(msg+"\n", args...)
		fmt.Println(err)
		os.Exit(1)
	}
}
