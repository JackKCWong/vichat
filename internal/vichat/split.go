package vichat

import (
	"github.com/henomis/lingoose/textsplitter"
)

func RecursiveTextSplit(text string, chunkSize, overlap int) []string {
	splitter := textsplitter.
		NewRecursiveCharacterTextSplitter(chunkSize, overlap)

	return splitter.SplitText(text)
}
