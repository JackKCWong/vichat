package main

import (
	"github.com/JackKCWong/vichat/internal/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vichat",
	Short: "vichat is a LLM chat cli",
	Args:  cobra.MinimumNArgs(0),
	Run: func(c *cobra.Command, args []string) {
		cmd.ChatCmd.Run(c, args)
	},
}

func init() {
	rootCmd.Flags().AddFlagSet(cmd.ChatCmd.Flags())
	rootCmd.AddCommand(cmd.TokCmd)
	rootCmd.AddCommand(cmd.ChatCmd)
	rootCmd.AddCommand(cmd.InstallCmd)
	rootCmd.AddCommand(cmd.SplitCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		println(err.Error())
	}
}
