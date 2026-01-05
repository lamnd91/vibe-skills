---
name: code-reviewer
description: Systematic code review for quality, correctness, and maintainability. Use when reviewing pull requests, code changes, diffs, or when asked to review/critique code. Covers functionality, architecture, performance, security, testing, and documentation with structured feedback using priority prefixes.
---

# Code Reviewer

Systematic approach to reviewing code changes.

## Review Process

1. **Understand context** - Read PR description, linked issues, related files
2. **Review by area** - Apply relevant checklists below
3. **Provide feedback** - Use comment format with priority prefixes

## Comment Format

| Prefix | Meaning | Action |
|--------|---------|--------|
| `[BLOCKING]` | Must fix before merge | Required |
| `[SUGGESTION]` | Improvement opportunity | Optional |
| `[QUESTION]` | Need clarification | Response needed |
| `[NIT]` | Minor style issue | Optional |

**Structure:**
```
[PREFIX] Brief issue description

Why: Explanation of the problem or risk
Fix: Suggested solution or alternative
```

**Example:**
```
[BLOCKING] SQL injection vulnerability in user search

Why: User input concatenated directly into query string
Fix: Use parameterized query

// Before
var sql = $"SELECT * FROM Users WHERE Name = '{input}'";

// After  
var sql = "SELECT * FROM Users WHERE Name = @name";
cmd.Parameters.AddWithValue("@name", input);
```

## Review Checklists

### Functionality
- Code accomplishes intended purpose
- Edge cases handled (null, empty, boundary values)
- Error handling comprehensive and user-friendly
- No regressions to existing functionality

### Code Quality
- Readable and self-documenting
- Functions have single responsibility
- No code duplication (DRY)
- Clear, consistent naming

### Architecture
- Follows established project patterns
- Dependencies appropriate and minimal
- No circular dependencies
- Changes in correct location/layer

### Performance
- No N+1 queries
- Resources properly released/cleaned up
- Database queries optimized
- Async operations where beneficial

### Security
- Input validation present
- No hardcoded secrets/credentials
- SQL injection prevention (parameterized queries)
- XSS prevention (output encoding)
- Authentication/authorization checks in place

### Testing
- Unit tests cover new/changed code
- Edge cases tested
- Tests readable and maintainable
- Integration tests if needed

### Documentation
- Public APIs documented
- Complex logic has comments explaining "why"
- Breaking changes documented

### Database & API Changes
- Migrations reversible
- API backward compatible or versioned
- Indexes added for query patterns

## Feedback Principles

- Point to exact lines with specific alternatives
- Explain *why* something is problematic
- Focus on code, not the author
- Acknowledge good patterns when found