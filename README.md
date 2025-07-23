# vbump

A simple command-line tool to automate version bumping in repositories that use a `VERSION` file and `make bump` command.

## Installation

```bash
go install github.com/lavabyrd/vbump/cmd/vbump@latest
```

## Usage

By default, `vbump` will increment the patch version:

```bash
vbump  # e.g., 1.12.2 -> 1.12.3
```

You can also bump minor or major versions:

```bash
vbump --minor  # e.g., 1.12.2 -> 1.13.0
vbump --major  # e.g., 1.12.2 -> 2.0.0
```

## Requirements

- A `VERSION` file containing a semantic version (X.Y.Z)
- A `make bump` target that accepts:
  - `VERSION` parameter for standard version bumping
- Git repository with no uncommitted changes
- Not on the main branch

## Development

### Running Tests

```bash
go test ./...
```

Or to run tests for a specific package:

```bash
go test ./cmd/vbump
```

### Building

```bash
go build ./cmd/vbump
```

## License

MIT
