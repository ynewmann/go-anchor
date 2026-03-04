// Package events provides event encoding and parsing for Anchor programs.
package events

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"go-solana-anchor/accounts"
	"go-solana-anchor/idl"
	"go-solana-anchor/internal/borsh"
)

// ParsedEvent holds a decoded event with its name and parsed data.
type ParsedEvent struct {
	Name   string
	Data   interface{}
	RawLog string
}

// EventDecoder decodes event fields from a Borsh decoder.
type EventDecoder func(*borsh.Decoder) (interface{}, error)

// ParseEventFromLogs parses a named event from transaction logs.
// logs are the "logMessages" from the transaction meta.
// eventName is the IDL event name (e.g. "MyEvent").
// decoder is called to decode the event payload (after the 8-byte discriminator).
func ParseEventFromLogs(logs []string, eventName string, idl *idl.IDL, decoder EventDecoder) (*ParsedEvent, error) {
	disc := accounts.EventDiscriminator(eventName)

	for _, log := range logs {
		// Anchor emits events as "Program data: <base64>" where base64 is discriminator + payload.
		if strings.HasPrefix(log, "Program data: ") {
			dataB64 := strings.TrimPrefix(log, "Program data: ")
			data, err := base64.StdEncoding.DecodeString(dataB64)
			if err != nil {
				continue
			}
			if len(data) < 8 {
				continue
			}
			// Check discriminator
			match := true
			for i := 0; i < 8; i++ {
				if data[i] != disc[i] {
					match = false
					break
				}
			}
			if !match {
				continue
			}

			d := borsh.NewDecoder(bytes.NewReader(data[8:]))
			parsed, err := decoder(d)
			if err != nil {
				return nil, fmt.Errorf("decode event %q: %w", eventName, err)
			}
			return &ParsedEvent{Name: eventName, Data: parsed, RawLog: log}, nil
		}

	}
	return nil, fmt.Errorf("event %q not found in logs", eventName)
}
