// Package semverfmt normalizes messy semantic version strings into
// canonical MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD] form.
package semverfmt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrEmpty is returned when the input has no usable version content
// after trimming whitespace and a leading "v".
var ErrEmpty = errors.New("empty version string")

// Normalize rewrites a single version string into canonical form. It
// tolerates the messiness that shows up in the wild: surrounding
// whitespace, a leading "v" or "V", missing minor/patch components
// ("1.2" becomes "1.2.0"), and leading zeros in numeric fields
// ("01.02.3" becomes "1.2.3"). It does not guess its way around
// versions that are missing a numeric component entirely or that have
// more than three dotted components in the core.
func Normalize(input string) (string, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return "", ErrEmpty
	}
	if s[0] == 'v' || s[0] == 'V' {
		s = s[1:]
	}
	if s == "" {
		return "", ErrEmpty
	}

	var build string
	hasBuild := false
	if i := strings.IndexByte(s, '+'); i >= 0 {
		hasBuild = true
		build = s[i+1:]
		s = s[:i]
	}

	var prerelease string
	hasPrerelease := false
	if i := strings.IndexByte(s, '-'); i >= 0 {
		hasPrerelease = true
		prerelease = s[i+1:]
		s = s[:i]
	}

	core, err := normalizeCore(s)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString(core)

	if hasPrerelease {
		pre, err := normalizeIdentifiers(prerelease, true)
		if err != nil {
			return "", fmt.Errorf("prerelease: %w", err)
		}
		b.WriteByte('-')
		b.WriteString(pre)
	}

	if hasBuild {
		bld, err := normalizeIdentifiers(build, false)
		if err != nil {
			return "", fmt.Errorf("build metadata: %w", err)
		}
		b.WriteByte('+')
		b.WriteString(bld)
	}

	return b.String(), nil
}

func normalizeCore(s string) (string, error) {
	rawParts := strings.Split(s, ".")
	if len(rawParts) > 3 {
		return "", fmt.Errorf("too many numeric components in %q", s)
	}

	parts := [3]string{"0", "0", "0"}

	for i, p := range rawParts {
		p = strings.TrimSpace(p)
		if p == "" {
			return "", fmt.Errorf("empty numeric component in %q", s)
		}
		n, err := strconv.ParseUint(p, 10, 64)
		if err != nil {
			return "", fmt.Errorf("component %q is not a non-negative integer", p)
		}
		parts[i] = strconv.FormatUint(n, 10)
	}

	return strings.Join(parts[:], "."), nil
}

// normalizeIdentifiers validates and rewrites a dot-separated identifier
// list (prerelease or build metadata). Leading zeros are stripped from
// numeric identifiers only when stripLeadingZeros is set: the semver
// spec forbids them in prerelease identifiers but allows them in build
// metadata, since build metadata carries no ordering meaning.
func normalizeIdentifiers(s string, stripLeadingZeros bool) (string, error) {
	raw := strings.Split(s, ".")
	out := make([]string, len(raw))
	for i, id := range raw {
		id = strings.TrimSpace(id)
		if id == "" {
			return "", fmt.Errorf("empty identifier in %q", s)
		}
		for _, c := range id {
			if !isIdentifierChar(c) {
				return "", fmt.Errorf("invalid character %q in identifier %q", c, id)
			}
		}
		if stripLeadingZeros && isNumeric(id) && len(id) > 1 && id[0] == '0' {
			n, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				return "", fmt.Errorf("identifier %q is not a valid number", id)
			}
			id = strconv.FormatUint(n, 10)
		}
		out[i] = id
	}
	return strings.Join(out, "."), nil
}

func isIdentifierChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-'
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
