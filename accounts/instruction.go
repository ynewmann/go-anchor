package accounts

import (
	"bytes"
	"go-solana-sdk/idl"
	"go-solana-sdk/internal/borsh"
)

// BuildInstructionData builds instruction data: discriminator + Borsh-encoded args.
// For instructions with no args, returns just the discriminator.
func BuildInstructionData(disc idl.Discriminator, args []byte) []byte {
	if len(args) == 0 {
		return disc[:]
	}
	out := make([]byte, 0, 8+len(args))
	out = append(out, disc[:]...)
	out = append(out, args...)
	return out
}

// EncodeInstructionArgs encodes instruction arguments to Borsh.
// For simple cases (e.g. increment with no args), pass nil.
// For struct args, the caller should use borsh.Encoder to produce args.
func EncodeInstructionArgs(fn func(*borsh.Encoder) error) ([]byte, error) {
	var buf bytes.Buffer
	e := borsh.NewEncoder(&buf)
	if err := fn(e); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
