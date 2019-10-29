package main

import (
	"bytes"
	filler "githashcrash/filler/base"
	"githashcrash/matcher/regexp"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	obj := []byte(`tree 5fe7b5921cf9f617615100dae6cc20747e6140e6
parent 0a681dcc33851637c3f5fbce17de93547c7d180b
author Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100
committer Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100

Better counting
`)
	placeholder := []byte("REPLACEME")
	res := run("^00000.*", obj, []byte("11"), 1, placeholder)
	if !strings.HasPrefix(res.sha1, "00000") {
		t.Fail()
	}
	if bytes.Equal(obj, res.object) {
		t.Fail()
	}
}

func TestReplace(t *testing.T) {
	obj := []byte(`tree 5fe7b5921cf9f617615100dae6cc20747e6140e6
parent 0a681dcc33851637c3f5fbce17de93547c7d180b
author Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100
committer Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100

Better counting

Hello world: REPLACEME abc
`)
	placeholder := []byte("REPLACEME")
	res := run("^00000.*", obj, []byte("11"), 1, placeholder)
	if !strings.HasPrefix(res.sha1, "00000") {
		t.Fail()
	}
	if bytes.Equal(obj, res.object) {
		t.Fail()
	}
	if bytes.Contains(res.object, placeholder) {
		t.Fail()
	}
}

func BenchmarkWorker(b *testing.B) {
	hashRe := "^00000.*"
	matcher := regexpmatcher.New(hashRe)
	obj := []byte(`tree 5fe7b5921cf9f617615100dae6cc20747e6140e6
parent 0a681dcc33851637c3f5fbce17de93547c7d180b
author Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100
committer Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100

Better counting

Hello world: REPLACEME abc
The end
`)
	placeholder := []byte("REPLACEME")
	for i := 0; i < b.N; i++ {
		results := make(chan Result)
		w := NewW()
		filler := filler.New([]byte("000"))
		go w.Work(matcher, filler, obj, placeholder, results)
		<-results
	}
}
