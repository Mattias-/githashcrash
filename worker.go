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
	"sync/atomic"
)

func worker(targetHash *regexp.Regexp, obj []byte, seed []byte, result chan string, testOps *uint64) {
	i := uint64(0)
	buf := make([]byte, binary.MaxVarintLen64-3)

	rawCollisionBuffer := make([]byte, binary.MaxVarintLen64)
	copy(rawCollisionBuffer, seed)

	collisionLength := base64.RawStdEncoding.EncodedLen(len(seed) + len(buf))
	collisionByteBuffer := make([]byte, collisionLength)

	newObjLen := len(obj) + collisionLength + 1
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

	for {
		binary.PutUvarint(buf, i)
		copy(rawCollisionBuffer[3:], buf)
		base64.RawStdEncoding.Encode(collisionByteBuffer, rawCollisionBuffer)

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
			result <- base64.RawStdEncoding.EncodeToString(rawCollisionBuffer)
			return
		}

		if i%100000 == 0 {
			atomic.AddUint64(testOps, 100000)
		}
		i++
	}
}
