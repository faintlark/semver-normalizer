// Command semverfmt reads version strings from stdin, one per line,
// and writes their normalized form to stdout.
package main

import (
	"fmt"
	"os"

	semverfmt "github.com/faintlark/semver-normalizer"
)

func main() {
	if err := semverfmt.Stream(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "semverfmt:", err)
		os.Exit(1)
	}
}
