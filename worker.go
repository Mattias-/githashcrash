package main

import (
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

func (w *Worker) work(targetHash *regexp.Regexp, obj []byte, seed []byte, result chan Result) {
	b64 := base64.RawStdEncoding
	// 3
	seedLen := len(seed)
	rawCollisionLen := 9
	// 12
	collisionLen := b64.EncodedLen(rawCollisionLen)

	rawCollisionBuffer := make([]byte, rawCollisionLen)
	copy(rawCollisionBuffer, seed)

	collisionByteBuffer := make([]byte, collisionLen)

	newObjLen := len(obj) + collisionLen + 1
	newObjectStart := append([]byte(fmt.Sprintf("commit %d\x00", newObjLen)), obj...)
	newObjectEnd := []byte("\n")

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
		binary.PutUvarint(rawCollisionBuffer[seedLen:], w.i)
		b64.Encode(collisionByteBuffer, rawCollisionBuffer)

		second := sha1.New()
		unmarshaler, ok := second.(encoding.BinaryUnmarshaler)
		if !ok {
			log.Fatal("second does not implement encoding.BinaryUnmarshaler")
		}
		if err := unmarshaler.UnmarshalBinary(state); err != nil {
			log.Fatal("unable to unmarshal hash:", err)
		}
		second.Write(collisionByteBuffer)
		second.Write(newObjectEnd)

		hsum := second.Sum(nil)
		hex.Encode(encodedBuffer, hsum)

		if targetHash.Match(encodedBuffer) {
			newObject := append(newObjectStart, collisionByteBuffer...)
			newObject = append(newObject, newObjectEnd...)
			result <- Result{
				b64.EncodeToString(rawCollisionBuffer),
				hex.EncodeToString(hsum),
				newObject}
			return
		}
	}
}
