package cmd

import (
	"context"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/JackKCWong/vichat/internal/vichat"
	"github.com/henomis/lingoose/chat"
	"github.com/henomis/lingoose/prompt"
	"github.com/spf13/cobra"
)

const DefaultTemperature = 0.7
const DefaultMaxTokens = 256

var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "read a chat from stdin and send to LLM chat",
	Run: func(cmd *cobra.Command, args []string) {

		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			slog.Error("failed to read input", "err", err)
			return
		}

		lines := strings.Split(string(input), "\n")
		var temperature float32 = DefaultTemperature
		maxTokens := DefaultMaxTokens
		if strings.HasPrefix(lines[0], "#") {
			temperature = getTemperature(lines[0])
			maxTokens = getMaxTokens(lines[0])
			lines = lines[0:]
		}

		chatClient := vichat.New().WithTemperature(temperature).WithMaxTokens(maxTokens)
		messages := chat.New(CreatePrompts(lines)...)
		res, err := chatClient.Chat(context.TODO(), messages)
		if err != nil {
			slog.Error("failed", "err", err.Error())
			return
		}

		println(res)
	},
}

func CreatePrompts(lines []string) []chat.PromptMessage {
	prompts := make([]chat.PromptMessage, 0)
	var messageType chat.MessageType
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
