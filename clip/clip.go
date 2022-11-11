package clip

import (
	"fmt"
	"os"

	"github.com/jbchouinard/wmt/config"
	"golang.design/x/clipboard"
)

var enabled bool

func init() {
	config.ValidValues["clipboard"] = map[string]bool{"yes": true, "no": true}
	config.DefaultValues["clipboard"] = "no"
	enabled = config.Get("clipboard") == "yes"

	if enabled {
		err := clipboard.Init()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Clipboard integration failed to init, disabling")
			fmt.Println("Re-enable with: wmt opt clipboard yes")
			config.Set("clipboard", "no")
			os.Exit(1)
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
