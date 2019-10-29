package worker

import (
	"github.com/Mattias-/githashcrash/pkg/filler"
	"github.com/Mattias-/githashcrash/pkg/matcher"
)

type Worker interface {
	Count() uint64
	Work(matcher.Matcher, filler.Filler, []byte, []byte, chan Result)
}
