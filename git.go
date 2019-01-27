package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
)

func printRecreate(result Result) {
	fmt.Println("Create with:")

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(result.object)
	w.Close()
	b64Content := base64.StdEncoding.EncodeToString(b.Bytes())
	fmt.Printf("echo '%s' | base64 -d >.git/objects/%s/%s; git reset %s\n", b64Content, result.sha1[:2], result.sha1[2:], result.sha1)
}
