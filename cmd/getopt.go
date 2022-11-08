package cmd

import (
	"fmt"

	"github.com/jbchouinard/multitool/config"
	"github.com/spf13/cobra"
)

var getOptCmd = &cobra.Command{
	Use:   "get",
	Short: "Get option value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.GetOption(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(getOptCmd)

}
