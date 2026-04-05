# Contributing to lx

Thank you for your interest in contributing! This document provides guidelines for setting up your development environment and submitting contributions.

## Local Setup

### Prerequisites

- **Go 1.26.1** or later
- **golangci-lint** (linter)
- **Git**

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/andre-carbajal/lx.git
   cd lx
   ```

2. Install golangci-lint:
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

3. Verify setup:
   ```bash
   make build    # Should compile without errors
   make test     # Should run tests
   make lint     # Should pass linter
   ```

## Development Workflow

### Building

```bash
make build       # Compiles to ./lx (or ./lx.exe on Windows)
```

### Testing

```bash
make test        # Runs all tests with coverage report
```

### Linting

```bash
make lint        # Runs golangci-lint
```

### Cleaning

```bash
make clean       # Removes build artifacts and coverage.out
```

## Code Standards

- All pull requests must:
  - Pass `make lint` without warnings
  - Pass `make test` with ≥80% code coverage
  - Include unit tests for new functionality
  - Have clear commit messages (Conventional Commits recommended)

## Commit Message Format

We follow Conventional Commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Examples:
- `feat(detector): add Windows API process detection`
- `fix(parser): handle quoted arguments with spaces`
- `test(parser): add 25 test cases for command parsing`
- `docs(readme): update roadmap`

## Pull Request Process

1. Create a feature branch: `git checkout -b feat/my-feature`
2. Make your changes and add tests
3. Run `make test` and `make lint` locally
4. Commit with clear messages
5. Push to your fork and create a PR to `main`
6. CI will run automatically — all checks must pass

## Questions?

Open an issue on GitHub or contact the maintainers.
