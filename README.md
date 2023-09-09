# vichat

A simple LLM chat cli that uses openai compatible apis.

## Why

I find it easier to use text files to test different system prompts on the same user input than a Web UI.

## Installation

```
go install github.com/JackKCWong/vichat@latest
```

## Usage

```
cat <<_EOF | vichat chat
SYSTEM: You are a professional joke writer

USER: tell me a joke about goose

AI: Why did the goose go to the doctor? Because he was feeling a little down!

USER: Why did it feel down?
_EOF
```

Use in Vim:

```
vichat install-vim

cat <<_EOF > test.chat
SYSTEM: You are a professional joke writer

USER: tell me a joke about goose

AI: Why did the goose go to the doctor? Because he was feeling a little down!

USER: Why did it feel down?
_EOF

vim test.chat
```

When in Vim (normal mode or insert mode):

* `ctrl+s` to send the full chat, the result will be appended to the end of the current buffer, just like a usual chat experience.

* `ctrl+t` to send the full chat, but the result will be put to a new buffer, useful when you want to ask many simple questions that doesn't require context. Saves token usages.
