package pda

import (
	"testing"

	"github.com/gagliardetto/solana-go"
)

func TestFindProgramAddress(t *testing.T) {
	programID, _ := solana.PublicKeyFromBase58("6khKp4BeJpCjBY1Eh39ybiqbfRnrn2UzWeUARjQLXYRC")
	seeds := [][]byte{[]byte("counter")}
	addr, bump, err := FindProgramAddress(seeds, programID)
	if err != nil {
		t.Fatal(err)
	}
	if bump == 0 {
		t.Error("expected non-zero bump")
	}
	if addr.Equals(solana.PublicKey{}) {
		t.Error("expected non-zero address")
	}
}
