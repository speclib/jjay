package workspace

import "testing"

func TestSlug(t *testing.T) {
	tests := []struct {
		name   string
		prompt string
		want   string
	}{
		{"typical", "add dark mode to the settings page", "dark-mode-settings-page"},
		{"stopwords dropped", "the a an to for of", "proposal"},
		{"punctuation stripped", "Fix login (OAuth) flow now!", "fix-login-oauth-flow"},
		{"token cap keeps first N salient", "one two three four five six", "one-two-three-four"},
		{"empty prompt", "", "proposal"},
		{"already kebab", "dark-mode", "dark-mode"},
	}
	for _, tt := range tests {
		if got := Slug(tt.prompt); got != tt.want {
			t.Errorf("Slug(%q) = %q, want %q", tt.prompt, got, tt.want)
		}
	}
}

func TestSlug_LengthCap(t *testing.T) {
	// Four long tokens that join past maxSlugLen; result must be capped and not
	// end on a dash.
	got := Slug("aaaaaaaaaa bbbbbbbbbb cccccccccc dddddddddd eeeeeeeeee")
	if len(got) > maxSlugLen {
		t.Errorf("Slug length %d exceeds cap %d: %q", len(got), maxSlugLen, got)
	}
	if got[len(got)-1] == '-' {
		t.Errorf("Slug must not end with a dash: %q", got)
	}
}

func TestUniqueSlug(t *testing.T) {
	tests := []struct {
		name  string
		base  string
		taken map[string]bool
		want  string
	}{
		{"no collision", "dark-mode", map[string]bool{}, "dark-mode"},
		{"one collision", "dark-mode", map[string]bool{"dark-mode": true}, "dark-mode-2"},
		{
			"chain of collisions",
			"dark-mode",
			map[string]bool{"dark-mode": true, "dark-mode-2": true, "dark-mode-3": true},
			"dark-mode-4",
		},
	}
	for _, tt := range tests {
		if got := UniqueSlug(tt.base, tt.taken); got != tt.want {
			t.Errorf("UniqueSlug(%q, %v) = %q, want %q", tt.base, tt.taken, got, tt.want)
		}
	}
}
