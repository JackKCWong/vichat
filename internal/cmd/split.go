package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var SplitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split a text into multiple chunks",
	Run: func(cmd *cobra.Command, args []string) {
		text, err := readAll(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read input: %q", err)
		}
		chunkSize, _ := cmd.Flags().GetInt("chunk-size")
		overlap, _ := cmd.Flags().GetInt("overlap")

		chunks := RecursiveTextSplit(text, chunkSize, overlap)
		for i := range chunks {
			fmt.Println("--------------------------------------------")
			fmt.Println(chunks[i])
		}
	},
}

func init() {
	SplitCmd.Flags().IntP("chunk-size", "n", 1000, "chunk size")
	SplitCmd.Flags().IntP("overlap", "o", 50, "overlap size")
}

func RecursiveTextSplit(text string, chunkSize, overlap int) []string {
	return nil
}
