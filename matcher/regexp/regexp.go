package regexpmatcher

import (
	"encoding/hex"
	"regexp"
)

type regexpmatcher struct {
	*regexp.Regexp
	encodedBuffer *[]byte
}

// New constructs a new regexpmatcher
func New(regexString string) *regexpmatcher {
	var targetHash = regexp.MustCompile(regexString)
	// Hex encoded SHA1 is 40 bytes
	encodedBuffer := make([]byte, 40)
	return &regexpmatcher{
		targetHash,
		&encodedBuffer,
	}
}

func (m *regexpmatcher) Match(hsum []byte) bool {
	hex.Encode(*m.encodedBuffer, hsum)
	return m.Regexp.Match(*m.encodedBuffer)
}
