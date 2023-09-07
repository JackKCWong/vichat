package cmd

import (
	"context"
	"log/slog"

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
		slog.Error("err", err.Error())
		return
	}

	println(res)
}
