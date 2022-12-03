package prefixmatcher

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/Mattias-/githashcrash/pkg/matcher"
)

type prefixmatcher struct {
	prefix []byte
}

func New(start string) matcher.Matcher {
	dst := make([]byte, hex.DecodedLen(len([]byte(start))))
	_, err := hex.Decode(dst, []byte(start))
	if err != nil {
		panic(err)
	}

	return &prefixmatcher{
		prefix: dst,
	}
}

func (m prefixmatcher) String() string {
	return fmt.Sprintf("prefix(%s)", m.prefix)
}

func (m *prefixmatcher) Match(hsum []byte) bool {
	return bytes.HasPrefix(hsum, m.prefix)
}
