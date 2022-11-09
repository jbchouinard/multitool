package errored

import "fmt"

func Check(err error, msg string) {
	if err != nil {
		fmt.Print(msg)
		panic(err)
	}
}

func Checkf(err error, msg string, args ...any) {
	fmt.Printf(msg, args...)
	panic(err)
}
