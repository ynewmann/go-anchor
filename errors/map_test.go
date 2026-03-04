package errors

import (
	"encoding/binary"
	"testing"

	"go-solana-anchor/idl"
)

func TestParseProgramError(t *testing.T) {
	idl := &idl.IDL{
		Errors: []idl.ErrorCode{
			{Code: 6000, Name: "InsufficientFunds", Msg: "Insufficient funds"},
		},
	}
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, 6000)
	err := ParseProgramError(data, idl)
	if err == nil {
		t.Fatal("expected error")
	}
	pe, ok := err.(*ProgramError)
	if !ok {
		t.Fatalf("expected *ProgramError, got %T", err)
	}
	if pe.Code != 6000 || pe.Name != "InsufficientFunds" {
		t.Errorf("got Code=%d Name=%s", pe.Code, pe.Name)
	}
}

func TestParseProgramError_Unknown(t *testing.T) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, 99999)
	err := ParseProgramError(data, &idl.IDL{})
	if err == nil {
		t.Fatal("expected error")
	}
}
