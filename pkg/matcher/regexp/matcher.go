package regexpmatcher

import (
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/Mattias-/githashcrash/pkg/matcher"
)

type regexpmatcher struct {
	regexp        *regexp.Regexp
	encodedBuffer *[]byte
	exp           string
}

func New(exp string) matcher.Matcher {
	var targetHash = regexp.MustCompile(exp)
	// Hex encoded SHA1 is 40 bytes
	encodedBuffer := make([]byte, 40)
	return &regexpmatcher{
		targetHash,
		&encodedBuffer,
		exp,
	}
}

func (m *regexpmatcher) String() string {
	return fmt.Sprintf("regexp(%s)", m.exp)
}

func (m *regexpmatcher) Match(hsum []byte) bool {
	hex.Encode(*m.encodedBuffer, hsum)
	return m.regexp.Match(*m.encodedBuffer)
}
