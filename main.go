package main

import (
	"log/slog"
	"os"

	"github.com/JackKCWong/vichat/internal/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vichat",
	Short: "vichat is a LLM chat cli",
	Run: func(c *cobra.Command, args []string) {
		cmd.ChatCmd.Run(c, args)
	},
}

func init() {
	rootCmd.Flags().AddFlagSet(cmd.ChatCmd.Flags())
	rootCmd.AddCommand(cmd.TokCmd)
	rootCmd.AddCommand(cmd.ChatCmd)
	rootCmd.AddCommand(cmd.InstallCmd)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := rootCmd.Execute(); err != nil {
		println(err.Error())
	}
}
