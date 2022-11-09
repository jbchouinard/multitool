package path

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jbchouinard/multitool/errored"
)

var WorkDir string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Print("could not get user home dir, using current dir")
		currentDir, err := os.Getwd()
		errored.Check(err, "path init")
		WorkDir = filepath.Join(currentDir, ".multitool")
	} else {
		WorkDir = filepath.Join(homeDir, ".multitool")
	}
	os.MkdirAll(WorkDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create config dir %q: %s", WorkDir, err)
	}
}
