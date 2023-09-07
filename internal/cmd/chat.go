package cmd

import (
	"context"
	"log/slog"
	"strings"

	"github.com/JackKCWong/vichat/internal/vichat"
	"github.com/henomis/lingoose/chat"
	"github.com/henomis/lingoose/prompt"
)

func sendChat() {
	messages := chat.New(
		chat.PromptMessage{
			Type:   chat.MessageTypeSystem,
			Prompt: prompt.New("You are a professional joke writer"),
		},
		chat.PromptMessage{
			Type:   chat.MessageTypeUser,
			Prompt: prompt.New("Write a joke about a goose"),
		},
	)

	chatClient := vichat.New()
	res, err := chatClient.Chat(context.TODO(), messages)
	if err != nil {
		slog.Error("failed", "err", err.Error())
		return
	}

	println(res)
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
