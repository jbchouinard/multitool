package clip

import (
	"fmt"

	"github.com/jbchouinard/multitool/config"
	"golang.design/x/clipboard"
)

var enabled bool

func init() {
	enabled = config.GetOption("clipboard") == "yes"

	if enabled {
		err := clipboard.Init()
		if err != nil {
			fmt.Println("Clipboard failed to init")
			fmt.Println("Disable with: multitool set clipboard no")
			panic(err)
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
