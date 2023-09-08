# vichat

A simple LLM chat cli that uses openai compatible apis.

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

