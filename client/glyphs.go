package client

const (
	unicodeRule     = "─" // ─ BOX DRAWINGS LIGHT HORIZONTAL
	unicodeEllipsis = "…" // … HORIZONTAL ELLIPSIS
	asciiRule       = "-"
	asciiEllipsis   = "..."
)

// ResolveBoxGlyphs returns the rule and ellipsis this console can carry.
func ResolveBoxGlyphs() (rule, ellipsis string) {
	if consoleIsUnicode() {
		return unicodeRule, unicodeEllipsis
	}
	return asciiRule, asciiEllipsis
}
