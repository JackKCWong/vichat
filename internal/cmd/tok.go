package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/JackKCWong/vichat/internal/vichat"
	"github.com/spf13/cobra"
)

var TokCmd = &cobra.Command{
	Use:   "tok",
	Short: "given a piece of text, tok estimate the num of tokens for a given model offline",
	Run: func(cmd *cobra.Command, args []string) {
		// do stuff here
		f := cmd.Flags()
		model, err := f.GetString("model")
		if err != nil {
			slog.Error("failed", "err", err)
		}

		text, err := io.ReadAll(os.Stdin)
		if err != nil {
			slog.Error("failed", "err", err)
		}

		toks, err := vichat.Tokenize(string(text), model)
		if err != nil {
			slog.Error("failed", "err", err)
			return
		}

		fmt.Println(len(toks))
	},
}

func init() {
	TokCmd.Flags().StringP("model", "m", "gpt-3.5-turbo", "gpt-3.5-turbo|gpt-4")
}
