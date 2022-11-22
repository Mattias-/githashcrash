package commitmsg

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/Mattias-/githashcrash/pkg/worker"
)

func PrintRecreate(result worker.Result) {
	log.Println("Create with:")

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(result.Object)
	w.Close()

	b64Content := base64.StdEncoding.EncodeToString(b.Bytes())
	fmt.Printf("mkdir -p .git/objects/%s; echo '%s' | base64 -d > .git/objects/%s/%s; git reset %s\n", result.Sha1[:2], b64Content, result.Sha1[:2], result.Sha1[2:], result.Sha1)
}
