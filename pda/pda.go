// Package pda provides Program Derived Address (PDA) helpers for Solana.
package pda

import (
	"github.com/gagliardetto/solana-go"
)

// FindProgramAddress implements the standard Solana PDA algorithm.
// It tries bumps from 255 down to 1 until create_program_address succeeds.
// Returns the derived address and the canonical bump seed.
func FindProgramAddress(seeds [][]byte, programID solana.PublicKey) (solana.PublicKey, uint8, error) {
	return solana.FindProgramAddress(seeds, programID)
}
