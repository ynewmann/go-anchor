// Package borsh provides Borsh serialization for Anchor compatibility.
// Layout matches Rust borsh::BorshSerialize / BorshDeserialize.
package borsh

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrShortRead = errors.New("borsh: short read")
)

// Encoder writes Borsh-encoded values.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns an encoder writing to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// EncodeU64 writes a little-endian uint64.
func (e *Encoder) EncodeU64(v uint64) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], v)
	_, err := e.w.Write(buf[:])
	return err
}

// EncodeU32 writes a little-endian uint32.
func (e *Encoder) EncodeU32(v uint32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], v)
	_, err := e.w.Write(buf[:])
	return err
}

// EncodeU8 writes a single byte.
func (e *Encoder) EncodeU8(v uint8) error {
	_, err := e.w.Write([]byte{v})
	return err
}

// EncodeBool writes 1 for true, 0 for false.
func (e *Encoder) EncodeBool(v bool) error {
	var b uint8
	if v {
		b = 1
	}
	return e.EncodeU8(b)
}

// EncodeString writes length (u32) + utf8 bytes.
func (e *Encoder) EncodeString(s string) error {
	b := []byte(s)
	if err := e.EncodeU32(uint32(len(b))); err != nil {
		return err
	}
	_, err := e.w.Write(b)
	return err
}

// EncodeBytes writes length (u32) + raw bytes.
func (e *Encoder) EncodeBytes(b []byte) error {
	if err := e.EncodeU32(uint32(len(b))); err != nil {
		return err
	}
	_, err := e.w.Write(b)
	return err
}

// Decoder reads Borsh-encoded values.
type Decoder struct {
	r io.Reader
}

// NewDecoder returns a decoder reading from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// DecodeU64 reads a little-endian uint64.
func (d *Decoder) DecodeU64() (uint64, error) {
	var buf [8]byte
	if _, err := io.ReadFull(d.r, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}

// DecodeU32 reads a little-endian uint32.
func (d *Decoder) DecodeU32() (uint32, error) {
	var buf [4]byte
	if _, err := io.ReadFull(d.r, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

// DecodeU8 reads a single byte.
func (d *Decoder) DecodeU8() (uint8, error) {
	var buf [1]byte
	if _, err := io.ReadFull(d.r, buf[:]); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// DecodeBool reads 1/0 as true/false.
func (d *Decoder) DecodeBool() (bool, error) {
	b, err := d.DecodeU8()
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

// DecodeString reads length (u32) + utf8 bytes.
func (d *Decoder) DecodeString() (string, error) {
	n, err := d.DecodeU32()
	if err != nil {
		return "", err
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(d.r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

// DecodeBytes reads length (u32) + raw bytes.
func (d *Decoder) DecodeBytes() ([]byte, error) {
	n, err := d.DecodeU32()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(d.r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}
