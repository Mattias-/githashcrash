package main

import (
	"encoding/base64"
	"log"
	"regexp"
	"testing"
)

func TestRun(b *testing.T) {
	//obj, _ := exec.Command("git", "cat-file", "-p", "HEAD").Output()
	obj := []byte(`tree 5fe7b5921cf9f617615100dae6cc20747e6140e6
parent 0a681dcc33851637c3f5fbce17de93547c7d180b
author Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100
committer Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100

Better counting
`)
	res := run("^00000.*", obj, "11", 1)
	log.Println(res.sha1)
	log.Println(base64.StdEncoding.EncodeToString(res.object))
}

func BenchmarkWorker(b *testing.B) {
	hashRe := "^00000.*"
	var targetHash = regexp.MustCompile(hashRe)
	obj := []byte(`tree 5fe7b5921cf9f617615100dae6cc20747e6140e6
parent 0a681dcc33851637c3f5fbce17de93547c7d180b
author Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100
committer Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100

Better counting
`)

	for i := 0; i < b.N; i++ {
		results := make(chan Result)
		w := &Worker{0}
		go w.work(targetHash, obj, []byte("000"), results)
		<-results
	}
}
