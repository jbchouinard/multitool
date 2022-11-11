package cmd

import (
	"time"

	"github.com/jbchouinard/wmt/history"
	"github.com/spf13/cobra"
)

var keepDays uint16

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge command history",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		asOf := time.Now().UTC().AddDate(0, 0, -int(keepDays))
		history.Purge(asOf)
	},
}

func init() {
	purgeCmd.Flags().Uint16Var(&keepDays, "keep-days", 7, "days of history to keep")
	rootCmd.AddCommand(purgeCmd)
}
