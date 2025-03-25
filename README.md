# vbump

A simple command-line tool to automate version bumping in repositories that use a `VERSION` file and `make bump` command.

## Installation

```bash
go install github.com/markpreston/vbump/cmd/vbump@latest
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

- A `VERSION` file in the current directory containing a semantic version (X.Y.Z)
- A `make bump` target that accepts a `VERSION` parameter
- Git repository with no uncommitted changes
- Not on the main branch

## License

MIT
