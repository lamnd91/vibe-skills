# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Build the CLI
go build -o vibe-skills ./cmd/vibe-skills

# Run tests
go test ./...

# Download dependencies
go mod tidy

# Run the built binary
./vibe-skills --help
./vibe-skills list
./vibe-skills install <skill-name>
```

## Architecture

This is a Go CLI tool that manages Claude Code skills. Skills are markdown files that get installed to `.claude/skills/` in user projects.

### Key Design Pattern: Embedded Skills

Skills are embedded into the binary at compile time using Go's `//go:embed` directive. The registry package (`internal/registry/`) contains the embedded skills filesystem:

```
internal/registry/skills/<stack>/<skill-name>/SKILL.md
```

When adding new skills, place them in `internal/registry/skills/` (not the top-level `skills/` directory).

### Package Structure

- **cmd/vibe-skills/main.go** - Entry point, calls `cli.Execute()`
- **internal/cli/** - Cobra commands (init, install, remove, list, search, version, self-update)
- **internal/registry/** - Skills registry with embedded files, skill lookup/search
- **internal/installer/** - Copies skills from registry to project's `.claude/skills/`
- **internal/config/** - Parses `.vibe-skills.yaml` config files
- **internal/updater/** - Self-update from GitHub releases
- **internal/version/** - Version info injected via ldflags

### Data Flow

1. User runs `vibe-skills install <skill>`
2. CLI creates `Registry` which loads embedded skills
3. `Installer` finds skill in registry and copies content to `.claude/skills/<skill>.md`

### Release Process

Uses GoReleaser triggered by git tags. Version info is injected via ldflags:
```bash
git tag v0.1.0
git push --tags
# GitHub Actions runs goreleaser
```
