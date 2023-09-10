# vichat

A simple LLM chat cli that uses openai compatible apis.

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

When in Vim (normal mode or insert mode):

* `ctrl+s` to send the full chat, the result will be appended to the end of the current buffer, just like a usual chat experience.

* `ctrl+t` to send the full chat, but the result will be put to a new buffer, useful when you want to ask many simple questions that doesn't require context. Saves token usages.


Vim tips:

* put this line in your `~/.vimrc` to enable code block highlight in markdown

```vim
let g:markdown_fenced_languages = ['html', 'js=javascript', 'rust', 'go', 'java']
```
