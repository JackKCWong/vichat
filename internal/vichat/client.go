package vichat 

import (
	"os"

	lingoose "github.com/henomis/lingoose/llm/openai"
	"github.com/sashabaranov/go-openai"
)

func New() *lingoose.OpenAI {
	cfg := openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))
	cfg.BaseURL = os.Getenv("OPENAI_API_BASE")
	client := openai.NewClientWithConfig(cfg)

	return lingoose.NewChat().WithClient(client)
}
