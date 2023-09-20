package cmd

import (
	"strings"
	"testing"
)

func TestCreatePrompts(t *testing.T) {
	script := `SYSTEM: You are
a helpful assistant.

USER: tell me a
joke about goose

AI: why is
the goose a comedian
`
	prompts := CreatePrompts(strings.Split(script, "\n"))

	if len(prompts) != 3 {
		t.Errorf("Expected 3 prompts, got %d", len(prompts))
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

	script = `tell me a joke about vim`
	prompts = CreatePrompts(strings.Split(script, "\n"))
	if len(prompts) != 1 {
		t.Errorf("Expected 1 prompts, got %d", len(prompts))
	}

	if prompts[0].Type != "user" {
		t.Errorf("Expected prompt type to be system, got %s", prompts[0].Type)
	}
}

func TestGetLLMParams(t *testing.T) {
	cfg1 := "# temperature=0.1, max_tokens=100"
	if getTemperature(cfg1) != 0.1 {
		t.Errorf("Expected temperature to be 0.1, got %f", getTemperature(cfg1))
	}

	if getMaxTokens(cfg1) != 100 {
		t.Errorf("Expected max_tokens to be 100, got %d", getMaxTokens(cfg1))
	}

	cfg2 := "# temperature: 0    max_tokens: 200"
	if getTemperature(cfg2) != 0 {
		t.Errorf("Expected temperature to be 0, got %f", getTemperature(cfg2))
	}

	if getMaxTokens(cfg2) != 200 {
		t.Errorf("Expected max_tokens to be 200, got %d", getMaxTokens(cfg2))
	}
}
