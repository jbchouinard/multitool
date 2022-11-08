package path

import (
	"log"
	"os"
	"path/filepath"
)

var WorkDir string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Print("could not get user home dir, using current dir")
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		WorkDir = filepath.Join(currentDir, ".multitool")
	} else {
		WorkDir = filepath.Join(homeDir, ".multitool")
	}
	os.MkdirAll(WorkDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create config dir %q: %s", WorkDir, err)
	}
}
