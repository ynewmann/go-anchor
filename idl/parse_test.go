package idl

import (
	"encoding/json"
	"testing"
)

func TestParse(t *testing.T) {
	raw := `{
		"address": "6khKp4BeJpCjBY1Eh39ybiqbfRnrn2UzWeUARjQLXYRC",
		"metadata": {
			"name": "counter",
			"version": "0.1.0",
			"spec": "0.1.0"
		},
		"instructions": [
			{
				"name": "increment",
				"discriminator": [11, 18, 104, 9, 104, 174, 59, 33],
				"accounts": [{"name": "counter", "writable": true}],
				"args": []
			}
		],
		"accounts": [
			{
				"name": "Counter",
				"discriminator": [255, 176, 4, 245, 188, 253, 124, 25]
			}
		],
		"types": [
			{
				"name": "Counter",
				"type": {
					"kind": "struct",
					"fields": [{"name": "count", "type": "u64"}]
				}
			}
		]
	}`

	idl, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if idl.Address != "6khKp4BeJpCjBY1Eh39ybiqbfRnrn2UzWeUARjQLXYRC" {
		t.Errorf("unexpected address: %s", idl.Address)
	}
	if idl.Metadata.Name != "counter" {
		t.Errorf("unexpected metadata name: %s", idl.Metadata.Name)
	}
	if len(idl.Instructions) != 1 {
		t.Fatalf("expected 1 instruction, got %d", len(idl.Instructions))
	}
	if idl.Instructions[0].Name != "increment" {
		t.Errorf("unexpected instruction name: %s", idl.Instructions[0].Name)
	}
	if idl.Instructions[0].Discriminator[0] != 11 {
		t.Errorf("unexpected discriminator[0]: %d", idl.Instructions[0].Discriminator[0])
	}
	if len(idl.Accounts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(idl.Accounts))
	}
	if idl.Accounts[0].Name != "Counter" {
		t.Errorf("unexpected account name: %s", idl.Accounts[0].Name)
	}
	if len(idl.Types) != 1 {
		t.Fatalf("expected 1 type, got %d", len(idl.Types))
	}
	if idl.Types[0].Name != "Counter" {
		t.Errorf("unexpected type name: %s", idl.Types[0].Name)
	}
}

func TestDiscriminatorRoundTrip(t *testing.T) {
	d := Discriminator{11, 18, 104, 9, 104, 174, 59, 33}
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var d2 Discriminator
	if err := json.Unmarshal(data, &d2); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	for i := 0; i < 8; i++ {
		if d[i] != d2[i] {
			t.Errorf("byte %d: got %d, want %d", i, d2[i], d[i])
		}
	}
}
