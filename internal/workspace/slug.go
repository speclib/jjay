package workspace

import (
	"strings"
)

// maxSlugTokens caps how many salient tokens a slug keeps, so a long prompt
// still yields a short, readable handle.
const maxSlugTokens = 4

// maxSlugLen caps the total slug length (characters) as a final guard against
// a handful of very long tokens producing an unwieldy name.
const maxSlugLen = 40

// stopwords are common low-signal words dropped from a prompt before slugging,
// so the salient tokens survive the token cap. Kept small and deterministic.
var stopwords = map[string]bool{
	"a": true, "an": true, "the": true, "to": true, "for": true, "of": true,
	"in": true, "on": true, "at": true, "and": true, "or": true, "with": true,
	"add": true, "is": true, "are": true, "be": true, "by": true, "as": true,
	"into": true, "from": true, "that": true, "this": true, "it": true,
}

// Slug turns a free-text prompt into a short kebab-case handle by deterministic
// code (no AI): lowercase, strip punctuation, split on whitespace, drop
// stopwords, keep the first maxSlugTokens salient tokens, join with "-", and
// cap the total length. If every token is a stopword (or the prompt is empty),
// it falls back to "proposal" so the result is never empty.
func Slug(prompt string) string {
	lowered := strings.ToLower(prompt)

	// Replace any non-alphanumeric run with a space, so punctuation becomes a
	// token boundary rather than sticking to a word.
	var b strings.Builder
	for _, r := range lowered {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}

	var tokens []string
	for _, tok := range strings.Fields(b.String()) {
		if stopwords[tok] {
			continue
		}
		tokens = append(tokens, tok)
		if len(tokens) == maxSlugTokens {
			break
		}
	}

	if len(tokens) == 0 {
		return "proposal"
	}

	slug := strings.Join(tokens, "-")
	if len(slug) > maxSlugLen {
		slug = strings.TrimRight(slug[:maxSlugLen], "-")
	}
	return slug
}

// UniqueSlug returns base, or base with a numeric suffix ("-2", "-3", …)
// appended, such that the result is not present in taken. taken holds names
// that are already in use (existing workspace and/or window names). Comparison
// is on the bare slug; callers pass the slug set, not the prefixed names.
func UniqueSlug(base string, taken map[string]bool) string {
	if !taken[base] {
		return base
	}
	for n := 2; ; n++ {
		candidate := base + "-" + itoa(n)
		if !taken[candidate] {
			return candidate
		}
	}
}

// itoa is a tiny base-10 formatter to avoid importing strconv for one call.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
