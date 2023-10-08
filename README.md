# vichat

A simple LLM chat cli (with Vim).

![demo](https://github.com/JackKCWong/vichat/blob/main/vichat.gif?raw=true)

## Why

I find it easier to use text files to test different system prompts on the same user input than a Web UI.

## Installation

```bash
go install github.com/JackKCWong/vichat@latest

# install vim plugin
vichat i 
```

## Usage

```bash
export OPENAI_API_KEY=<your LLM api key>
export OPENAI_API_BASE=<your LLM api base>
```

Chat directly on cli and open the respone in Vim:

```bash
vichat chat [options] tell me a joke about Go
# or
vichat chat -t 0.1 -m 100
tell me a joke about Golang
^D
# or output to terminal directly
vichat chat -o tell me a joke about Vim

# or just simply
vichat tell me a joke about Vim
```

Chat with history:
```bash
cat <<_EOF | vichat chat
SYSTEM: You are a professional joke writer

USER: tell me a joke about goose

AI: Why did the goose go to the doctor? Because he was feeling a little down!

USER: Why did it feel down?
_EOF
```

Use in Vim:

```bash
cat <<_EOF > test.chat
SYSTEM: You are a professional joke writer

USER: tell me a joke about goose
_EOF

vim test.chat
```

When in Vim:

* `ctrl+s` (n or i mode) to send the full chat, the result will be appended to the end of the current buffer, just like a usual chat experience.

* `ctrl+t` (n or i mode) to send the full chat, but the result will be put to a new buffer, useful when you want to ask many simple questions that doesn't require context. Saves token usages.

* `ctrl+a` (n mode) to jump to the end and start asking a new question.

* `ctrl+t` (n mode) to estimate the number of tokens using titoken gpt-like encoding.

* `ctrl+k` (n or v mode) to count the number of tokens.

* `ctrl+c` (n or v mode) to chunk the text using RecursiveTextSplitter.

* `ctrl+n` (n mode) to start a new chat with the same system prompt.

* `q` (n mode) close current chat.


Vim tips:

* put this line in your `~/.vimrc` to enable code block highlight in markdown

```vim
syntax on
let g:markdown_fenced_languages = ['html', 'js=javascript', 'rust', 'go', 'java']
```

