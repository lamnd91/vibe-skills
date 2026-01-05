# Pull Request

Best practices for creating clear and reviewable pull requests.

## PR Title

Use a clear, descriptive title that summarizes the change:

```
feat: Add user authentication with OAuth2
fix: Resolve race condition in order processing
refactor: Extract payment logic to separate service
```

## PR Description Template

```markdown
## Summary
Brief description of what this PR does and why.

## Changes
- List of specific changes made
- Another change
- Yet another change

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Refactoring (no functional changes)
- [ ] Documentation update

## Testing
Describe how you tested these changes:
- Unit tests added/updated
- Manual testing steps
- Test environment used

## Screenshots (if applicable)
Add screenshots for UI changes.

## Checklist
- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes

## Related Issues
Closes #123
Related to #456
```

## Best Practices

### Keep PRs Small
- Aim for < 400 lines of code changed
- Split large changes into multiple PRs
- One logical change per PR

### Make it Reviewable
- Provide context in the description
- Add inline comments for complex logic
- Include before/after screenshots for UI changes

### Respond to Feedback
- Address all comments before merging
- Explain your decisions if you disagree
- Push fixes as separate commits for easy re-review

### Before Requesting Review
1. Self-review your changes
2. Run tests locally
3. Rebase on latest main/master
4. Resolve any conflicts
5. Verify CI passes
