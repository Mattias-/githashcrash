package filler

type Filler interface {
	Fill(uint64)
	OutputBuffer() *[]byte
}
