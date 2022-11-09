package cmd

import (
	"fmt"

	"github.com/jbchouinard/multitool/config"
	"github.com/spf13/cobra"
)

var getOptCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get option value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "all" {
			options := config.GetAll()
			for _, opt := range options {
				fmt.Printf("%-12s %s\n", opt.Key, opt.Value)
			}
		} else {
			fmt.Println(config.Get(args[0]))
		}
	},
}

func init() {
	rootCmd.AddCommand(getOptCmd)

}
