package idlcmd

import (
	"encoding/json"
	"fmt"
	"os"

	"go-solana-anchor/idl"
)

// ConvertIDL performs basic legacy IDL to v0.30 format conversion.
func ConvertIDL(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}

	// Legacy IDL uses "version" and different structure. Normalize to v0.30.
	// Try parsing as v0.30 first
	parsed, err := idl.Parse(data)
	if err == nil {
		// Already valid v0.30
		enc, err := json.MarshalIndent(parsed, "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile(outputPath, enc, 0644)
	}

	// Basic legacy conversion: map common fields
	out := &idl.IDL{}
	if v, ok := raw["address"].(string); ok {
		out.Address = v
	}
	if m, ok := raw["metadata"].(map[string]interface{}); ok {
		if v, ok := m["name"].(string); ok {
			out.Metadata.Name = v
		}
		if v, ok := m["version"].(string); ok {
			out.Metadata.Version = v
		}
		out.Metadata.Spec = "0.30.0"
	}
	// Instructions
	if arr, ok := raw["instructions"].([]interface{}); ok {
		for _, it := range arr {
			item, ok := it.(map[string]interface{})
			if !ok {
				continue
			}
			inst := idl.Instruction{}
			if v, ok := item["name"].(string); ok {
				inst.Name = v
			}
			if arr, ok := item["accounts"].([]interface{}); ok {
				for _, a := range arr {
					acc, ok := a.(map[string]interface{})
					if !ok {
						continue
					}
					ia := idl.InstructionAccount{}
					if v, ok := acc["name"].(string); ok {
						ia.Name = v
					}
					if v, ok := acc["writable"].(bool); ok {
						ia.Writable = v
					}
					if v, ok := acc["signer"].(bool); ok {
						ia.Signer = v
					}
					inst.Accounts = append(inst.Accounts, ia)
				}
			}
			if arr, ok := item["args"].([]interface{}); ok {
				for _, f := range arr {
					field, ok := f.(map[string]interface{})
					if !ok {
						continue
					}
					inst.Args = append(inst.Args, idl.Field{
						Name: getStr(field, "name"),
						Type: field["type"],
					})
				}
			}
			// Legacy discriminator may be array
			if arr, ok := item["discriminator"].([]interface{}); ok && len(arr) == 8 {
				for i := 0; i < 8; i++ {
					if n, ok := toUint8(arr[i]); ok {
						inst.Discriminator[i] = n
					}
				}
			}
			out.Instructions = append(out.Instructions, inst)
		}
	}
	// Accounts
	if arr, ok := raw["accounts"].([]interface{}); ok {
		for _, it := range arr {
			item, ok := it.(map[string]interface{})
			if !ok {
				continue
			}
			acc := idl.Account{}
			if v, ok := item["name"].(string); ok {
				acc.Name = v
			}
			if arr, ok := item["discriminator"].([]interface{}); ok && len(arr) == 8 {
				for i := 0; i < 8; i++ {
					if n, ok := toUint8(arr[i]); ok {
						acc.Discriminator[i] = n
					}
				}
			}
			out.Accounts = append(out.Accounts, acc)
		}
	}
	// Types
	if arr, ok := raw["types"].([]interface{}); ok {
		for _, it := range arr {
			item, ok := it.(map[string]interface{})
			if !ok {
				continue
			}
			td := idl.TypeDef{}
			if v, ok := item["name"].(string); ok {
				td.Name = v
			}
			if ty, ok := item["type"].(map[string]interface{}); ok {
				td.Type.Kind = getStr(ty, "kind")
				if arr, ok := ty["fields"].([]interface{}); ok {
					for _, f := range arr {
						field, ok := f.(map[string]interface{})
						if !ok {
							continue
						}
						td.Type.Fields = append(td.Type.Fields, idl.Field{
							Name: getStr(field, "name"),
							Type: field["type"],
						})
					}
				}
			}
			out.Types = append(out.Types, td)
		}
	}
	// Errors
	if arr, ok := raw["errors"].([]interface{}); ok {
		for _, it := range arr {
			item, ok := it.(map[string]interface{})
			if !ok {
				continue
			}
			ec := idl.ErrorCode{}
			if v, ok := toUint64(item["code"]); ok {
				ec.Code = v
			}
			if v, ok := item["name"].(string); ok {
				ec.Name = v
			}
			if v, ok := item["msg"].(string); ok {
				ec.Msg = v
			}
			out.Errors = append(out.Errors, ec)
		}
	}

	enc, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, enc, 0644)
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func toUint8(v interface{}) (uint8, bool) {
	switch x := v.(type) {
	case float64:
		if x >= 0 && x <= 255 {
			return uint8(x), true
		}
	case int:
		if x >= 0 && x <= 255 {
			return uint8(x), true
		}
	}
	return 0, false
}

func toUint64(v interface{}) (uint64, bool) {
	switch x := v.(type) {
	case float64:
		if x >= 0 {
			return uint64(x), true
		}
	case int:
		if x >= 0 {
			return uint64(x), true
		}
	}
	return 0, false
}
