package main

import (
	"bytes"
	"crypto/sha1"
	"encoding"
	"encoding/hex"
	"fmt"
	"githashcrash/filler/base"
	"githashcrash/matcher/regexp"
	"log"
	"regexp"
)

func (w *Worker) work(targetHash *regexp.Regexp, obj []byte, seed []byte, placeholder []byte, result chan Result) {
	matcher := regexpmatcher.New(targetHash.String())
	outputBuffer, filler := basefiller.New(seed)

	// Split on placeholder
	z := bytes.SplitN(obj, placeholder, 2)
	var before []byte
	before = z[0]
	var after []byte
	if len(z) == 2 {
		after = z[1]
	} else {
		after = []byte("\n")
	}

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

	// Hex encoded SHA1 is 40 bytes
	encodedBuffer := make([]byte, 40)

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
		hex.Encode(encodedBuffer, hsum)

		if matcher.Match(encodedBuffer) {
			newObject := append(newObjectStart, *outputBuffer...)
			newObject = append(newObject, newObjectEnd...)
			result <- Result{
				hex.EncodeToString(hsum),
				newObject}
			return
		}
	}
}
