package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/JackKCWong/vichat/internal/vichat"
	"github.com/spf13/cobra"
)

var SplitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split a text into multiple chunks",
	Run: func(cmd *cobra.Command, args []string) {
		text, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read input: %q", err)
		}
		chunkSize, _ := cmd.Flags().GetInt("chunk-size")
		overlap, _ := cmd.Flags().GetInt("overlap")

		chunks := vichat.RecursiveTextSplit(string(text), chunkSize, overlap)
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
