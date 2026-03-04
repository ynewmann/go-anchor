package idlcmd

import (
	"fmt"
	"os"

	"go-solana-anchor/idl"
)

// ValidateIDL reads and validates an IDL JSON file.
func ValidateIDL(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	parsed, err := idl.Parse(data)
	if err != nil {
		return fmt.Errorf("parse IDL: %w", err)
	}
	// Basic validation
	if parsed.Address == "" {
		return fmt.Errorf("IDL missing address")
	}
	if len(parsed.Instructions) == 0 {
		return fmt.Errorf("IDL has no instructions")
	}
	return nil
}
