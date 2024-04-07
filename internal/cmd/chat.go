package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mattn/go-isatty"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"

	_ "embed"
	"encoding/csv"
)

const DefaultTemperature = 0.7
const DefaultMaxTokens = 1000
const DefaultSystemPrompt = "You are a helpful assistant, you help people by answering their questions politely and precisely."

//go:embed prompts.csv
var awesomePrompts []byte

func init() {
	ChatCmd.Flags().IntP("max_tokens", "m", DefaultMaxTokens, "max token for response")
	ChatCmd.Flags().Float32P("temperature", "t", DefaultTemperature, "temperature, higher means more randomness.")
	ChatCmd.Flags().BoolP("render", "r", false, "render markdown to terminal")
	ChatCmd.Flags().BoolP("func", "f", false, "use functions")
	ChatCmd.Flags().StringP("system-prompt", "p", "assistant", "point to a system prompt file")
	ChatCmd.Flags().StringP("outdir", "o", ".", "dir to keep chat history")
	ChatCmd.Flags().BoolP("stream", "s", false, "use streaming response")
}

var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "read a chat from stdin and send to LLM chat",
	Run: func(cmd *cobra.Command, args []string) {

		if hasDifference(vimPlugins, "vim/ftdetect", os.ExpandEnv("$HOME/.vim/ftdetect")) {
			log.Fatalf("vichat vim plugin appears to be out of sync. run `vichat i` to install it again.")
		}

		if hasDifference(vimPlugins, "vim/ftplugin", os.ExpandEnv("$HOME/.vim/ftplugin")) {
			log.Fatalf("vichat vim plugin appears to be out of sync. run `vichat i` to install it again.")
		}

		if hasDifference(vimPlugins, "vim/syntax", os.ExpandEnv("$HOME/.vim/syntax")) {
			log.Fatalf("vichat vim plugin appears to be out of sync. run `vichat i` to install it again.")
		}

		var opts = cmd.Flags()
		var temperature float32 = DefaultTemperature
		var maxTokens int = DefaultMaxTokens

		if m, err := opts.GetInt("max_tokens"); err == nil {
			maxTokens = m
		}

		if t, err := opts.GetFloat32("temperature"); err == nil {
			temperature = t
		}

		var input string
		var lines []string
		var isSimpleChat bool = false
		if !isatty.IsTerminal(os.Stdin.Fd()) {
			stdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("failed to read input: %q", err)
				return
			}

			input = string(stdin)
			lines = strings.Split(string(input), "\n")
			if strings.HasPrefix(lines[0], "#") {
				temperature = getTemperature(lines[0])
				maxTokens = getMaxTokens(lines[0])
				lines = lines[0:]
			}
		} else {
			input = strings.Join(args, " ")
			lines = []string{input}
			isSimpleChat = true
		}

		cfg := openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))
		cfg.BaseURL = os.Getenv("OPENAI_API_BASE")
		llmClient := openai.NewClientWithConfig(cfg)

		messages := CreatePrompts(lines)
		if len(messages) == 0 {
			log.Fatalf("invalid input")
			return
		}

		if isSimpleChat {
			var promptStr []byte
			var err error
			optPormpt, _ := opts.GetString("system-prompt")
			if optPormpt == "assistant" {
				promptStr = []byte(DefaultSystemPrompt)
			} else {
				promptStr, err = os.ReadFile(optPormpt)
				if err != nil {
					prd := csv.NewReader(bytes.NewReader(awesomePrompts))
					embedPrompts, err := prd.ReadAll()
					if err == nil {
						index := make([]string, len(embedPrompts))
						for i := range embedPrompts {
							index[i] = strings.ToLower(embedPrompts[i][0])
						}

						matches := fuzzy.RankFind(optPormpt, index)
						sort.Sort(matches)

						hit := matches[0].OriginalIndex
						promptStr = []byte(embedPrompts[hit][1])
					}
				}
			}

			messages = append([]openai.ChatCompletionMessage{{
				Role:    openai.ChatMessageRoleSystem,
				Content: string(promptStr),
			}}, messages...)
		}

		if temperature == 0 {
			// workaround https://github.com/sashabaranov/go-openai?tab=readme-ov-file#why-dont-we-get-the-same-answer-when-specifying-a-temperature-field-of-0-and-asking-the-same-question
			temperature = math.SmallestNonzeroFloat32
		}

		isRenderOutput, _ := opts.GetBool("render")
		stream, _ := opts.GetBool("stream")

		if isRenderOutput {
			resp, err := llmClient.CreateChatCompletion(context.Background(),
				openai.ChatCompletionRequest{
					Model:       openai.GPT3Dot5Turbo,
					Temperature: temperature,
					MaxTokens:   maxTokens,
					Messages:    messages,
					Stream:      false,
				})
			if err != nil {
				log.Fatalf("failed to send chat: %q", err.Error())
				return
			}

			fmt.Printf("\n%s\n", markdown.Render(resp.Choices[0].Message.Content, 90, 4))
		} else {
			if isSimpleChat {
				// open the full chat in vim
				dir, err := opts.GetString("outdir")
				if err != nil {
					dir = os.TempDir()
				}

				nonWords := regexp.MustCompile(`\W+`)
				filename := nonWords.ReplaceAllString(input, "_")
				if len(filename) > 50 {
					filename = filename[:50]
				}

				tmpf, err := os.CreateTemp(dir, fmt.Sprintf("%s-*.chat", filename))
				if err != nil {
					log.Fatalf("failed to create temp file: %q", err)
				}

				fmt.Fprintf(tmpf, "# temperature=%.1f, max_tokens=%d\n\n", temperature, maxTokens)
				printPrompts(tmpf, messages)
				tmpf.Close()

				// invoke vim using cmd and open tmpf
				var cmd *exec.Cmd
				if input == "" {
					cmd = exec.Command("vim", "-c", "norm! GkA", tmpf.Name())
				} else {
					if stream {
						cmd = exec.Command("vim", "-c", "redraw|call SetStream(1)|Chat", tmpf.Name())
					} else {
						cmd = exec.Command("vim", "-c", "redraw|call SetStream(0)|Chat", tmpf.Name())
					}
				}

				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout

				cmd.Run()
			} else {
				if stream {
					resp, err := llmClient.CreateChatCompletionStream(context.Background(),
						openai.ChatCompletionRequest{
							Model:       openai.GPT3Dot5Turbo,
							Temperature: temperature,
							MaxTokens:   maxTokens,
							Messages:    messages,
						})

					if err != nil {
						log.Fatalf("failed to stream chat: %q", err.Error())
						return
					}

					defer resp.Close()
					for {
						response, err := resp.Recv()
						if errors.Is(err, io.EOF) {
							fmt.Println()
							return
						}

						if err != nil {
							fmt.Printf("\nStream error: %v\n", err)
							return
						}

						fmt.Printf(response.Choices[0].Delta.Content)
					}

				} else {
					resp, err := llmClient.CreateChatCompletion(context.Background(),
						openai.ChatCompletionRequest{
							Model:       openai.GPT3Dot5Turbo,
							Temperature: temperature,
							MaxTokens:   maxTokens,
							Messages:    messages,
							Stream:      false,
						})
					if err != nil {
						log.Fatalf("failed to send chat: %q", err.Error())
						return
					}

					res := resp.Choices[0].Message.Content
					if isRenderOutput {
						fmt.Printf("\n%s\n", markdown.Render(res, 90, 4))
					} else {
						fmt.Printf("%s\n\n", res)
					}
				}
			}
		}
	},
}

