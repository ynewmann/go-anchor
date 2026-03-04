// Package errors provides error code mapping for Anchor programs.
package errors

import (
	"encoding/binary"
	"fmt"

	"go-solana-anchor/idl"
)

// ProgramError wraps an Anchor program error with its IDL metadata.
type ProgramError struct {
	Code    uint64
	Name    string
	Message string
}

func (e *ProgramError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("program error %d (%s): %s", e.Code, e.Name, e.Message)
	}
	return fmt.Sprintf("program error %d (%s)", e.Code, e.Name)
}

// ParseProgramError maps program error data to a Go error using the IDL.
// data is the error bytes from a failed transaction (typically 8-byte u64 error code).
func ParseProgramError(data []byte, idl *idl.IDL) error {
	if len(data) < 8 {
		return fmt.Errorf("program error data too short: %d bytes", len(data))
	}
	code := binary.LittleEndian.Uint64(data[:8])

	for _, ec := range idl.Errors {
		if ec.Code == code {
			return &ProgramError{
				Code:    code,
				Name:    ec.Name,
				Message: ec.Msg,
			}
		}
	}
	return fmt.Errorf("unknown program error: %d", code)
}
