package cmd

import (
	"fmt"

	"github.com/jbchouinard/wmt/env"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Env variables commands",
}

var envGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get env variables",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		val, _ := useEnv().Get(args[0])
		fmt.Println(val)
	},
}

var envSetCmd = &cobra.Command{
	Use:   "set KEY VALUE",
	Short: "Set env variables",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		useEnv().Set(args[0], args[1])
		useEnv().Save(false)
	},
}

var envUnsetCmd = &cobra.Command{
	Use:   "unset KEY",
	Short: "Unset env variable",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		useEnv().Unset(args[0])
		useEnv().Save(false)
	},
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List environment variables",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		vars := useEnv().List()
		for _, v := range vars {
			fmt.Println(v)
		}
	},
}

var envCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Print name of current env",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(useEnv().Name())
	},
}

var envUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Set current env in use (default: global)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			env.SetCurrentEnv(args[0])
		} else {
			env.SetCurrentEnv("global")
		}
	},
}

var useGlobal bool

func useEnv() *env.Env {
	if useGlobal {
		return env.Global
	} else {
		return env.Current
	}
}

func init() {
	envCmd.AddCommand(envGetCmd)
	envCmd.AddCommand(envSetCmd)
	envCmd.AddCommand(envUnsetCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envUseCmd)
	envCmd.AddCommand(envCurrentCmd)

	envCmd.PersistentFlags().BoolVarP(&useGlobal, "global", "g", false, "on global env")
	rootCmd.AddCommand(envCmd)
}