func CreatePrompts(lines []string) []openai.ChatCompletionMessage {
	prompts := make([]openai.ChatCompletionMessage, 0)
	var messageType = ""
	var message strings.Builder
	for _, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "SYSTEM: ") {
			if messageType != "" {
				prompts = append(prompts, openai.ChatCompletionMessage{
					Role:    messageType,
					Content: message.String(),
				})
				message.Reset()
			}
			messageType = openai.ChatMessageRoleSystem
			message.WriteString(line[8:])
			message.WriteString("\n")
		} else if strings.HasPrefix(line, "AI: ") {
			if messageType != "" {
				prompts = append(prompts, openai.ChatCompletionMessage{
					Role:    messageType,
					Content: message.String(),
				})
				message.Reset()
			}
			messageType = openai.ChatMessageRoleAssistant
			message.WriteString(line[4:])
			message.WriteString("\n")
		} else if strings.HasPrefix(line, "USER: ") {
			if messageType != "" {
				prompts = append(prompts, openai.ChatCompletionMessage{
					Role:    messageType,
					Content: message.String(),
				})
				message.Reset()
			}
			messageType = openai.ChatMessageRoleUser
			message.WriteString(line[6:])
			message.WriteString("\n")
		} else {
			message.WriteString(line)
			message.WriteString("\n")
		}
	}

	if messageType == "" {
		messageType = openai.ChatMessageRoleUser
	}

	prompts = append(prompts, openai.ChatCompletionMessage{
		Role:    messageType,
		Content: message.String(),
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

type TimeQuery struct {
	SecondsAgo int `json:"secondsAgo"`
}

type TimeResp struct {
	Time  string `json:"time"`
	Query string `json:"query"`
}

func getRelativeTime(query TimeQuery) TimeResp {
	return TimeResp{Time: time.Now().
		Add(time.Duration(-query.SecondsAgo) * time.Second).
		Format("2006-01-02 15:04:05"),
		Query: (time.Duration(query.SecondsAgo) * time.Second).String() + " ago",
	}
}

func printPrompts(w io.Writer, prompts []openai.ChatCompletionMessage) {
	for _, p := range prompts {
		prefix := ""
		switch p.Role {
		case openai.ChatMessageRoleSystem:
			prefix = "SYSTEM"
		case openai.ChatMessageRoleUser:
			prefix = "USER"
		case openai.ChatMessageRoleAssistant:
			prefix = "AI"
		}

		fmt.Fprintf(w, "%s: %s\n\n", prefix, strings.Trim(p.Content, "\r\n"))
	}
}
