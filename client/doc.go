// Package client provides the Anchor program client and instruction builder.
//
// # Phase 5: Constraints (Client-Side)
//
// Anchor account constraints (init, mut, signer, seeds, bump, has_one, address,
// owner, close, etc.) are enforced on-chain. This client does not validate
// constraints when building instructions—it is the caller's responsibility to:
//
//   - Ensure init accounts are created with the correct space and discriminator
//   - Include writable accounts when the IDL specifies mut
//   - Ensure signers sign the transaction
//   - Derive PDAs with the correct seeds and bump when required
//   - Validate account relationships (has_one, address, owner) if needed for correctness
//
// Use this client to build and send instructions; rely on the program for
// constraint validation. Client-side constraint checks can be added for DX
// and debugging if desired.
package client
