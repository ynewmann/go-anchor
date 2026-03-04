package events

import (
	"encoding/base64"
	"testing"

	"go-solana-sdk/accounts"
	"go-solana-sdk/idl"
	"go-solana-sdk/internal/borsh"
)

func TestParseEventFromLogs(t *testing.T) {
	disc := accounts.EventDiscriminator("MyEvent")
	// Borsh-encoded bytes: u32 length(3) + 3 bytes
	payload := []byte{3, 0, 0, 0, 1, 2, 3}
	full := append(disc[:], payload...)
	logLine := "Program data: " + base64.StdEncoding.EncodeToString(full)
	logs := []string{"Program 111... invoke [1]", logLine}

	ev, err := ParseEventFromLogs(logs, "MyEvent", &idl.IDL{}, func(d *borsh.Decoder) (interface{}, error) {
		return d.DecodeBytes()
	})
	if err != nil {
		t.Fatal(err)
	}
	if ev.Name != "MyEvent" {
		t.Errorf("got name %q", ev.Name)
	}
}

func TestParseEventFromLogs_NotFound(t *testing.T) {
	logs := []string{"Program 111... invoke [1]", "Program log: something"}
	_, err := ParseEventFromLogs(logs, "NoSuchEvent", &idl.IDL{}, func(*borsh.Decoder) (interface{}, error) {
		return nil, nil
	})
	if err == nil {
		t.Error("expected error for missing event")
	}
}
