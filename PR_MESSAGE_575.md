# #575 Add support for custom base fees in simulation requests

## Summary
- Added/verified builder support for custom base fee overrides in simulation requests via `WithMockBaseFee(baseFee uint32)`.
- Ensured `Build()` propagates the override into `SimulationRequest.MockBaseFee`.
- Ensured `Reset()` clears any configured custom base fee override.
- Added focused tests for:
  - non-zero override
  - zero-value override
  - reset-clearing behavior

## Why
This allows developers to override the default baseline inclusion fee so local simulations can model surge-pricing and fee-sufficiency conditions deterministically.

## Files changed
- `internal/simulator/builder_mock_base_fee_test.go`

## Validation
- `go test ./internal/simulator -run TestSimulationRequestBuilder_WithMockBaseFee -count=1`
- `go test ./internal/simulator -run TestSimulationRequestBuilder_WithMockBaseFee_ZeroValue -count=1`
- `go test ./internal/simulator -run TestSimulationRequestBuilder_ResetClearsMockBaseFee -count=1`

## Attachment (Proof)
<!-- Paste uploaded image URL below -->
![proof](<PASTE_ATTACHMENT_URL_HERE>)

### How to get the attachment URL
1. Take a screenshot of successful test output (or relevant terminal proof).
2. Open your GitHub PR description (or a PR comment) and drag-drop the image.
3. GitHub uploads it and inserts a Markdown image URL.
4. Copy that URL and replace `<PASTE_ATTACHMENT_URL_HERE>` above.
