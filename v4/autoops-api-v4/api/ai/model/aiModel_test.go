package model

import (
	"encoding/json"
	"testing"
)

func TestAiChatRequestUnmarshalSupportsSessionID(t *testing.T) {
	raw := []byte(`{
		"modelId": 1,
		"messages": [{"role":"user","content":"hello"}],
		"sessionID": "legacy-session"
	}`)

	var req AiChatRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if req.SessionID != "legacy-session" {
		t.Fatalf("expected sessionID fallback, got %q", req.SessionID)
	}
}

func TestAiChatRequestUnmarshalPrefersSessionId(t *testing.T) {
	raw := []byte(`{
		"modelId": 1,
		"messages": [{"role":"user","content":"hello"}],
		"sessionId": "new-session",
		"sessionID": "legacy-session"
	}`)

	var req AiChatRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if req.SessionID != "new-session" {
		t.Fatalf("expected sessionId to win, got %q", req.SessionID)
	}
}
