package accounts

import (
	"bytes"
	"go-solana-anchor/idl"
	"go-solana-anchor/internal/borsh"
	"testing"
)

func TestBuildInstructionData_NoArgs(t *testing.T) {
	disc := idl.Discriminator{11, 18, 104, 9, 104, 174, 59, 33}
	got := BuildInstructionData(disc, nil)
	if len(got) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(got))
	}
	for i := 0; i < 8; i++ {
		if got[i] != disc[i] {
			t.Errorf("byte %d: got %d, want %d", i, got[i], disc[i])
		}
	}
}

func TestBuildInstructionData_EmptyArgs(t *testing.T) {
	disc := idl.Discriminator{11, 18, 104, 9, 104, 174, 59, 33}
	got := BuildInstructionData(disc, []byte{})
	if len(got) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(got))
	}
}

func TestBuildInstructionData_WithArgs(t *testing.T) {
	disc := idl.Discriminator{11, 18, 104, 9, 104, 174, 59, 33}
	args := []byte{1, 2, 3, 4, 5}
	got := BuildInstructionData(disc, args)
	if len(got) != 8+5 {
		t.Fatalf("expected 13 bytes, got %d", len(got))
	}
	for i := 0; i < 8; i++ {
		if got[i] != disc[i] {
			t.Errorf("disc byte %d: got %d, want %d", i, got[i], disc[i])
		}
	}
	for i := 0; i < 5; i++ {
		if got[8+i] != args[i] {
			t.Errorf("args byte %d: got %d, want %d", i, got[8+i], args[i])
		}
	}
}

func TestEncodeInstructionArgs(t *testing.T) {
	args, err := EncodeInstructionArgs(func(e *borsh.Encoder) error {
		return e.EncodeU64(42)
	})
	if err != nil {
		t.Fatalf("EncodeInstructionArgs: %v", err)
	}
	if len(args) != 8 {
		t.Fatalf("expected 8 bytes for u64, got %d", len(args))
	}
	var buf bytes.Buffer
	e := borsh.NewEncoder(&buf)
	_ = e.EncodeU64(42)
	expected := buf.Bytes()
	for i := range expected {
		if args[i] != expected[i] {
			t.Errorf("byte %d: got %d, want %d", i, args[i], expected[i])
		}
	}
}
