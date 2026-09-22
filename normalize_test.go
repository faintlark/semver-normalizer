package semverfmt

import (
	"errors"
	"testing"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"already canonical", "1.2.3", "1.2.3"},
		{"lowercase v prefix", "v1.2.3", "1.2.3"},
		{"uppercase v prefix", "V1.2.3", "1.2.3"},
		{"missing patch", "1.2", "1.2.0"},
		{"major only", "1", "1.0.0"},
		{"leading zeros in core", "01.02.3", "1.2.3"},
		{"surrounding whitespace", "  1.2.3  ", "1.2.3"},
		{"prerelease passthrough", "1.2.3-alpha", "1.2.3-alpha"},
		{"prerelease leading zero stripped", "1.2.3-alpha.01", "1.2.3-alpha.1"},
		{"prerelease bare zero kept", "1.2.3-0", "1.2.3-0"},
		{"prerelease zero run collapses", "1.2.3-00", "1.2.3-0"},
		{"build metadata keeps leading zeros", "1.2.3+001", "1.2.3+001"},
		{"prerelease and build combined", "V2.0.0-RC.01+build.007", "2.0.0-RC.1+build.007"},
		{"prerelease with build", "1.2.3-alpha+build", "1.2.3-alpha+build"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Normalize(tc.input)
			if err != nil {
				t.Fatalf("Normalize(%q) returned error: %v", tc.input, err)
			}
			if got != tc.want {
				t.Errorf("Normalize(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestNormalizeEmpty(t *testing.T) {
	cases := []string{"", "   ", "v", "V", "\t\n"}

	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			_, err := Normalize(input)
			if !errors.Is(err, ErrEmpty) {
				t.Errorf("Normalize(%q) error = %v, want ErrEmpty", input, err)
			}
		})
	}
}

func TestNormalizeErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"too many core components", "1.2.3.4"},
		{"non-numeric component", "1.a.3"},
		{"empty numeric component", "1..3"},
		{"trailing dot in core", "1.2.3."},
		{"invalid character in prerelease", "1.2.3-alpha_beta"},
		{"empty prerelease identifier", "1.2.3-alpha..beta"},
		{"trailing dash with nothing after it", "1.2.3-"},
		{"invalid character in build metadata", "1.2.3+build_1"},
		{"trailing plus with nothing after it", "1.2.3+"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Normalize(tc.input)
			if err == nil {
				t.Errorf("Normalize(%q) = %q, want error", tc.input, got)
			}
		})
	}
}
