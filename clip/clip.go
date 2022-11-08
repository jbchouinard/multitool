package clip

import (
	"log"

	"github.com/jbchouinard/multitool/config"
	"golang.design/x/clipboard"
)

var enabled bool

func init() {
	enabled = config.GetOption("clipboard") == "yes"

	if enabled {
		err := clipboard.Init()
		if err != nil {
			log.Fatalf("Clipboard failed to init: %s\nDisable with: multitool set clipboard no", err)
		}
	}

}

func WriteBytes(content []byte) {
	if enabled {
		done := clipboard.Write(clipboard.FmtText, content)
		<-done
	}
}

func Write(content string) {
	if enabled {
		WriteBytes([]byte(content))
	}
}
