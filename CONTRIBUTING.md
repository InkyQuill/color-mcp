# Contributing to Color MCP

Thank you for your interest in contributing to Color MCP!

## Development Setup

1. Fork and clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/color-mcp.git
   cd color-mcp
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

## Commit Convention

We use [Conventional Commits](https://www.conventionalcommits.org/) for our commit messages:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that don't affect code meaning (formatting, etc.)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvement
- **test**: Adding or updating tests
- **build**: Changes to build system or dependencies
- **ci**: CI/CD changes
- **chore**: Other changes that don't modify src or test files
- **revert**: Revert a previous commit

### Scopes

- **convert**: Color conversion functionality
- **compare**: Color comparison functionality
- **detect**: Format detection functionality
- **ci**: CI/CD changes
- **docs**: Documentation changes

### Examples

```
feat(compare): add WCAG contrast ratio calculation

fix: correct OKLCH to RGB conversion for edge cases

docs: update README with new comparison feature

test: add coverage for hue difference calculation
```

## Pull Request Process

1. Create a branch from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. Make your changes and write tests

3. Ensure all tests pass:
   ```bash
   go test -v -race -cover ./...
   ```

4. Format your code:
   ```bash
   gofmt -w .
   ```

5. Commit using conventional commits:
   ```bash
   git commit -m "feat(scope): description"
   ```

6. Push and create a pull request

## Testing

Run the full test suite:
```bash
go test -v -race -cover ./...
```

Run specific test:
```bash
go test -v ./internal -run TestCompareColors
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write table-driven tests where appropriate
- Add comments for exported functions
- Aim for >80% test coverage
