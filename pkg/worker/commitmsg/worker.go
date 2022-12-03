package commitmsg

import (
	"bytes"
	"crypto/sha1"
	"encoding"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Mattias-/githashcrash/pkg/worker"
)

type Matcher interface {
	Match([]byte) bool
}

type Filler interface {
	Fill(uint64)
	OutputBuffer() *[]byte
}

type Worker struct {
	count       uint64
	matcher     Matcher
	filler      Filler
	object      []byte
	placeholder []byte
}

func (w *Worker) Count() uint64 {
	return w.count
}

func NewWorker(matcher Matcher, filler Filler, obj, placeholder []byte) *Worker {
	return &Worker{
		count:       0,
		matcher:     matcher,
		filler:      filler,
		object:      obj,
		placeholder: placeholder,
	}
}

func (w *Worker) Work(rc chan worker.Result) {
	outputBuffer := w.filler.OutputBuffer()

	// Split on placeholder
	z := bytes.SplitN(w.object, w.placeholder, 2)
	before := z[0]
	var after []byte
	if len(z) == 2 {
		after = z[1]
	} else {
		// If no placeholder is found place it last.
		after = []byte("\n")
	}

	// The new object: Before placeholder, placeholder, After placeholder
	newObjLen := len(before) + len(*outputBuffer) + len(after)

	b := bytes.NewBufferString(fmt.Sprintf("commit %d\x00", newObjLen))
	b.Write(before)
	newObjectStart := b.Bytes()
	newObjectEnd := after

	// Freeze the hash of the object bytes before the placeholder
	first := sha1.New()
	first.Write(newObjectStart)
	marshaler, ok := first.(encoding.BinaryMarshaler)
	if !ok {
		log.Fatal("first does not implement encoding.BinaryMarshaler")
	}
	state, err := marshaler.MarshalBinary()
	if err != nil {
		log.Fatal("unable to marshal hash:", err)
	}

	for {
		w.count++
		w.filler.Fill(w.count)

		second := sha1.New()
		unmarshaler := second.(encoding.BinaryUnmarshaler)
		unmarshaler.UnmarshalBinary(state)
		second.Write(*outputBuffer)
		second.Write(newObjectEnd)
		hsum := second.Sum(nil)

		if w.matcher.Match(hsum) {
			b := bytes.NewBuffer(newObjectStart)
			b.Write(*outputBuffer)
			b.Write(newObjectEnd)

			rc <- result{
				sha1:   hex.EncodeToString(hsum),
				object: b.Bytes(),
			}
			return
		}
	}
}
