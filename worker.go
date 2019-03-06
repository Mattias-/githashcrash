package main

import (
	"bytes"
	"crypto/sha1"
	"encoding"
	"encoding/hex"
	"fmt"
	"githashcrash/filler/base"
	"log"
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

func (w *Worker) work(matcher Matcher, obj []byte, seed []byte, placeholder []byte, result chan Result) {
	outputBuffer, filler := basefiller.New(seed)

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
		filler.Fill(w.i)

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

		if matcher.Match(hsum) {
			newObject := append(newObjectStart, *outputBuffer...)
			newObject = append(newObject, newObjectEnd...)
			result <- Result{
				hex.EncodeToString(hsum),
				newObject}
			return
		}
	}
}
