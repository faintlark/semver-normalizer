# semver-normalizer

Version strings collected from real sources are rarely clean semver.
Git tags have a `v` prefix. Old changelog entries drop the patch
number. Some build system somewhere zero-pads the minor version.
`semver-normalizer` takes that mess and rewrites it into canonical
`MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]` form:

```
v1.2      -> 1.2.0
 1.02.3   -> 1.2.3
V2.0.0-RC.01+build.007 -> 2.0.0-RC.1+build.007
```

(Build metadata identifiers keep leading zeros — semver says build
metadata carries no ordering meaning, so there's nothing to normalize
away there. Prerelease identifiers do get leading zeros stripped,
since semver forbids them.)

## Command line

```
go build -o semverfmt ./cmd/semverfmt
printf 'v1.2\n 1.02.3 \nV2.0.0-RC.01+build.007\n' | ./semverfmt
```

```
1.2.0
1.2.3
2.0.0-RC.1+build.007
```

Blank lines and lines that fail to parse are handled per line: a bad
line reports its line number and stops the run, everything already
written to stdout stays written.

## As a library

```go
import semverfmt "github.com/faintlark/semver-normalizer"

out, err := semverfmt.Normalize("v1.2")
// out == "1.2.0"
```

For anything larger than a single string, `Stream` reads from an
`io.Reader` and writes to an `io.Writer`, one version per line. It's
built around `bufio.Reader.ReadString`, not `bufio.Scanner`, so it
never has to hold more than the current line in memory — a file of a
hundred version tags and a file of a hundred million cost the same
amount of memory to process:

```go
err := semverfmt.Stream(os.Stdin, os.Stdout)
```

## Status

Early. The core normalization rules above are implemented and tested
by hand; see the roadmap in commit history for what's still missing
(a real test suite, comparison/sorting, a `--strict` mode that refuses
to guess at malformed input instead of erroring).

## License

MIT, see LICENSE.
