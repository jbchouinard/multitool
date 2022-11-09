package cmd

import (
	"fmt"

	"github.com/jbchouinard/wmt/config"
	"github.com/jbchouinard/wmt/errored"
	"github.com/jbchouinard/wmt/history"
	"github.com/spf13/cobra"
)

var setOptCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set option value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Set(args[0], args[1])
		errored.Check(err, "set "+args[0])
		history.Add("set", fmt.Sprintf("%s %s", args[0], args[1]))
	},
}

func init() {
	rootCmd.AddCommand(setOptCmd)
}
