# go-anchor

Go client library and tooling for interacting with [Solana Anchor](https://www.anchor-lang.com/) programs. Provides IDL parsing, instruction building, account decoding, CPI helpers, and a CLI—all compatible with the official Rust Anchor framework.

## Features

- **IDL parsing** — Parse Anchor v0.30+ IDL JSON (instructions, accounts, types, events, errors)
- **Discriminators** — Account, instruction, and event discriminators matching Anchor's `sha256` convention
- **Instruction builder** — Build instructions with correct discriminators and Borsh-encoded args
- **Account decoding** — Decode account data with discriminator validation
- **PDA derivation** — `FindProgramAddress` for Program Derived Addresses
- **CPI** — Cross-program invocation context and helpers
- **Events & errors** — Parse events from logs and map program errors from IDL
- **SPL** — Token and Associated Token Account helpers
- **CLI** — `go-anchor idl fetch|validate|convert`

## Installation

```bash
go get github.com/ynewmann/go-anchor
```

The module name is `go-solana-sdk`; use it in imports: `go-solana-sdk/idl`, `go-solana-sdk/client`, etc.

## Usage

### Program client

```go
import (
    "github.com/gagliardetto/solana-go/rpc"
    "go-solana-sdk/client"
    "go-solana-sdk/idl"
)

idlData, _ := os.ReadFile("idl.json")
parsed, _ := idl.Parse(idlData)

rpcClient := rpc.New("https://api.mainnet-beta.solana.com")
programID := solana.MustPublicKeyFromBase58(parsed.Address)
prog := client.NewProgram(parsed, programID, rpcClient)

// Build an instruction
ix, err := prog.BuildInstruction("increment", nil, map[string]solana.PublicKey{
    "counter": counterPubkey,
})
```

### CLI

```bash
# Build the CLI
go build -o go-anchor ./cmd/go-anchor

# Fetch IDL from chain
go-anchor idl fetch <program_id> -o idl.json

# Validate IDL file
go-anchor idl validate idl.json

# Convert legacy IDL to v0.30
go-anchor idl convert legacy.json -o v30.json
```

Set `RPC_URL` for a custom RPC endpoint (default: mainnet-beta).

## Project layout

```
├── idl/           # IDL parsing and types
├── client/        # Program client, instruction builder
├── accounts/      # Discriminators, instruction encoding
├── cpi/           # CPI context and helpers
├── pda/           # PDA derivation
├── events/        # Event parsing
├── errors/        # Error code mapping
├── pkg/spl/       # SPL token helpers
├── internal/borsh # Borsh serialization
├── cmd/go-anchor/ # CLI
└── docs/          # Development plan
```

## Dependencies

- [github.com/gagliardetto/solana-go](https://github.com/gagliardetto/solana-go) — RPC, keys, transactions

## License

Apache-2.0
