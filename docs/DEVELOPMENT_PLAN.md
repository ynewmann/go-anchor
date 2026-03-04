# Go Solana Anchor — Development Plan

> A development roadmap for building a Golang equivalent of the [official Rust Solana Anchor](https://www.anchor-lang.com/) framework, enabling Anchor-like program interaction and tooling in Go.

---

## Executive Summary

The Rust Anchor framework provides:

1. **Program-side**: Macros and types for writing Solana programs (`#[program]`, `#[derive(Accounts)]`, `#[account]`, constraints)
2. **Client-side**: IDL-based client generation, instruction building, account decoding, CPI support
3. **Tooling**: Build, deploy, IDL generation, key management

**Important**: Solana programs run on the Sealevel VM and compile to SBF (eBPF-derived) bytecode. Rust has first-class SBF support. Go does not officially target SBF for on-chain program deployment. Therefore, this plan focuses on:

- **Primary goal**: Full client-side Anchor parity in Go (IDL parsing, instruction/account encoding, CPI, events, errors)
- **Secondary track**: Research and experimentation for Go→SBF program development (if/when tooling matures)

---

## Rust Anchor Feature Map

| Rust Feature | Description | Go Equivalent |
|--------------|-------------|---------------|
| `declare_id!` | Program ID constant | `ProgramID` type + config |
| `#[program]` | Instruction dispatch module | N/A (program-side) |
| `#[derive(Accounts)]` | Account validation struct | Account struct tags + validation |
| `#[account]` | Account type + serialization | Struct tags + Borsh/custom (de)serializer |
| `Context<T>` | Instruction context (accounts, bumps, program_id) | `Context` struct |
| Account constraints | init, mut, signer, seeds, bump, etc. | Validation functions/tags |
| IDL | Auto-generated from program | Parse IDL JSON → Go types |
| `CpiContext` | Cross-program invocation | CPI helper types |
| Events | `emit!()` | Event encoder + log parsing |
| Errors | Custom error codes | Error type + IDL mapping |

---

## Phase 1: Foundation (Weeks 1–3)

### 1.1 Project Structure

```
go-solana-anchor/
├── idl/           # IDL parsing and types
├── client/        # Program client, instruction builder
├── accounts/      # Account types, encoding, discriminators
├── cpi/           # CPI context and helpers
├── events/        # Event encoding and parsing
├── errors/        # Error code mapping
├── cmd/           # CLI (idl fetch, convert, validate)
└── examples/      # Example programs and clients
```

### 1.2 IDL Support

- [ ] Define Go structs for IDL JSON (Anchor v0.30+ spec)
- [ ] Support `address`, `metadata`, `instructions`, `accounts`, `types`, `events`, `errors`
- [ ] Parse instruction discriminators (8-byte u64 or byte slice)
- [ ] Type mapping: `u8`→`uint8`, `u64`→`uint64`, `i64`→`int64`, `pubkey`→`solana.PublicKey`, `string`→`string`, structs, enums, vecs
- [ ] Handle legacy IDL format (optional conversion layer)

**Reference**: [Anchor IDL Spec](https://www.anchor-lang.com/docs/basics/idl), [Solana IDL Guide](https://solana.com/developers/guides/advanced/idls)

### 1.3 Dependencies

- `github.com/gagliardetto/solana-go` or `github.com/blocto/solana-go-sdk` for RPC, keys, transactions
- `encoding/json` for IDL parsing
- Custom or existing Borsh implementation for Anchor serialization

---

## Phase 2: Account & Instruction Encoding (Weeks 4–6)

### 2.1 Account Discriminator

- [ ] Implement discriminator = first 8 bytes of `sha256("account:<TypeName>")` (Anchor convention)
- [ ] Validate account data starts with expected discriminator
- [ ] Support zero-copy / zero discriminator for large accounts

### 2.2 Instruction Discriminator

- [ ] `sha256("global:<InstructionName>")[0:8]` for instruction selector
- [ ] Build instruction data: discriminator + Borsh-encoded args

### 2.3 Borsh Serialization

- [ ] Borsh encoder/decoder for Go (or adopt `github.com/near/borsh-go`)
- [ ] Support: primitives, `Option`, `Vec`, structs, enums
- [ ] Ensure byte layout matches Anchor/Rust Borsh

### 2.4 Account Types

- [ ] `Account[T]` – owned, deserialized account
- [ ] `AccountLoader[T]` – zero-copy for large accounts
- [ ] `UncheckedAccount` – raw `AccountInfo`
- [ ] `Signer` – validation helper
- [ ] `Program` – program account validation

---

## Phase 3: Client & Instruction Builder (Weeks 7–9)

### 3.1 Program Client

```go
// Target API
client, err := anchor.NewClient(idl, programID, connection, wallet)
// or from IDL bytes
client, err := anchor.NewClientFromIDL(idlJSON, provider)
```

### 3.2 Instruction Builder

```go
// Goal: type-safe instruction building from IDL
ix, err := client.
    Methods.Increment().
    Accounts(counterPubkey).
    Build()
```

- [ ] Generate or reflect instruction accounts from IDL
- [ ] Build `Instruction` with correct discriminator + args
- [ ] Support optional accounts, remaining accounts

### 3.3 Account Fetching

- [ ] `client.Accounts.Counter.Fetch(ctx, pubkey)`
- [ ] `client.Accounts.Counter.FetchMultiple(ctx, pubkeys)`
- [ ] `client.Accounts.Counter.Decode(data)` for raw bytes

---

## Phase 4: CPI & Advanced Features (Weeks 10–12)

### 4.1 CPI Helpers

- [ ] `CpiContext` with program, accounts, optional signer seeds
- [ ] `invoke` / `invoke_signed` wrappers
- [ ] Helper to build CPI instruction from IDL + method name

### 4.2 PDA Derivation

- [ ] `FindProgramAddress(seeds, programID)` – match Solana’s `Pubkey::find_program_address`
- [ ] `FindBump(seeds, programID)` for bump seed
- [ ] Seeds: `[]byte`, `PublicKey`, `uint64`, strings (e.g. `"user"`, `"vault"`)

### 4.3 Events

- [ ] Event discriminator: `sha256("event:<EventName>")[0:8]`
- [ ] Parse logs from transaction for event data
- [ ] Decode event fields per IDL `events` definition

### 4.4 Errors

- [ ] Map program error codes to Go errors
- [ ] Parse `ProgramFailed` / custom error from transaction
- [ ] IDL `errors` → named error types

---

## Phase 5: Constraints Validation (Client-Side) (Weeks 13–14)

Client-side constraints are advisory (actual validation happens on-chain). Implement for DX and debugging:

- [ ] `init` – ensure account is new / will be created
- [ ] `mut` – mark account writable in instruction
- [ ] `signer` – ensure wallet signs
- [ ] `seeds` / `bump` – derive PDA and validate address
- [ ] `has_one` – validate account field matches another account key
- [ ] `address` – validate account matches expected pubkey
- [ ] `owner` – validate account owner
- [ ] `close` – ensure close target is correct
- [ ] `constraint` – custom expression (eval limited to simple checks or documentation)

---

## Phase 6: Tooling & CLI (Weeks 15–16)

### 6.1 CLI Commands

- [ ] `go-anchor idl fetch <program_id> -o idl.json`
- [ ] `go-anchor idl validate idl.json`
- [ ] `go-anchor idl convert legacy.json -o v30.json` (legacy → v0.30)
- [ ] `go-anchor gen client -i idl.json -o pkg/generated/` (optional codegen)

### 6.2 Testing

- [ ] Unit tests for IDL parsing, encoding, discriminators
- [ ] Integration tests against real Anchor programs (e.g. counter, token)
- [ ] Fuzz or property tests for Borsh compatibility with Rust

---

## Phase 7: SPL Integration (Weeks 17–18)

- [ ] Token Program: mint, token account, transfer, etc.
- [ ] Associated Token Account creation
- [ ] Use constraints equivalent to `token::mint`, `token::authority`, `associated_token::*`

---

## Optional: Go On-Chain Programs (Research)

Solana programs are compiled to SBF. Options to explore:

1. **TinyGo / LLVM SBF target** – Check if Go can target SBF via LLVM
2. **Solana Labs tooling** – Monitor for Go/SBF support
3. **Interpreted layer** – Embed a small VM (high overhead, likely impractical)
4. **Codama / IDL-first** – Design programs via IDL; generate Go clients; keep programs in Rust

Recommendation: Focus on client-side parity first; track Solana/LLVM/Go ecosystem for any SBF support.

---

## Success Criteria

- [ ] Parse Anchor v0.30 IDL and legacy IDL
- [ ] Build and send instructions to any Anchor program
- [ ] Fetch and decode accounts by discriminator
- [ ] Support CPI with PDA signers
- [ ] Parse events and errors from transactions
- [ ] CLI for IDL fetch/validate/convert
- [ ] Feature parity checklist vs Rust Anchor client

---

## References

- [Anchor High-Level Overview](https://www.anchor-lang.com/docs/high-level-overview)
- [Anchor Account Constraints](https://www.anchor-lang.com/docs/account-constraints)
- [Anchor IDL](https://www.anchor-lang.com/docs/basics/idl)
- [Solana IDLs Guide](https://solana.com/developers/guides/advanced/idls)
- [solana-go](https://github.com/gagliardetto/solana-go)
- [anchor-go](https://github.com/gagliardetto/anchor-go) (client generation)
