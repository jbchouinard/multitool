package cmd

import (
	"fmt"
	"os"

	"github.com/jbchouinard/multitool/config"
	"github.com/jbchouinard/multitool/history"
	"github.com/spf13/cobra"
)

var unsetOptCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Un-set option value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Unset(args[0]); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		history.Add("unset", args[0])
	},
}

func init() {
	rootCmd.AddCommand(unsetOptCmd)
}
