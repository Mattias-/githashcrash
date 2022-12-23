package commitmsg

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"

	"github.com/Mattias-/githashcrash/pkg/worker"
)

type result struct {
	sha1   string
	object []byte
}

// Compile time check that interface is implemented
var _ worker.Result = result{}

func (r result) Sha1() string {
	return r.sha1
}

func (r result) Object() []byte {
	return r.object
}

func (r result) ShellRecreateCmd() string {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(r.object)
	w.Close()
	if err != nil {
		panic(err)
	}

	b64Content := base64.StdEncoding.EncodeToString(b.Bytes())
	return fmt.Sprintf("mkdir -p .git/objects/%s; echo '%s' | base64 -d > .git/objects/%s/%s; git reset %s\n", r.sha1[:2], b64Content, r.sha1[:2], r.sha1[2:], r.sha1)
}
