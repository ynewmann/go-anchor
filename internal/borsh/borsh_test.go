package borsh

import (
	"bytes"
	"testing"
)

func TestU64RoundTrip(t *testing.T) {
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	if err := e.EncodeU64(12345678); err != nil {
		t.Fatal(err)
	}
	d := NewDecoder(&buf)
	v, err := d.DecodeU64()
	if err != nil {
		t.Fatal(err)
	}
	if v != 12345678 {
		t.Errorf("got %d, want 12345678", v)
	}
}

func TestStringRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	if err := e.EncodeString("hello"); err != nil {
		t.Fatal(err)
	}
	d := NewDecoder(&buf)
	s, err := d.DecodeString()
	if err != nil {
		t.Fatal(err)
	}
	if s != "hello" {
		t.Errorf("got %q, want hello", s)
	}
}
