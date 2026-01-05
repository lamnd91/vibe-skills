# Contributing to Vibe Skills

Thank you for your interest in contributing to Vibe Skills! This document provides guidelines and instructions for contributing.

## Ways to Contribute

- **Add new skills**: Create skills for technologies you know well
- **Improve existing skills**: Enhance documentation, fix errors
- **Report bugs**: Submit issues for problems you encounter
- **Suggest features**: Share ideas for improving the CLI
- **Improve documentation**: Help make docs clearer

## Creating a New Skill

### 1. Choose the Right Stack

Skills are organized by technology stack:

```
skills/
├── common/      # Universal skills (code review, commits, PRs)
├── dotnet/      # .NET ecosystem
├── database/    # Database technologies
├── frontend/    # Frontend frameworks
├── devops/      # DevOps and infrastructure
```

### 2. Create the Skill Directory

```bash
mkdir -p internal/registry/skills/<stack>/<skill-name>
```

### 3. Write SKILL.md

Create a `SKILL.md` file with your skill content:

```markdown
# Skill Name

Brief description of what this skill teaches Claude to do.

## Overview
Explain the purpose and context.

## Guidelines
- Specific instruction 1
- Specific instruction 2
- Specific instruction 3

## Examples
Show concrete examples of how to apply this skill.

## Best Practices
List recommended practices.

## Common Pitfalls
Describe what to avoid.
```

### 4. Skill Guidelines

- **Be concise**: Claude has limited context, keep skills focused
- **Be specific**: Provide concrete examples, not abstract concepts
- **Be actionable**: Include clear instructions Claude can follow
- **Use formatting**: Headers, lists, and code blocks improve readability
- **Include examples**: Show both good and bad patterns

### 5. Test Locally

```bash
# Build the CLI
go build -o vibe-skills ./cmd/vibe-skills

# List skills to verify yours appears
./vibe-skills list

# Install your skill
./vibe-skills install <skill-name>

# Verify the skill file
cat .claude/skills/<skill-name>.md
```

## Submitting Changes

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-new-skill`
3. Make your changes
4. Test locally
5. Commit with conventional commits: `git commit -m "feat(skills): add kubernetes deployment skill"`
6. Push to your fork
7. Open a pull request

### PR Requirements

- [ ] Skill follows the directory structure
- [ ] SKILL.md is well-formatted
- [ ] Skill appears in `vibe-skills list`
- [ ] Skill installs correctly
- [ ] No breaking changes to existing skills

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(skills): add new database optimization skill
fix(cli): resolve install path issue on Windows
docs(readme): update installation instructions
```

## Development Setup

### Prerequisites

- Go 1.22 or later
- Git

### Building

```bash
# Clone the repo
git clone https://github.com/cuongtl/vibe-skills.git
cd vibe-skills

# Install dependencies
go mod download

# Build
go build -o vibe-skills ./cmd/vibe-skills

# Run tests
go test ./...
```

### Project Structure

```
vibe-skills/
├── cmd/vibe-skills/       # CLI entry point
├── internal/
│   ├── cli/               # Cobra commands
│   ├── config/            # Config file handling
│   ├── installer/         # Skill installation logic
│   ├── registry/          # Embedded skills
│   │   └── skills/        # Skill content (embedded)
│   ├── updater/           # Self-update logic
│   └── version/           # Version info
├── scripts/               # Installation scripts
└── docs/                  # Documentation
```

## Code Style

- Follow standard Go formatting (`gofmt`)
- Run `golangci-lint` before submitting
- Keep functions small and focused
- Add comments for non-obvious logic

## Getting Help

- Open an issue for questions
- Join discussions in existing issues
- Tag maintainers if needed

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
