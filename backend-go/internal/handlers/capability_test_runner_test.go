package handlers

import "testing"

func TestShouldRunRedirectVerification(t *testing.T) {
	tests := []struct {
		name               string
		protocols          []string
		sourceTab          string
		channelServiceType string
		want               bool
	}{
		{
			name:               "cross protocol runs redirect verification even for unrelated base protocol",
			protocols:          []string{"messages"},
			sourceTab:          "responses",
			channelServiceType: "chat",
			want:               true,
		},
		{
			name:               "cross protocol runs redirect verification for source protocol",
			protocols:          []string{"responses"},
			sourceTab:          "responses",
			channelServiceType: "chat",
			want:               true,
		},
		{
			name:               "explicit virtual protocol runs redirect verification",
			protocols:          []string{"responses->chat"},
			sourceTab:          "responses",
			channelServiceType: "chat",
			want:               true,
		},
		{
			name:               "cross protocol runs redirect verification even when protocols contain other virtual protocol",
			protocols:          []string{"messages->chat"},
			sourceTab:          "responses",
			channelServiceType: "chat",
			want:               true,
		},
		{
			name:               "all base protocols run redirect verification for cross protocol channel",
			protocols:          []string{"messages", "responses", "chat", "gemini"},
			sourceTab:          "responses",
			channelServiceType: "chat",
			want:               true,
		},
		{
			name:               "same source and channel protocol does not run redirect verification implicitly",
			protocols:          []string{"responses"},
			sourceTab:          "responses",
			channelServiceType: "responses",
			want:               false,
		},
		{
			name:               "same source unrelated virtual protocol does not run redirect verification",
			protocols:          []string{"messages->chat"},
			sourceTab:          "responses",
			channelServiceType: "responses",
			want:               false,
		},
		{
			name:               "same source explicit virtual protocol runs redirect verification",
			protocols:          []string{"responses->responses"},
			sourceTab:          "responses",
			channelServiceType: "responses",
			want:               true,
		},
		{
			name:               "empty source tab does not run redirect verification",
			protocols:          []string{"responses"},
			sourceTab:          "",
			channelServiceType: "chat",
			want:               false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldRunRedirectVerification(tt.protocols, tt.sourceTab, tt.channelServiceType)
			if got != tt.want {
				t.Fatalf("shouldRunRedirectVerification(%v, %q, %q) = %v, want %v", tt.protocols, tt.sourceTab, tt.channelServiceType, got, tt.want)
			}
		})
	}
}
