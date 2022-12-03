package regexpmatcher

import (
	"encoding/hex"
	"testing"
)

func TestMatch(t *testing.T) {
	m := New("^30.*")
	s := ""
	if m.Match([]byte(s)) {
		t.Errorf("Expected hex(%s)=%s not to match %s", s, hex.EncodeToString([]byte(s)), m.String())
	}

	s = "0000"
	if !m.Match([]byte(s)) {
		t.Errorf("Expected hex(%s)=%s to match %s", s, hex.EncodeToString([]byte(s)), m.String())

	}
}
