// Package spl provides SPL Token and Associated Token Account helpers.
package spl

import (
	"github.com/gagliardetto/solana-go"
)

// TokenProgramID is the SPL Token program.
var TokenProgramID = solana.TokenProgramID

// AssociatedTokenProgramID is the SPL Associated Token Account program.
var AssociatedTokenProgramID = solana.SPLAssociatedTokenAccountProgramID

// FindAssociatedTokenAddress returns the ATA for the given wallet and mint.
func FindAssociatedTokenAddress(wallet, mint solana.PublicKey) (solana.PublicKey, uint8, error) {
	return solana.FindAssociatedTokenAddress(wallet, mint)
}

// CreateAssociatedTokenAccountInstruction creates an instruction to create
// an associated token account. Payer funds the account creation.
// Wallet is the ATA owner; mint is the token mint.
func CreateAssociatedTokenAccountInstruction(
	payer, wallet, mint solana.PublicKey,
) solana.Instruction {
	ata, _, _ := solana.FindAssociatedTokenAddress(wallet, mint)
	// Create instruction has empty data.
	accounts := solana.AccountMetaSlice{
		solana.Meta(payer).WRITE().SIGNER(),
		solana.Meta(ata).WRITE(),
		solana.Meta(wallet),
		solana.Meta(mint),
		solana.Meta(solana.SystemProgramID),
		solana.Meta(solana.TokenProgramID),
	}
	return solana.NewInstruction(solana.SPLAssociatedTokenAccountProgramID, accounts, nil)
}
