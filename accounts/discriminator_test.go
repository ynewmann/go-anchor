package accounts

import (
	"encoding/hex"
	"go-solana-anchor/idl"
	"testing"
)

func TestAccountDiscriminator(t *testing.T) {
	// Anchor "account:Counter" - known from counter example
	d := AccountDiscriminator("Counter")
	if len(d) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(d))
	}
	// Verify it matches the IDL counter example discriminator
	expected := idl.Discriminator{255, 176, 4, 245, 188, 253, 124, 25}
	for i := 0; i < 8; i++ {
		if d[i] != expected[i] {
			t.Errorf("byte %d: got %d, want %d (full: %s)", i, d[i], expected[i], hex.EncodeToString(d[:]))
		}
	}
}

func TestInstructionDiscriminator(t *testing.T) {
	d := InstructionDiscriminator("increment")
	if len(d) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(d))
	}
	expected := idl.Discriminator{11, 18, 104, 9, 104, 174, 59, 33}
	for i := 0; i < 8; i++ {
		if d[i] != expected[i] {
			t.Errorf("byte %d: got %d, want %d", i, d[i], expected[i])
		}
	}
}

func TestEventDiscriminator(t *testing.T) {
	d := EventDiscriminator("MyEvent")
	if len(d) != 8 {
		t.Fatalf("expected 8 bytes, got %d", len(d))
	}
}

func TestValidateAccountDiscriminator(t *testing.T) {
	disc := AccountDiscriminator("Counter")
	valid := append(disc[:], []byte("more data")...)
	if !ValidateAccountDiscriminator(valid, disc) {
		t.Error("expected valid")
	}
	if ValidateAccountDiscriminator([]byte("short"), disc) {
		t.Error("expected invalid for short data")
	}
	wrong := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
	if ValidateAccountDiscriminator(wrong, disc) {
		t.Error("expected invalid for wrong discriminator")
	}
}
