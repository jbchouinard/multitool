package path

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jbchouinard/wmt/errored"
)

var WorkDir string

func getWorkDir() string {
	envDir, ok := os.LookupEnv("WMT_DIR")
	if ok {
		return envDir
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(homeDir, ".wmt")
	} else {
		fmt.Fprint(os.Stderr, "could not get user home dir")
	}

	currentDir, err := os.Getwd()
	errored.Check(err, "path init failed")
	return filepath.Join(currentDir, ".wmt")
}

func init() {
	WorkDir = getWorkDir()
	err := os.MkdirAll(WorkDir, os.ModePerm)
	if err != nil {
		errored.Fatalf("failed to create config dir %q: %s", WorkDir, err)
	}
}
