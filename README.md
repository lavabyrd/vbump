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

### Protocol-specific Version Bumping

For protocol-specific version bumping, use the `--protocol` flag. This will:

- Look for the VERSION file in `plugins/<protocol>/VERSION`
- Handle protocol-prefixed versions (e.g., `solana-0.3.1`)
- Pass both `VERSION` and `PROTOCOL` variables to make

```bash
vbump --protocol=solana  # e.g., solana-0.3.1 -> solana-0.3.2
```

The make command will receive:

```bash
VERSION=0.3.2 PROTOCOL=solana make bump
```

## Requirements

- A `VERSION` file containing a semantic version (X.Y.Z)
  - For protocol-specific versions, the file should be in `plugins/<protocol>/VERSION` and contain a version with protocol prefix (e.g., `solana-0.3.1`)
- A `make bump` target that accepts:
  - `VERSION` parameter for standard version bumping
  - `PROTOCOL` parameter for protocol-specific version bumping
- Git repository with no uncommitted changes
- Not on the main branch

## License

MIT
