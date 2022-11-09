package cmd

import (
	"fmt"

	"github.com/jbchouinard/wmt/config"
	"github.com/jbchouinard/wmt/errored"
	"github.com/jbchouinard/wmt/history"
	"github.com/spf13/cobra"
)

var optCmd = &cobra.Command{
	Use:   "opt",
	Short: "View or change settings",
}

var listOptCmd = &cobra.Command{
	Use:   "list",
	Short: "List all option values",
	Run: func(cmd *cobra.Command, args []string) {
		options := config.GetAll()
		for _, opt := range options {
			fmt.Printf("%-12s %s\n", opt.Key, opt.Value)
		}
	},
}

var getOptCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get option value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Get(args[0]))
	},
}

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
	optCmd.AddCommand(unsetOptCmd)
	optCmd.AddCommand(setOptCmd)
	optCmd.AddCommand(getOptCmd)
	optCmd.AddCommand(listOptCmd)
	rootCmd.AddCommand(optCmd)
}
