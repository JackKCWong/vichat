package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/JackKCWong/vichat/internal/vichat"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/henomis/lingoose/chat"
	"github.com/henomis/lingoose/prompt"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

const DefaultTemperature = 0.7
const DefaultMaxTokens = 1000
const DefaultSystemPrompt = "SYSTEM: You are a helpful assistant."

func init() {
	ChatCmd.Flags().IntP("max_tokens", "m", DefaultMaxTokens, "max token for response")
	ChatCmd.Flags().Float32P("temperature", "t", DefaultTemperature, "temperature, higher means more randomness.")
	ChatCmd.Flags().BoolP("term", "o", false, "print to terminal")
	// ChatCmd.Flags().StringP("role", "r", "assistant", "each role maps to a system prompt defined in ~/.vichat.yaml")
}

var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "read a chat from stdin and send to LLM chat",
	Run: func(cmd *cobra.Command, args []string) {

		input := ""
		if len(args) == 0 {
			stdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				slog.Error("failed to read input", "err", err)
				return
			}

			input = string(stdin)
		} else {
			input = strings.Join(args, " ")
		}

		var f = cmd.Flags()
		var temperature float32 = DefaultTemperature
		var maxTokens int = DefaultMaxTokens

		lines := strings.Split(string(input), "\n")
		if strings.HasPrefix(lines[0], "#") {
			temperature = getTemperature(lines[0])
			maxTokens = getMaxTokens(lines[0])
			lines = lines[0:]
		} else {
			if m, err := f.GetInt("max_tokens"); err == nil {
				maxTokens = m
			}

			if t, err := f.GetFloat32("temperature"); err == nil {
				temperature = t
			}
		}

		chatClient := vichat.New().WithTemperature(temperature).WithMaxTokens(maxTokens)
		prompts := CreatePrompts(lines)
		if len(prompts) == 0 {
			slog.Error("invalid input")
			return
		}

		if prompts[0].Type != chat.MessageTypeSystem {
			prompts = append([]chat.PromptMessage{{
				Type:   chat.MessageTypeSystem,
				Prompt: prompt.New(DefaultSystemPrompt),
			}}, prompts...)
		}

		messages := chat.New(prompts...)
		res, err := chatClient.Chat(context.TODO(), messages)
		if err != nil {
			slog.Error("failed", "err", err.Error())
			return
		}

		term, _ := f.GetBool("term")

		if term && isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd()) {
			if strings.Contains(res, "```") {
				fmt.Println()
				fmt.Println(string(markdown.Render(res, 90, 4)))
				return
			}
		}

		fmt.Println(res)
	},
}

func CreatePrompts(lines []string) []chat.PromptMessage {
	prompts := make([]chat.PromptMessage, 0)
	var messageType chat.MessageType = chat.MessageTypeUser
	var message strings.Builder
	for _, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "SYSTEM: ") {
			if messageType != "" {
				prompts = append(prompts, chat.PromptMessage{
					Type:   messageType,
					Prompt: prompt.New(message.String()),
				})
				message.Reset()
			}
			messageType = chat.MessageTypeSystem
			message.WriteString(line[8:])
			message.WriteString("\n")
		} else if strings.HasPrefix(line, "AI: ") {
			if messageType != "" {
				prompts = append(prompts, chat.PromptMessage{
					Type:   messageType,
					Prompt: prompt.New(message.String()),
				})
				message.Reset()
			}
			messageType = chat.MessageTypeAssistant
			message.WriteString(line[4:])
			message.WriteString("\n")
		} else if strings.HasPrefix(line, "USER: ") {
			if messageType != "" {
				prompts = append(prompts, chat.PromptMessage{
					Type:   messageType,
					Prompt: prompt.New(message.String()),
				})
				message.Reset()
			}
			messageType = chat.MessageTypeUser
			message.WriteString(line[6:])
			message.WriteString("\n")
		} else {
			message.WriteString(line)
			message.WriteString("\n")
		}
	}

	prompts = append(prompts, chat.PromptMessage{
		Type:   messageType,
		Prompt: prompt.New(message.String()),
	})

	return prompts
}

func getTemperature(text string) float32 {
	// match temperature from a string using regex pattern temperature\W+([\d\.]+), extract the matched group
	// and assign it to temperature variable.
	re := regexp.MustCompile(`temperature\W+([\d\.]+)`)
	kv := re.FindStringSubmatch(text)
	if len(kv) > 1 {
		t, err := strconv.ParseFloat(kv[1], 32)
		if err == nil {
			return float32(t)
		}
	}

	return DefaultTemperature
}

func getMaxTokens(text string) int {
	// match max_tokens from a string using regex pattern max_tokens\W+([\d\.]+), extract the matched group
	// and assign it to temperature variable.
	re := regexp.MustCompile(`max_tokens\W+(\d+)`)
	kv := re.FindStringSubmatch(text)
	if len(kv) > 1 {
		t, err := strconv.Atoi(kv[1])
		if err == nil {
			return t
		}
	}

	return DefaultMaxTokens
}
