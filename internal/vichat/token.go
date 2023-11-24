package vichat

import (
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"github.com/pkoukk/tiktoken-go-loader"
)

func Tokenize(text string, encoding string) ([]int, error) {
	tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())

	tkm, err := tiktoken.EncodingForModel(encoding)
	if err != nil {
		return nil, fmt.Errorf("getEncoding: %v", err)
	}

	// encode
	toks := tkm.Encode(text, nil, nil)
	return toks, nil
}

func Decode(tok []int, encoding string) ([]string, error) {
	tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())

	tkm, err := tiktoken.EncodingForModel(encoding)

	if err != nil {
		return nil, fmt.Errorf("getEncoding: %v", err)
	}

	res := make([]string, 0, len(tok))
	for i := range tok {
		res = append(res, tkm.Decode(tok[i:i+1]))
	}

	return res, nil
}
