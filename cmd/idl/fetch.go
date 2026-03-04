// Package idlcmd provides IDL CLI commands.
package idlcmd

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"go-solana-anchor/idl"
)

// FetchIDL fetches the IDL from chain for the given program ID and writes to output path.
func FetchIDL(rpcEndpoint, programIDStr, outputPath string) error {
	programID, err := solana.PublicKeyFromBase58(programIDStr)
	if err != nil {
		return fmt.Errorf("invalid program ID: %w", err)
	}

	client := rpc.New(rpcEndpoint)
	defer client.Close()

	// Derive IDL account: base = FindProgramAddress([], programID), idlAddr = CreateWithSeed(base, "anchor:idl", programID)
	base, _, err := solana.FindProgramAddress([][]byte{}, programID)
	if err != nil {
		return fmt.Errorf("find base PDA: %w", err)
	}
	idlAddr, err := solana.CreateWithSeed(base, "anchor:idl", programID)
	if err != nil {
		return fmt.Errorf("create IDL address: %w", err)
	}

	resp, err := client.GetAccountInfo(context.Background(), idlAddr)
	if err != nil {
		return fmt.Errorf("get account info: %w", err)
	}
	if resp == nil || resp.Value == nil {
		return fmt.Errorf("IDL account not found (program may not have IDL deployed)")
	}

	data := resp.Value.Data.GetBinary()

	if len(data) < 8+32+4 {
		return fmt.Errorf("IDL account data too short")
	}
	// Layout: 8-byte discriminator, 32-byte authority, then Borsh Vec<u8>: u32 len + bytes (gzipped IDL JSON)
	payload := data[8:]
	_ = payload[:32] // authority
	payload = payload[32:]
	if len(payload) < 4 {
		return fmt.Errorf("IDL data too short")
	}
	dataLen := binary.LittleEndian.Uint32(payload[:4])
	payload = payload[4:]
	if uint32(len(payload)) < dataLen {
		return fmt.Errorf("IDL data length mismatch")
	}
	compressed := payload[:dataLen]

	var jsonBytes []byte
	zr, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		// Try uncompressed (some IDLs may not be gzipped)
		jsonBytes = compressed
	} else {
		defer zr.Close()
		jsonBytes, err = io.ReadAll(zr)
		if err != nil {
			return fmt.Errorf("decompress IDL: %w", err)
		}
	}

	var raw idl.IDL
	if err := json.Unmarshal(jsonBytes, &raw); err != nil {
		return fmt.Errorf("parse IDL JSON: %w", err)
	}
	return writeIDL(outputPath, &raw)
}

func writeIDL(path string, idl *idl.IDL) error {
	enc, err := json.MarshalIndent(idl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, enc, 0644)
}
