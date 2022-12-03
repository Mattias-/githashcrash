package base

import (
	"encoding/base64"
	"encoding/binary"
)

// Base construct a short filling bytes by switching base
type base struct {
	seed         []byte
	seedLen      int
	inputBuffer  []byte
	outputBuffer []byte
}

var b64encoder = base64.RawStdEncoding

func New(seed []byte) *base {
	seedLen := 3
	if len(seed) > seedLen {
		seed = seed[:2]
	}

	rawCollisionLen := 9
	// collisionLen = 12
	collisionLen := b64encoder.EncodedLen(rawCollisionLen)

	inputBuffer := make([]byte, rawCollisionLen)
	outputBuffer := make([]byte, collisionLen)

	// inputBuffer always start with seed
	copy(inputBuffer, seed)
	return &base{
		seed,
		seedLen,
		inputBuffer,
		outputBuffer,
	}
}

func (b base) OutputBuffer() *[]byte {
	return &b.outputBuffer
}

// Fill output buffer with new value
func (b base) Fill(value uint64) {
	binary.PutUvarint(b.inputBuffer[b.seedLen:], value)
	// base64 encoding reduce byte size
	// from i to o
	b64encoder.Encode(b.outputBuffer, b.inputBuffer)
}
