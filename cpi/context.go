// Package cpi provides CPI (cross-program invocation) context and helpers.
package cpi

import (
	"github.com/gagliardetto/solana-go"
)

// SignerSeed is a seed used for PDA signing in invoke_signed.
// When a program performs CPI with invoke_signed, it provides seeds so the runtime
// can let PDAs sign. The client uses these when building instructions that
// include PDA signer accounts.
type SignerSeed struct {
	// Seeds for the PDA (excluding the bump).
	Seeds [][]byte
	// ProgramID that owns the PDA.
	ProgramID solana.PublicKey
}

// CpiContext holds context for building a CPI instruction.
type CpiContext struct {
	// ProgramID is the program being invoked.
	ProgramID solana.PublicKey
	// Accounts maps account names to PublicKeys for the CPI.
	Accounts map[string]solana.PublicKey
	// SignerSeeds are PDAs that will sign when the calling program uses invoke_signed.
	// The client includes these PDAs in the account list; the program provides
	// the seeds to the runtime.
	SignerSeeds []SignerSeed
}

// NewCpiContext creates a CpiContext for invoking the given program.
func NewCpiContext(programID solana.PublicKey, accounts map[string]solana.PublicKey) *CpiContext {
	return &CpiContext{
		ProgramID:   programID,
		Accounts:    accounts,
		SignerSeeds: nil,
	}
}

// WithSignerSeeds adds signer seeds for invoke_signed.
func (c *CpiContext) WithSignerSeeds(seeds []SignerSeed) *CpiContext {
	c.SignerSeeds = seeds
	return c
}

// Invoke builds a CPI instruction. The returned instruction targets ProgramID
// with the given accounts and data. Callers add this to a transaction.
func Invoke(programID solana.PublicKey, accounts solana.AccountMetaSlice, data []byte) solana.Instruction {
	return solana.NewInstruction(programID, accounts, data)
}

// InvokeSigned builds a CPI instruction that will be invoked with invoke_signed
// by the calling program. The accounts must include the PDA signers; their
// corresponding seeds are used on-chain by the calling program.
func InvokeSigned(programID solana.PublicKey, accounts solana.AccountMetaSlice, data []byte, _ []SignerSeed) solana.Instruction {
	// Signer seeds are used by the program at runtime, not by the client.
	// The client just builds the instruction with the correct accounts.
	return solana.NewInstruction(programID, accounts, data)
}
