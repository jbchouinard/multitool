package cmd

import (
	"fmt"

	"github.com/jbchouinard/wmt/history"
	"github.com/spf13/cobra"
)

var count int

var historyCmd = &cobra.Command{
	Use:   "history command",
	Short: "Show history",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entries := history.GetLast(args[0], count)
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("%s: %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"), entry.Value)
		}
	},
}

func init() {
	historyCmd.Flags().IntVarP(&count, "count", "n", 10, "max entries shown")
	rootCmd.AddCommand(historyCmd)
}
