package cmd

import (
	"fmt"
	"log"
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
			log.Fatalf("failed to read model: %q", err)
		}

		text, err := readAll(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read input: %q", err)
		}

		tokIDs, err := vichat.Tokenize(text, model)
		if err != nil {
			log.Fatalf("failed to tokenize: %q", err)
			return
		}

		if verbose, _ := f.GetBool("verbose"); verbose {
			toks, err := vichat.Decode(tokIDs, model)
			if err != nil {
				log.Fatalf("failed to decode: %q", err)
				return
			}

			for i := range tokIDs {
				fmt.Printf("%d:\t%q\n", tokIDs[i], toks[i])
			}
		} else {
			fmt.Println(len(tokIDs))
		}
	},
}

func init() {
	TokCmd.Flags().StringP("model", "m", "gpt-3.5-turbo", "gpt-3.5-turbo|gpt-4")
	TokCmd.Flags().BoolP("verbose", "v", false, "output tokens and their IDs")
}
