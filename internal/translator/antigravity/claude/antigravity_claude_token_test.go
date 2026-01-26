package claude

import (
	"context"
	"testing"
)

func TestConvertAntigravityResponseToClaude_TokenCalculation(t *testing.T) {
	requestJSON := []byte(`{"model": "test-model"}`)

	// Response with usage metadata
	// promptTokenCount: 100
	// cachedContentTokenCount: 20
	// candidatesTokenCount: 10
	// thoughtsTokenCount: 5
	// totalTokenCount: 115 (100 + 10 + 5)

	responseJSON := []byte(`{
		"response": {
			"candidates": [{"finishReason": "STOP"}],
			"usageMetadata": {
				"promptTokenCount": 100,
				"cachedContentTokenCount": 20,
				"candidatesTokenCount": 10,
				"thoughtsTokenCount": 5,
				"totalTokenCount": 115
			}
		}
	}`)

	var param any
	ctx := context.Background()

	// Initialize params
	param = &Params{
		HasFirstResponse: true,
		HasContent: true, // Force final events
        ResponseType: 1, // Content type
	}

	ConvertAntigravityResponseToClaude(ctx, "test-model", requestJSON, requestJSON, responseJSON, &param)

	params := param.(*Params)

	// Verify PromptTokenCount
	// Before fix: 100 - 20 = 80
	// After fix: 100
	if params.PromptTokenCount != 100 {
		t.Errorf("Expected PromptTokenCount to be 100, got %d", params.PromptTokenCount)
	}

	// Verify CandidatesTokenCount
	// Before fix: Total(115) - Prompt(80) - Thoughts(5) = 30 (WRONG)
	// After fix: Total(115) - Prompt(100) - Thoughts(5) = 10 (CORRECT)
	if params.CandidatesTokenCount != 10 {
		t.Errorf("Expected CandidatesTokenCount to be 10, got %d", params.CandidatesTokenCount)
	}
}
