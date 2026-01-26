package gemini

import (
	"context"
	"testing"
)

func TestConvertAntigravityResponseToGemini_AltFormat(t *testing.T) {
	// Case 1: alt is empty string
	// Expected: extract "response" field
	ctx := context.WithValue(context.Background(), "alt", "")
	rawJSON := []byte(`{"response": {"candidates": [{"content": {"parts": [{"text": "Hello"}]}}]}}`)

	result := ConvertAntigravityResponseToGemini(ctx, "model", nil, nil, rawJSON, nil)
	if len(result) != 1 {
		t.Fatalf("Expected 1 chunk, got %d", len(result))
	}
	expected1 := `{"candidates": [{"content": {"parts": [{"text": "Hello"}]}}]}`
	if result[0] != expected1 {
		t.Errorf("Expected %s, got %s", expected1, result[0])
	}

	// Case 2: alt is non-empty string (e.g. "json")
	// Expected: parse array of objects, extract "response" from each
	ctx2 := context.WithValue(context.Background(), "alt", "json")
	// Input is an array of response wrappers
	rawArrayJSON := []byte(`[
		{"response": {"text": "Part 1"}},
		{"response": {"text": "Part 2"}}
	]`)

	result2 := ConvertAntigravityResponseToGemini(ctx2, "model", nil, nil, rawArrayJSON, nil)
	if len(result2) != 1 {
		t.Fatalf("Expected 1 chunk, got %d", len(result2))
	}

	// This is where the bug should manifest.
	// Current code uses uninitialized 'chunk', so it likely returns "[]"
	expected2 := `[{"text": "Part 1"},{"text": "Part 2"}]`
	if result2[0] != expected2 {
		t.Errorf("Expected %s, got %s", expected2, result2[0])
	}
}
