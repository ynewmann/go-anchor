package client

import (
	"bytes"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"go-solana-sdk/accounts"
	"go-solana-sdk/idl"
	"go-solana-sdk/internal/borsh"
)

// ArgsEncoder is a function that encodes instruction arguments to Borsh.
type ArgsEncoder func(*borsh.Encoder) error

// BuildInstruction builds a solana.Instruction for the given instruction name.
// accounts maps IDL account names (e.g. "counter", "authority") to PublicKeys.
// args can be nil (no args), []byte (raw encoded), or ArgsEncoder (encode via callback).
func (p *Program) BuildInstruction(name string, args interface{}, accountsMap map[string]solana.PublicKey) (solana.Instruction, error) {
	var inst *idl.Instruction
	for i := range p.IDL.Instructions {
		if p.IDL.Instructions[i].Name == name {
			inst = &p.IDL.Instructions[i]
			break
		}
	}
	if inst == nil {
		return nil, fmt.Errorf("instruction %q not found in IDL", name)
	}

	argsBytes, err := encodeArgs(args)
	if err != nil {
		return nil, fmt.Errorf("encode args: %w", err)
	}

	data := accounts.BuildInstructionData(inst.Discriminator, argsBytes)

	metas := make(solana.AccountMetaSlice, 0, len(inst.Accounts))
	for _, acc := range inst.Accounts {
		pubkey, ok := accountsMap[acc.Name]
		if !ok && !acc.Optional {
			return nil, fmt.Errorf("missing required account %q", acc.Name)
		}
		if !ok {
			continue // skip optional accounts not provided
		}
		meta := solana.Meta(pubkey)
		if acc.Writable {
			meta = meta.WRITE()
		}
		if acc.Signer {
			meta = meta.SIGNER()
		}
		metas = append(metas, meta)
	}

	return solana.NewInstruction(p.ProgramID, metas, data), nil
}

func encodeArgs(args interface{}) ([]byte, error) {
	if args == nil {
		return nil, nil
	}
	switch v := args.(type) {
	case []byte:
		return v, nil
	case ArgsEncoder:
		var buf bytes.Buffer
		e := borsh.NewEncoder(&buf)
		if err := v(e); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	default:
		return nil, fmt.Errorf("args must be nil, []byte, or ArgsEncoder, got %T", args)
	}
}
