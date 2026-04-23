package main

import "testing"

func TestApplyPresetPreservesUnmanagedFields(t *testing.T) {
	root := map[string]any{
		"permissions": map[string]any{
			"allow_file_access": true,
		},
		"env": map[string]any{
			"FOO":            "bar",
			"API_TIMEOUT_MS": "1",
			"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": 0,
		},
	}

	applyPreset(root, providerPresets["openrouter"], "sk-test", "")

	env := root["env"].(map[string]any)
	if env["FOO"] != "bar" {
		t.Fatalf("expected unmanaged env to be preserved")
	}
	if _, ok := env["API_TIMEOUT_MS"]; ok {
		t.Fatalf("expected stale managed key to be removed")
	}
	if got := env["ANTHROPIC_BASE_URL"]; got != "https://openrouter.ai/api" {
		t.Fatalf("unexpected base url: %v", got)
	}
	if got := env["ANTHROPIC_AUTH_TOKEN"]; got != "sk-test" {
		t.Fatalf("unexpected auth token: %v", got)
	}
}

func TestApplyPresetOverrideModel(t *testing.T) {
	root := map[string]any{}
	applyPreset(root, providerPresets["minimax"], "sk-test", "custom-model")

	env := root["env"].(map[string]any)
	for _, key := range []string{
		"ANTHROPIC_MODEL",
		"ANTHROPIC_DEFAULT_HAIKU_MODEL",
		"ANTHROPIC_DEFAULT_SONNET_MODEL",
		"ANTHROPIC_DEFAULT_OPUS_MODEL",
	} {
		if got := env[key]; got != "custom-model" {
			t.Fatalf("expected %s to be custom-model, got %v", key, got)
		}
	}
}

func TestDetectProvider(t *testing.T) {
	cases := []struct {
		baseURL string
		model   string
		want    string
	}{
		{baseURL: "https://api.minimaxi.com/anthropic", want: "minimax"},
		{baseURL: "https://openrouter.ai/api", want: "openrouter"},
		{baseURL: "https://opencode.ai/zen/go", model: "opencode-go/minimax-m2.7", want: "opencode-go"},
		{baseURL: "https://example.com", want: "custom"},
	}

	for _, tc := range cases {
		if got := detectProvider(tc.baseURL, tc.model); got != tc.want {
			t.Fatalf("detectProvider(%q, %q) = %q, want %q", tc.baseURL, tc.model, got, tc.want)
		}
	}
}
