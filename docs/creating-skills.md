# Creating Skills

This guide explains how to create new skills for Vibe Skills.

## What is a Skill?

A skill is a markdown document that instructs Claude Code on how to handle specific tasks or follow certain practices. When installed, skills are placed in `.claude/skills/` directory and automatically loaded by Claude Code.

## Skill Structure

Each skill is a `SKILL.md` file located at:
```
internal/registry/skills/<stack>/<skill-name>/SKILL.md
```

## Writing Effective Skills

### 1. Clear Purpose

Start with a clear description:

```markdown
# SQL Optimization

Optimize SQL queries for better performance and reduced resource usage.
```

### 2. Actionable Instructions

Provide specific, actionable guidance:

```markdown
## Query Optimization Rules

1. Always use parameterized queries
2. Avoid SELECT * - specify columns explicitly
3. Use appropriate indexes for WHERE clauses
4. Limit result sets with TOP/LIMIT
```

### 3. Examples

Include concrete examples:

```markdown
## Examples

### Bad
```sql
SELECT * FROM users WHERE name LIKE '%john%'
```

### Good
```sql
SELECT id, name, email
FROM users
WHERE name LIKE @searchTerm
```
```

### 4. Context-Aware

Consider when the skill applies:

```markdown
## When to Apply

- Writing new database queries
- Reviewing existing queries
- Debugging slow endpoints
```

## Stacks

Choose the appropriate stack for your skill:

| Stack | Description | Examples |
|-------|-------------|----------|
| `common` | Universal practices | code review, commits, PRs |
| `dotnet` | .NET ecosystem | EF Core, ASP.NET, C# |
| `database` | Databases | SQL, indexing, optimization |
| `frontend` | Frontend | React, Vue, CSS |
| `devops` | Operations | Docker, K8s, CI/CD |

## Template

```markdown
# Skill Name

Brief one-line description.

## Overview

2-3 sentences explaining the purpose and value.

## Guidelines

- Guideline 1
- Guideline 2
- Guideline 3

## Examples

### Good Example
```code
// example
```

### Bad Example
```code
// anti-pattern
```

## Best Practices

1. Best practice 1
2. Best practice 2

## Common Mistakes

- Mistake to avoid 1
- Mistake to avoid 2
```

## Testing Your Skill

1. Add your skill to the appropriate directory
2. Build the CLI: `go build -o vibe-skills ./cmd/vibe-skills`
3. Verify it appears: `./vibe-skills list`
4. Install it: `./vibe-skills install <skill-name>`
5. Check the output: `cat .claude/skills/<skill-name>.md`

## Submission

See [CONTRIBUTING.md](../CONTRIBUTING.md) for how to submit your skill.
