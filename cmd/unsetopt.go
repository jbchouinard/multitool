package cmd

import (
	"github.com/jbchouinard/wmt/config"
	"github.com/jbchouinard/wmt/errored"
	"github.com/jbchouinard/wmt/history"
	"github.com/spf13/cobra"
)

var unsetOptCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Un-set option value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Unset(args[0])
		errored.Check(err, "unset "+args[0])
		history.Add("unset", args[0])
	},
}

func init() {
	rootCmd.AddCommand(unsetOptCmd)
}
