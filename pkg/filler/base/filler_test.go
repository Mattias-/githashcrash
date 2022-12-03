package base

import "testing"

func TestBase_bufferlen(t *testing.T) {
	f1 := New([]byte("abc"))
	buf1 := *(f1.OutputBuffer())

	f1.Fill(0)
	if len(buf1) != 12 {
		t.Errorf("Expected len %d, got %d", 4, len(buf1))
	}

	f1.Fill(1_000_000_000_000)
	if len(buf1) != 12 {
		t.Errorf("Expected len %d, got %d", 4, len(buf1))
	}
}

func TestBase_seeds(t *testing.T) {
	f1 := New([]byte{})
	buf1 := *(f1.OutputBuffer())
	f1.Fill(0)

	f2 := New([]byte(""))
	buf2 := *(f2.OutputBuffer())
	f2.Fill(0)
	if buf1[0] != buf2[0] {
		t.Errorf("Expected buffer values to be equal %d, %d", buf1[0], buf2[0])
	}
	if buf1[1] != buf2[1] {
		t.Errorf("Expected buffer values to be equal %d, %d", buf1[1], buf2[1])
	}
	if buf1[2] != buf2[2] {
		t.Errorf("Expected buffer values to be equal %d, %d", buf1[2], buf2[2])
	}
}
func TestBase_seeds2(t *testing.T) {
	f1 := New([]byte("aaaaaaaaaa"))
	buf2 := *(f1.OutputBuffer())
	f1.Fill(0)
	_ = buf2
}
