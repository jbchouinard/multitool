package cmd

import (
	"fmt"

	"github.com/jbchouinard/multitool/config"
	"github.com/jbchouinard/multitool/history"
	"github.com/spf13/cobra"
)

var setOptCmd = &cobra.Command{
	Use:   "set",
	Short: "Set option value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		config.SetOption(args[0], args[1])
		history.Add("set", fmt.Sprintf("%s %s", args[0], args[1]))
	},
}

func init() {
	rootCmd.AddCommand(setOptCmd)
}
