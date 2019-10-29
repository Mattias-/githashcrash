package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"githashcrash/worker"
	"log"
)

func printRecreate(result worker.Result) {
	log.Println("Create with:")

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(result.Object)
	w.Close()

	b64Content := base64.StdEncoding.EncodeToString(b.Bytes())
	fmt.Printf("echo '%s' | base64 -d >.git/objects/%s/%s; git reset %s\n", b64Content, result.Sha1[:2], result.Sha1[2:], result.Sha1)
}
