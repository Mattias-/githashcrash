package startswithmatcher

import (
	"bytes"
	"encoding/hex"
)

type startswithmatcher struct {
	start []byte
}

// New constructs a new regexpmatcher
func New(start string) *startswithmatcher {
	dst := make([]byte, hex.DecodedLen(len([]byte(start))))
	_, err := hex.Decode(dst, []byte(start))
	if err != nil {
		panic(err)
	}

	return &startswithmatcher{
		dst,
	}
}

func (m *startswithmatcher) Match(hsum []byte) bool {
	return bytes.HasPrefix(hsum, m.start)
}
