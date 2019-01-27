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

func (w *Worker) worker(targetHash *regexp.Regexp, obj []byte, seed []byte, result chan string) {
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

		hex.Encode(encodedBuffer, second.Sum(nil))

		if targetHash.Match(encodedBuffer) {
			result <- b64.EncodeToString(rawCollisionBuffer)
			return
		}
	}
}
