package commitmsg

import (
	"bytes"
	"strings"
	"testing"

	filler "github.com/Mattias-/githashcrash/pkg/filler/base"
	matcher "github.com/Mattias-/githashcrash/pkg/matcher/startswith"
	"github.com/Mattias-/githashcrash/pkg/worker"
)

func TestWorker(t *testing.T) {
	hashRe := "0000"
	matcher := matcher.New(hashRe)
	obj := []byte(`tree 5fe7b5921cf9f617615100dae6cc20747e6140e6
parent 0a681dcc33851637c3f5fbce17de93547c7d180b
author Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100
committer Mattias Appelgren <mattias@ppelgren.se> 1548587085 +0100

Better counting

Hello world: REPLACEME abc
The end
`)
	placeholder := []byte("REPLACEME")
	results := make(chan worker.Result)
	w := NewW()
	filler := filler.New([]byte("000"))
	go w.Work(matcher, filler, obj, placeholder, results)
	var r = <-results

	if !strings.HasPrefix(r.Sha1, "0000") {
		t.Fail()
	}
	if bytes.Equal(obj, r.Object) {
		t.Fail()
	}
	if bytes.Contains(r.Object, placeholder) {
		t.Fail()
	}

}

func BenchmarkWorker(b *testing.B) {
	hashRe := "0000"
	matcher := matcher.New(hashRe)
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
		results := make(chan worker.Result)
		w := NewW()
		filler := filler.New([]byte("000"))
		go w.Work(matcher, filler, obj, placeholder, results)
		var r = <-results

		if !strings.HasPrefix(r.Sha1, "0000") {
			b.Fail()
		}
		if bytes.Equal(obj, r.Object) {
			b.Fail()
		}
		if bytes.Contains(r.Object, placeholder) {
			b.Fail()
		}

	}
}
