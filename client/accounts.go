package client

import (
	"bytes"
	"fmt"

	"go-solana-sdk/accounts"
	"go-solana-sdk/idl"
	"go-solana-sdk/internal/borsh"
)

// AccountDecoder is a function that decodes account data from a Borsh decoder.
type AccountDecoder func(*borsh.Decoder) (interface{}, error)

// DecodeAccount validates the account discriminator and decodes account data.
// name is the IDL account type name (e.g. "Counter").
// decoder is called with a decoder positioned after the 8-byte discriminator to decode the payload.
// Returns the decoded value or an error if discriminator is invalid.
func (p *Program) DecodeAccount(name string, data []byte, decoder AccountDecoder) (interface{}, error) {
	var acc *idl.Account
	for i := range p.IDL.Accounts {
		if p.IDL.Accounts[i].Name == name {
			acc = &p.IDL.Accounts[i]
			break
		}
	}
	if acc == nil {
		return nil, fmt.Errorf("account type %q not found in IDL", name)
	}

	if !accounts.ValidateAccountDiscriminator(data, acc.Discriminator) {
		return nil, fmt.Errorf("account discriminator mismatch for %q", name)
	}

	d := borsh.NewDecoder(bytes.NewReader(data[8:]))
	return decoder(d)
}
