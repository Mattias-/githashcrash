package main

import (
	"bytes"
	"crypto/sha1"
	"encoding"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
)

type Worker struct {
	i uint64
}

type Result struct {
	magic  string
	sha1   string
	object []byte
}

func getMatchFunc(targetHash *regexp.Regexp) func([]byte) bool {
	return targetHash.Match
}

func getFiller(seed []byte) (*[]byte, *[]byte, func(uint64)) {
	b64 := base64.RawStdEncoding
	// seedLen = 3 (as specified in advance)
	seedLen := len(seed)
	rawCollisionLen := 9
	// collisionLen = 12
	collisionLen := b64.EncodedLen(rawCollisionLen)

	inputBuffer := make([]byte, rawCollisionLen)
	outputBuffer := make([]byte, collisionLen)

	// inputBuffer always start with seed
	copy(inputBuffer, seed)

	// The filler function adds count and compacts and updates outputBuffer
	return &inputBuffer, &outputBuffer, func(count uint64) {
		// The last part of the input buffer contains the count
		binary.PutUvarint(inputBuffer[seedLen:], count)
		// base64 encoding reduce byte size
		// from i to o
		b64.Encode(outputBuffer, inputBuffer)
	}
}

func (w *Worker) work(targetHash *regexp.Regexp, obj []byte, seed []byte, placeholder []byte, result chan Result) {
	matcher := getMatchFunc(targetHash)
	b64 := base64.RawStdEncoding

	// Get filler function, it copies inputBuffer to outputBuffer after doing some modifications, like adding current count
	// Length is preserved
	inputBuffer, outputBuffer, filler := getFiller(seed)

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
		filler(w.i)

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

		if matcher(encodedBuffer) {
			newObject := append(newObjectStart, *outputBuffer...)
			newObject = append(newObject, newObjectEnd...)
			result <- Result{
				b64.EncodeToString(*inputBuffer),
				hex.EncodeToString(hsum),
				newObject}
			return
		}
	}
}
