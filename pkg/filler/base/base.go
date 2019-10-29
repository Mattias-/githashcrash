package base

import (
	"encoding/base64"
	"encoding/binary"
)

// Base construct a short filling bytes by switching base
type base struct {
	seed         []byte
	seedLen      int
	encoding     base64.Encoding
	inputBuffer  []byte
	outputBuffer []byte
}

func New(seed []byte) *base {
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
	return &base{
		seed,
		seedLen,
		*b64,
		inputBuffer,
		outputBuffer,
	}
}

func (b base) OutputBuffer() *[]byte {
	return &b.outputBuffer
}

func (b base) Fill(count uint64) {
	binary.PutUvarint(b.inputBuffer[b.seedLen:], count)
	// base64 encoding reduce byte size
	// from i to o
	b.encoding.Encode(b.outputBuffer, b.inputBuffer)
}
