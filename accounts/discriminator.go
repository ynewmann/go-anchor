// Package accounts provides account types, encoding, and discriminators for Anchor.
package accounts

import (
	"crypto/sha256"
	"go-solana-sdk/idl"
)

// AccountDiscriminator returns the 8-byte discriminator for an account type.
// Matches Anchor: sha256("account:<TypeName>")[0:8]
func AccountDiscriminator(typeName string) idl.Discriminator {
	return hashPrefix("account:" + typeName)
}

// InstructionDiscriminator returns the 8-byte discriminator for an instruction.
// Matches Anchor: sha256("global:<InstructionName>")[0:8]
func InstructionDiscriminator(instructionName string) idl.Discriminator {
	return hashPrefix("global:" + instructionName)
}

// EventDiscriminator returns the 8-byte discriminator for an event.
// Matches Anchor: sha256("event:<EventName>")[0:8]
func EventDiscriminator(eventName string) idl.Discriminator {
	return hashPrefix("event:" + eventName)
}

func hashPrefix(s string) idl.Discriminator {
	h := sha256.Sum256([]byte(s))
	var d idl.Discriminator
	copy(d[:], h[:8])
	return d
}

// ValidateAccountDiscriminator checks that data starts with the expected discriminator.
func ValidateAccountDiscriminator(data []byte, expected idl.Discriminator) bool {
	if len(data) < 8 {
		return false
	}
	for i := 0; i < 8; i++ {
		if data[i] != expected[i] {
			return false
		}
	}
	return true
}
