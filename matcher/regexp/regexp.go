package regexpmatcher

import (
	"regexp"
)

type regexpmatcher struct {
	*regexp.Regexp
}

// New constructs a new regexpmatcher
func New(regexString string) *regexpmatcher {
	var targetHash = regexp.MustCompile(regexString)
	return &regexpmatcher{targetHash}
}

func (m *regexpmatcher) Match(s []byte) bool {
	return m.Regexp.Match(s)
}
