// Package client provides the Anchor program client and instruction builder.
package client

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"go-solana-sdk/idl"
)

// Program is an Anchor program client.
type Program struct {
	IDL       *idl.IDL
	ProgramID solana.PublicKey
	RPC       *rpc.Client
}

// NewProgram creates a new Program client.
func NewProgram(idl *idl.IDL, programID solana.PublicKey, rpcClient *rpc.Client) *Program {
	return &Program{
		IDL:       idl,
		ProgramID: programID,
		RPC:       rpcClient,
	}
}
