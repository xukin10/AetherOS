# Contributing to AetherOS

## Development Flow

Issue → Branch → Pull Request → Review → Merge → Release

**NEVER commit directly to main.**

## Git Branch Strategy

- main - Production
- develop - Development
- feature/* - Feature branches
- fix/* - Bug fixes
- release/* - Release branches

## Commit Convention

feat(scope): description
fix(scope): description
docs: description
refactor: description

## Pull Request Requirements

Each PR must:
- Have a single objective
- Compile successfully
- Include tests
- Update documentation

## Definition of Done

- [ ] Code compiles
- [ ] gofmt passes
- [ ] golangci-lint passes
- [ ] go test passes
- [ ] Documentation updated
- [ ] PR reviewed and approved