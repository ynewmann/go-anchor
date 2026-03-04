// Package idl provides parsing and types for Anchor IDL (Interface Description Language).
package idl

import (
	"encoding/json"
	"errors"
)

// IDL is the root structure for an Anchor program's interface (v0.30+ spec).
type IDL struct {
	Address      string            `json:"address"`
	Metadata     Metadata          `json:"metadata"`
	Instructions []Instruction     `json:"instructions"`
	Accounts     []Account         `json:"accounts"`
	Types        []TypeDef         `json:"types,omitempty"`
	Events       []Event           `json:"events,omitempty"`
	Errors       []ErrorCode       `json:"errors,omitempty"`
	Constants    []Constant        `json:"constants,omitempty"`
	Docs         []string          `json:"docs,omitempty"`
}

// Metadata holds program metadata.
type Metadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Spec    string `json:"spec"`
}

// Instruction describes a callable instruction.
type Instruction struct {
	Name          string                `json:"name"`
	Discriminator Discriminator         `json:"discriminator"`
	Accounts      []InstructionAccount  `json:"accounts"`
	Args          []Field               `json:"args"`
	Docs          []string              `json:"docs,omitempty"`
}

// InstructionAccount describes an account required by an instruction.
type InstructionAccount struct {
	Name     string `json:"name"`
	Writable bool   `json:"writable,omitempty"`
	Signer   bool   `json:"signer,omitempty"`
	Optional bool   `json:"optional,omitempty"`
}

// Account describes a custom account type.
type Account struct {
	Name          string        `json:"name"`
	Discriminator Discriminator `json:"discriminator"`
	Docs          []string      `json:"docs,omitempty"`
}

// TypeDef describes a type (struct, enum, etc.).
type TypeDef struct {
	Name string    `json:"name"`
	Type TypeDefTy `json:"type"`
	Docs []string  `json:"docs,omitempty"`
}

// TypeDefTy is the type definition body.
type TypeDefTy struct {
	Kind      string         `json:"kind"`
	Fields    []Field        `json:"fields,omitempty"`
	Variants  []EnumVariant  `json:"variants,omitempty"`
	Alias     *string        `json:"alias,omitempty"`
}

// Field describes a struct field or instruction arg.
type Field struct {
	Name string      `json:"name"`
	Type interface{} `json:"type"`
	Docs []string    `json:"docs,omitempty"`
}

// EnumVariant describes an enum variant.
type EnumVariant struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields,omitempty"`
}

// Event describes an emit event.
type Event struct {
	Name          string        `json:"name"`
	Discriminator Discriminator `json:"discriminator"`
	Fields        []Field       `json:"fields,omitempty"`
	Docs          []string      `json:"docs,omitempty"`
}

// ErrorCode describes a program error.
type ErrorCode struct {
	Code uint64   `json:"code"`
	Name string   `json:"name"`
	Msg  string   `json:"msg,omitempty"`
	Docs []string `json:"docs,omitempty"`
}

// Constant describes a program constant.
type Constant struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Discriminator is an 8-byte instruction/account/event selector.
type Discriminator [8]byte

// UnmarshalJSON implements json.Unmarshaler for Discriminator.
func (d *Discriminator) UnmarshalJSON(data []byte) error {
	var raw []json.Number
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 8 {
		return errors.New("discriminator must be 8 bytes")
	}
	for i := 0; i < 8; i++ {
		n, err := raw[i].Int64()
		if err != nil || n < 0 || n > 255 {
			return errors.New("discriminator bytes must be 0-255")
		}
		d[i] = byte(n)
	}
	return nil
}

// MarshalJSON implements json.Marshaler for Discriminator.
func (d Discriminator) MarshalJSON() ([]byte, error) {
	arr := make([]int, 8)
	for i := 0; i < 8; i++ {
		arr[i] = int(d[i])
	}
	return json.Marshal(arr)
}

// Bytes returns the discriminator as a slice.
func (d Discriminator) Bytes() []byte {
	return d[:]
}

// Parse unmarshals IDL from JSON bytes.
func Parse(data []byte) (*IDL, error) {
	var idl IDL
	if err := json.Unmarshal(data, &idl); err != nil {
		return nil, err
	}
	return &idl, nil
}
