package semverfmt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Stream reads version strings from r, one per line, normalizes each,
// and writes the result to w. Blank lines and line endings are passed
// through unchanged so the output keeps the same shape as the input.
//
// It reads with bufio.Reader.ReadString('\n') instead of bufio.Scanner
// on purpose: a Scanner holds its token size limit fixed up front and
// errors out on a single very long line, while ReadString just grows
// the buffer for whatever line it is currently on and drops it before
// moving to the next. Either way only one line is ever held in memory,
// so a file with a million version tags costs the same few hundred
// bytes as a file with ten.
func Stream(r io.Reader, w io.Writer) error {
	reader := bufio.NewReader(r)
	writer := bufio.NewWriter(w)

	lineNum := 0
	for {
		lineNum++
		line, err := reader.ReadString('\n')

		if len(line) > 0 {
			endsWithNewline := strings.HasSuffix(line, "\n")
			content := strings.TrimRight(line, "\r\n")

			out := content
			if strings.TrimSpace(content) != "" {
				normalized, nerr := Normalize(content)
				if nerr != nil {
					return fmt.Errorf("line %d: %w", lineNum, nerr)
				}
				out = normalized
			}

			if _, werr := writer.WriteString(out); werr != nil {
				return werr
			}
			if endsWithNewline {
				if werr := writer.WriteByte('\n'); werr != nil {
					return werr
				}
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
