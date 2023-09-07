package cmd_test

import (
	"testing"

	"github.com/JackKCWong/vichat/internal/cmd"
)

func TestCreatePrompts(t *testing.T) {
	prompts := cmd.CreatePrompts(`SYSTEM: You are
a helpful assistant.

USER: tell me a
joke about goose

AI: why is
the goose a comedian
`)

	if len(prompts) != 3 {
		t.Errorf("Expected 2 prompts, got %d", len(prompts))
	}

	if prompts[0].Type != "system" {
		t.Errorf("Expected prompt type to be system, got %s", prompts[0].Type)
	}

	if prompts[0].Prompt.String() != "You are\na helpful assistant.\n" {
		t.Errorf("Expected prompt string, got %s", prompts[0].Prompt.String())
	}

	if prompts[1].Type != "user" {
		t.Errorf("Expected prompt type to be user, got %s", prompts[1].Type)
	}

	if prompts[2].Type != "assistant" {
		t.Errorf("Expected prompt type to be ai, got %s", prompts[2].Type)
	}
	if prompts[2].Prompt.String() != "why is\nthe goose a comedian\n" {
		t.Errorf("Expected prompt string, got %s", prompts[2].Prompt.String())
	}

}
