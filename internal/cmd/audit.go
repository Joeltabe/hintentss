// Copyright 2025 Erst Users
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// AttestationCertificate represents a single X.509 certificate in the
// hardware attestation chain. Certificates are ordered leaf-to-root.
type AttestationCertificate struct {
	PEM     string `json:"pem"`
	Subject string `json:"subject"`
	Issuer  string `json:"issuer"`
	Serial  string `json:"serial"`
}

// HardwareAttestation contains the full attestation chain retrieved from
// an HSM or hardware security token. When present in an AuditLog it
// provides cryptographic proof that the signing key resides on a
// hardware device and is non-exportable.
type HardwareAttestation struct {
	Certificates    []AttestationCertificate `json:"certificates"`
	TokenInfo       string                  `json:"token_info"`
	KeyNonExportable bool                   `json:"key_non_exportable"`
	RetrievedAt     string                  `json:"retrieved_at"`
}

// AuditLog represents the signed audit trail of a transaction simulation
type AuditLog struct {
	Version         string    `json:"version"`
	Timestamp       time.Time `json:"timestamp"`
	TransactionHash string    `json:"transaction_hash"`
	TraceHash       string    `json:"trace_hash"`
	Signature       string    `json:"signature"`
	PublicKey       string    `json:"public_key"`
	Payload         Payload   `json:"payload"`
}

// Payload contains the actual trace data.
// This struct is serialized using canonical JSON (sorted keys) to ensure
// deterministic hashing across different platforms and Go versions.
type Payload struct {
	EnvelopeXdr   string   `json:"envelope_xdr"`
	ResultMetaXdr string   `json:"result_meta_xdr"`
	Events        []string `json:"events"`
	Logs          []string `json:"logs"`
}

// Generate creates a signed audit log from the simulation results
func Generate(txHash string, envelopeXdr, resultMetaXdr string, events, logs []string, privateKeyHex string) (*AuditLog, error) {
	// 1. Construct Payload
	payload := Payload{
		EnvelopeXdr:   envelopeXdr,
		ResultMetaXdr: resultMetaXdr,
		Events:        events,
		Logs:          logs,
	}

	// 2. Construct the hash input.
	// When hardware attestation is present, it is included in the hash
	// so that stripping it would invalidate the signature.
	type hashInput struct {
		Payload             Payload              `json:"payload"`
		HardwareAttestation *HardwareAttestation `json:"hardware_attestation,omitempty"`
	}

	hi := hashInput{Payload: payload}
	if opts != nil && opts.HardwareAttestation != nil {
		hi.HardwareAttestation = opts.HardwareAttestation
	}

	payloadBytes, err := json.Marshal(hi)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 3. Calculate Trace Hash (SHA256)
	hash := sha256.Sum256(payloadBytes)
	traceHashHex := hex.EncodeToString(hash[:])

	// 4. Parse Private Key
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key hex: %w", err)
	}

	if len(privKeyBytes) != ed25519.PrivateKeySize && len(privKeyBytes) != ed25519.SeedSize {
		return nil, fmt.Errorf("invalid private key length: %d", len(privKeyBytes))
	}

	var privateKey ed25519.PrivateKey
	if len(privKeyBytes) == ed25519.SeedSize {
		privateKey = ed25519.NewKeyFromSeed(privKeyBytes)
	} else {
		privateKey = ed25519.PrivateKey(privKeyBytes)
	}

	// 5. Sign the Trace Hash
	// We sign the hash of the payload to ensure integrity.
	signature := ed25519.Sign(privateKey, hash[:])

	return &AuditLog{
		Version:         "1.0.0",
		Timestamp:       time.Now().UTC(),
		TransactionHash: txHash,
		TraceHash:       traceHashHex,
		Signature:       hex.EncodeToString(signature),
		PublicKey:       hex.EncodeToString(privateKey.Public().(ed25519.PublicKey)),
		Payload:         payload,
	}, nil
}
