package path

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jbchouinard/multitool/errored"
)

var WorkDir string

func getWorkDir() string {
	envDir, ok := os.LookupEnv("MULTITOOL_DIR")
	if ok {
		return envDir
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(homeDir, ".multitool")
	} else {
		fmt.Fprint(os.Stderr, "could not get user home dir")
	}

	currentDir, err := os.Getwd()
	errored.Check(err, "path init failed")
	return filepath.Join(currentDir, ".multitool")
}

func init() {
	WorkDir = getWorkDir()
	err := os.MkdirAll(WorkDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create config dir %q: %s", WorkDir, err)
	}
}
