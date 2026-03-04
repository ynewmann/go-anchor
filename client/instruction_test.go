package client

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"go-solana-sdk/idl"
)

func TestBuildInstruction(t *testing.T) {
	raw := `{
		"address": "6khKp4BeJpCjBY1Eh39ybiqbfRnrn2UzWeUARjQLXYRC",
		"metadata": {"name": "counter", "version": "0.1.0", "spec": "0.1.0"},
		"instructions": [
			{
				"name": "increment",
				"discriminator": [11, 18, 104, 9, 104, 174, 59, 33],
				"accounts": [{"name": "counter", "writable": true}],
				"args": []
			}
		],
		"accounts": [{"name": "Counter", "discriminator": [255, 176, 4, 245, 188, 253, 124, 25]}]
	}`
	idl, err := idl.Parse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	programID, _ := solana.PublicKeyFromBase58("6khKp4BeJpCjBY1Eh39ybiqbfRnrn2UzWeUARjQLXYRC")
	prog := NewProgram(idl, programID, rpc.New("https://api.mainnet-beta.solana.com"))

	counterAddr, _ := solana.PublicKeyFromBase58("11111111111111111111111111111111")
	ix, err := prog.BuildInstruction("increment", nil, map[string]solana.PublicKey{"counter": counterAddr})
	if err != nil {
		t.Fatal(err)
	}
	if ix.ProgramID() != programID {
		t.Errorf("wrong program ID")
	}
	data, _ := ix.Data()
	if len(data) != 8 {
		t.Errorf("expected 8 bytes (discriminator only), got %d", len(data))
	}
}
