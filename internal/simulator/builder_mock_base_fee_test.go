// Copyright 2025 Erst Users
// SPDX-License-Identifier: Apache-2.0

package simulator

import "testing"

func TestSimulationRequestBuilder_WithMockBaseFee(t *testing.T) {
	req, err := NewSimulationRequestBuilder().
		WithEnvelopeXDR("envelope").
		WithResultMetaXDR("result").
		WithMockBaseFee(5000).
		Build()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if req.MockBaseFee == nil {
		t.Fatal("expected MockBaseFee to be set")
	}
	if *req.MockBaseFee != 5000 {
		t.Fatalf("expected MockBaseFee=5000, got %d", *req.MockBaseFee)
	}
}

func TestSimulationRequestBuilder_WithMockBaseFee_ZeroValue(t *testing.T) {
	req, err := NewSimulationRequestBuilder().
		WithEnvelopeXDR("envelope").
		WithResultMetaXDR("result").
		WithMockBaseFee(0).
		Build()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if req.MockBaseFee == nil {
		t.Fatal("expected MockBaseFee to be set")
	}
	if *req.MockBaseFee != 0 {
		t.Fatalf("expected MockBaseFee=0, got %d", *req.MockBaseFee)
	}
}

func TestSimulationRequestBuilder_ResetClearsMockBaseFee(t *testing.T) {
	req, err := NewSimulationRequestBuilder().
		WithEnvelopeXDR("envelope").
		WithResultMetaXDR("result").
		WithMockBaseFee(123).
		Reset().
		WithEnvelopeXDR("envelope2").
		WithResultMetaXDR("result2").
		Build()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if req.MockBaseFee != nil {
		t.Fatalf("expected MockBaseFee to be cleared, got %v", req.MockBaseFee)
	}
}
