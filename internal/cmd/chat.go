package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/JackKCWong/vichat/internal/vichat"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/henomis/lingoose/chat"
	"github.com/henomis/lingoose/prompt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	_ "embed"
	"encoding/csv"
)

const DefaultTemperature = 0.7
const DefaultMaxTokens = 1000
const DefaultSystemPrompt = "You are a helpful assistant."

//go:embed prompts.csv
var awesomePrompts []byte

func init() {
	ChatCmd.Flags().IntP("max_tokens", "m", DefaultMaxTokens, "max token for response")
	ChatCmd.Flags().Float32P("temperature", "t", DefaultTemperature, "temperature, higher means more randomness.")
	ChatCmd.Flags().BoolP("render", "r", false, "render markdown to terminal")
	ChatCmd.Flags().BoolP("func", "f", false, "use functions")
	ChatCmd.Flags().StringP("system-prompt", "s", "system.prompt", "point to a system prompt file")
}

var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "read a chat from stdin and send to LLM chat",
	Run: func(cmd *cobra.Command, args []string) {

		if hasDifference(vimPlugins, "vim/ftdetect", os.ExpandEnv("$HOME/.vim/ftdetect")) {
			log.Printf("WARNING: your vichat vim plugin appears to be out of sync. run `vichat i` to install it again.")
		}

		if hasDifference(vimPlugins, "vim/ftplugin", os.ExpandEnv("$HOME/.vim/ftplugin")) {
			log.Printf("WARNING: your vichat vim plugin appears to be out of sync. run `vichat i` to install it again.")
		}

		if hasDifference(vimPlugins, "vim/syntax", os.ExpandEnv("$HOME/.vim/syntax")) {
			log.Printf("WARNING: your vichat vim plugin appears to be out of sync. run `vichat i` to install it again.")
		}

		var f = cmd.Flags()
		var temperature float32 = DefaultTemperature
		var maxTokens int = DefaultMaxTokens

		if m, err := f.GetInt("max_tokens"); err == nil {
			maxTokens = m
		}

		if t, err := f.GetFloat32("temperature"); err == nil {
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

		llm := vichat.New().WithTemperature(temperature).WithMaxTokens(maxTokens)
		prompts := CreatePrompts(lines)
		if len(prompts) == 0 {
			log.Fatalf("invalid input")
			return
		}

		if isSimpleChat {
			var promptStr []byte
			var err error
			prf, _ := f.GetString("system-prompt")
			if prf == "system.prompt" {
				promptStr = []byte(DefaultSystemPrompt)
			} else {
				promptStr, err = os.ReadFile(prf)
				if err != nil {
					prd := csv.NewReader(bytes.NewReader(awesomePrompts))
					embedPrompts, err := prd.ReadAll()
					if err == nil {
						index := make([]string, len(embedPrompts))
						for i := range embedPrompts {
							index[i] = strings.ToLower(embedPrompts[i][0])
						}

						matches := fuzzy.RankFind(prf, index)
						sort.Sort(matches)

						hit := matches[0].OriginalIndex
						promptStr = []byte(embedPrompts[hit][1])
					}
				}
			}

			prompts = append([]chat.PromptMessage{{
				Type:   chat.MessageTypeSystem,
				Prompt: prompt.New(string(promptStr)),
			}}, prompts...)
		}

		if ok, _ := f.GetBool("func"); ok {
			if err := llm.BindFunction(
				getRelativeTime,
				"getRelativeTime",
				`Use this function to find out what time is it using a relative duration of seconds. 
				Translate the time into a num of seconds before calling the function. 
				e.g. 1 hour ago = getRelativeTime(3600)
					 now = getRelativeTime(0)
				`,
			); err != nil {
				log.Fatalf("failed to bind function: %q", err.Error())
				return
			}
		}

		isTermOutput, _ := f.GetBool("term")
		messages := chat.New(prompts...)
		if isatty.IsTerminal(os.Stdout.Fd()) {
			if isTermOutput {
				resp, err := llm.Chat(context.TODO(), messages)
				if err != nil {
					log.Fatalf("failed to send chat: %q", err.Error())
					return
				}

				resp = string(markdown.Render(resp, 90, 4))

				fmt.Println()
				fmt.Println(resp)
				fmt.Println()
			} else if isSimpleChat {
				// open the full chat in vim
				dir, err := os.Getwd()
				if err != nil {
					dir = os.TempDir()
				}
				tmpf, err := os.CreateTemp(dir, "*.chat")
				if err != nil {
					log.Fatalf("failed to create temp file: %q", err)
				}

				fmt.Fprintf(tmpf, "# temperature=%.1f, max_tokens=%d\n\n", temperature, maxTokens)
				printPrompts(tmpf, prompts)
				tmpf.Close()

				// invoke vim using cmd and open tmpf
				cmd := exec.Command("vim", "-c", "redraw|Chat", tmpf.Name())
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout

				cmd.Run()
			}
		} else {
			// in vim mode, stream the output
			fmt.Printf("AI: ")

			err := llm.ChatStream(context.TODO(), func(s string) {
				fmt.Print(s)
			}, messages)

			if err != nil {
				log.Fatalf("failed to stream chat: %q", err.Error())
				return
			}

			fmt.Println("\n\nUSER: ")
		}
	},
}

func CreatePrompts(lines []string) []chat.PromptMessage {
	prompts := make([]chat.PromptMessage, 0)
	var messageType chat.MessageType = ""
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

	if messageType == "" {
		messageType = chat.MessageTypeUser
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

func printPrompts(w io.Writer, prompts []chat.PromptMessage) {
	for _, p := range prompts {
		prefix := ""
		switch p.Type {
		case chat.MessageTypeSystem:
			prefix = "SYSTEM"
		case chat.MessageTypeUser:
			prefix = "USER"
		case chat.MessageTypeAssistant:
			prefix = "AI"
		}

		fmt.Fprintf(w, "%s: %s\n\n", prefix, strings.Trim(p.Prompt.String(), "\r\n"))
	}
}
