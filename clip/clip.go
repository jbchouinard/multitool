package clip

import (
	"fmt"

	"github.com/jbchouinard/wmt/config"
	"golang.design/x/clipboard"
)

var enabled bool

func init() {
	config.ValidValues["clipboard"] = map[string]bool{"yes": true, "no": true}
	config.DefaultValues["clipboard"] = "yes"
	enabled = config.Get("clipboard") == "yes"

	if enabled {
		err := clipboard.Init()
		if err != nil {
			fmt.Println("Clipboard failed to init")
			fmt.Println("Disable with: wmt set clipboard no")
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
