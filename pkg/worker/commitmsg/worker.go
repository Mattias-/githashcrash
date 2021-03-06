package commitmsg

import (
	"bytes"
	"crypto/sha1"
	"encoding"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Mattias-/githashcrash/pkg/filler"
	"github.com/Mattias-/githashcrash/pkg/matcher"
	"github.com/Mattias-/githashcrash/pkg/worker"
)

func split2(h, needle []byte) ([]byte, []byte) {
	// Split on placeholder
	z := bytes.SplitN(h, needle, 2)
	var before []byte
	before = z[0]
	var after []byte
	if len(z) == 2 {
		after = z[1]
	} else {
		// If no placeholder is found place it last.
		after = []byte("\n")
	}
	return before, after
}

type worker2 struct {
	i uint64
}

func (w *worker2) Count() uint64 {
	return w.i
}

func NewW() worker.Worker {
	return &worker2{0}
}

func (w *worker2) Work(m matcher.Matcher, f filler.Filler, obj []byte, placeholder []byte, result chan worker.Result) {
	outputBuffer := f.OutputBuffer()

	// Split on placeholder
	before, after := split2(obj, placeholder)

	newObjLen := len(before) + len(*outputBuffer) + len(after)
	newObjectStart := append([]byte(fmt.Sprintf("commit %d\x00", newObjLen)), before...)
	newObjectEnd := after

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

	for ; ; w.i++ {
		f.Fill(w.i)

		second := sha1.New()
		unmarshaler, ok := second.(encoding.BinaryUnmarshaler)
		if !ok {
			log.Fatal("second does not implement encoding.BinaryUnmarshaler")
		}
		if err := unmarshaler.UnmarshalBinary(state); err != nil {
			log.Fatal("unable to unmarshal hash:", err)
		}
		second.Write(*outputBuffer)
		second.Write(newObjectEnd)
		hsum := second.Sum(nil)

		if m.Match(hsum) {
			newObject := append(newObjectStart, *outputBuffer...)
			newObject = append(newObject, newObjectEnd...)
			result <- worker.Result{
				hex.EncodeToString(hsum),
				newObject}
			return
		}
	}
}
