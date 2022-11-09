package editor

import (
	"os"
	"os/exec"
	"strings"

	"github.com/jbchouinard/wmt/config"
	"github.com/jbchouinard/wmt/errored"
)

var editorCommand string
var editorArgs []string

func init() {
	config.DefaultValues["editor"] = "nano"
	editor := strings.Split(config.Get("editor"), " ")
	if len(editor) == 0 {
		errored.Fatal("invalid editor command")
	}
	editorCommand = editor[0]
	editorArgs = editor[1:]
}

func Edit(filename string, create bool) error {
	if create {
		if f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
			return err
		} else {
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	args := append(editorArgs, filename)
	return exec.Command(editorCommand, args...).Run()
}
