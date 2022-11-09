package cmd

import (
	"fmt"
	"os"

	"github.com/jbchouinard/multitool/config"
	"github.com/jbchouinard/multitool/history"
	"github.com/spf13/cobra"
)

var setOptCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set option value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Set(args[0], args[1]); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		history.Add("set", fmt.Sprintf("%s %s", args[0], args[1]))
	},
}

func init() {
	rootCmd.AddCommand(setOptCmd)
}
