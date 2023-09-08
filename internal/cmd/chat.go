package cmd

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/JackKCWong/vichat/internal/vichat"
	"github.com/henomis/lingoose/chat"
	"github.com/henomis/lingoose/prompt"
	"github.com/spf13/cobra"
)

var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "read a chat from stdin and send to LLM chat",
	Run: func(cmd *cobra.Command, args []string) {

		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			slog.Error("failed to read input", "err", err)
			return
		}

		messages := chat.New(CreatePrompts(string(input))...)
		chatClient := vichat.New()
		res, err := chatClient.Chat(context.TODO(), messages)
		if err != nil {
			slog.Error("failed", "err", err.Error())
			return
		}

		println(res)
	},
}

func CreatePrompts(text string) []chat.PromptMessage {
	// read text line by line, if a line starts with SYSTEM / AI / USER, create a new prompt of the corresponding type
	lines := strings.Split(text, "\n")
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
