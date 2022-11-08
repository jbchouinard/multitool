package errored

import "log"

func Fatal(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
