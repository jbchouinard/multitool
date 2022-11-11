package cmd

import (
	"fmt"

	"github.com/jbchouinard/wmt/config"
	"github.com/jbchouinard/wmt/errored"
	"github.com/spf13/cobra"
)

var optCmd = &cobra.Command{
	Use:   "opt [key] [value]",
	Short: "View or change options",
	Args:  cobra.MaximumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch len(args) {
		case 0:
			options := config.GetAll()
			for _, opt := range options {
				fmt.Printf("%-12s %s\n", opt.Key, opt.Value)
			}
		case 1:
			fmt.Println(config.Get(args[0]))
		case 2:
			k := args[0]
			v := args[1]
			if v == "" || v == "_" {
				err = config.Unset(k)
			} else {
				err = config.Set(k, v)
			}
		}
		errored.Check(err, "")
	},
}

func init() {
	rootCmd.AddCommand(optCmd)
}
