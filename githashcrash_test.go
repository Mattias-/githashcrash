package main

import (
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
	run("^00000.*", obj, "11", 8)
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
		results := make(chan string)
		w := &Worker{0}
		go w.worker(targetHash, obj, []byte("000"), results)
		<-results
	}
}
