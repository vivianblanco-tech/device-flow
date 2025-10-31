# Contributing to Laptop Tracking System

Thank you for your interest in contributing to the Laptop Tracking System! This document provides guidelines and instructions for contributing.

## Development Process

We follow a structured development process based on Test-Driven Development (TDD) and the plan outlined in `plan.md`.

### Workflow

1. **Check the Plan**: Review `plan.md` for the current phase and tasks
2. **Create a Branch**: Create a feature branch from `develop`
3. **Write Tests First**: Follow TDD principles (Red → Green → Refactor)
4. **Implement Feature**: Write the minimum code to pass tests
5. **Document**: Add comments and update documentation
6. **Test**: Ensure all tests pass
7. **Commit**: Use conventional commit messages
8. **Push & PR**: Push to your branch and create a pull request

## Getting Started

1. **Fork the Repository**
   ```bash
   git clone https://github.com/yourusername/laptop-tracking-system.git
   cd laptop-tracking-system
   ```

2. **Set Up Development Environment**
   ```bash
   make dev-setup
   ```
   See `docs/SETUP.md` for detailed setup instructions.

3. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for code formatting
- Run `go vet` to catch common errors
- Use meaningful variable and function names
- Keep functions small and focused (< 50 lines preferred)
- Add comments for exported functions and complex logic

### File Organization

- Keep files under 300 lines of code
- Split large files into multiple smaller, focused files
- Use clear, descriptive file names
- Group related functionality together

### Project Structure

```
internal/
├── models/       - Data models and business logic
├── handlers/     - HTTP request handlers (thin layer)
├── middleware/   - HTTP middleware
├── auth/         - Authentication logic
├── validator/    - Input validation
├── email/        - Email service
└── jira/         - JIRA integration
```

## Test-Driven Development (TDD)

We strictly follow TDD principles:

### The TDD Cycle

1. **🟥 RED**: Write a failing test
   ```go
   func TestUserValidation(t *testing.T) {
       user := &User{Email: "invalid"}
       err := user.Validate()
       if err == nil {
           t.Error("Expected validation error for invalid email")
       }
   }
   ```

2. **🟩 GREEN**: Write minimum code to pass
   ```go
   func (u *User) Validate() error {
       if !strings.Contains(u.Email, "@") {
           return errors.New("invalid email")
       }
       return nil
   }
   ```

3. **🔄 REFACTOR**: Improve the code
   ```go
   func (u *User) Validate() error {
       if err := validateEmail(u.Email); err != nil {
           return fmt.Errorf("email validation failed: %w", err)
       }
       return nil
   }
   ```

### Testing Guidelines

- Write tests before implementation
- Aim for 80%+ code coverage
- Test edge cases and error conditions
- Use table-driven tests for multiple scenarios
- Mock external dependencies (database, APIs, email)
- Keep tests fast and independent

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package tests
go test ./internal/models -v

# Run with race detection
go test -race ./...
```

## Commit Message Convention

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation only changes
- `style:` Code style changes (formatting, semicolons, etc.)
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `chore:` Build process or auxiliary tool changes
- `perf:` Performance improvements
- `ci:` CI configuration changes

### Examples

```bash
feat(auth): add Google OAuth authentication

Implement Google OAuth 2.0 flow with domain restriction to @bairesdev.com.
Includes session management and user creation/lookup.

Closes #42
```

```bash
fix(database): correct connection pool settings

Increase max open connections to 25 to handle concurrent requests better.
```

```bash
test(models): add user validation tests

Add comprehensive tests for user model validation including email format,
password requirements, and role validation.
```

## Pull Request Process

1. **Update Documentation**: Ensure README and docs are updated
2. **Add Tests**: All new features must have tests
3. **Pass CI**: Ensure all CI checks pass
4. **Code Review**: Request review from maintainers
5. **Address Feedback**: Make requested changes
6. **Squash Commits**: Keep PR history clean

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe testing performed

## Checklist
- [ ] Tests pass locally
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] No new warnings
- [ ] Added tests for new features
```

## Code Review Guidelines

### For Authors

- Keep PRs small and focused
- Write clear PR descriptions
- Respond to feedback promptly
- Be open to suggestions

### For Reviewers

- Be respectful and constructive
- Focus on code, not the person
- Explain the "why" behind suggestions
- Approve when ready, request changes when needed

## Database Migrations

### Creating Migrations

```bash
make migrate-create name=add_users_table
```

This creates two files:
- `NNNN_add_users_table.up.sql`
- `NNNN_add_users_table.down.sql`

### Migration Guidelines

- Make migrations idempotent
- Always provide `up` and `down` migrations
- Test rollback scenarios
- Don't modify existing migrations
- Add comments explaining complex changes

### Example Migration

```sql
-- 000002_add_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    role user_role NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- 000002_add_users_table.down.sql
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

## Documentation

### Code Comments

- Add package-level comments
- Document all exported functions
- Explain complex algorithms
- Avoid obvious comments

```go
// Package auth provides authentication and authorization functionality.
// It supports both password-based and OAuth authentication.
package auth

// Authenticate verifies user credentials and returns a session token.
// It returns an error if credentials are invalid or the user is locked out.
func Authenticate(email, password string) (*Session, error) {
    // Implementation...
}
```

### Documentation Updates

Update relevant docs when making changes:
- `README.md` - General information
- `docs/SETUP.md` - Setup instructions
- `CONTRIBUTING.md` - This file
- API documentation (if applicable)

## Common Tasks

### Adding a New Feature

1. Check `plan.md` for the feature's phase
2. Create feature branch: `git checkout -b feature/feature-name`
3. Write tests first (TDD)
4. Implement feature
5. Update documentation
6. Create PR

### Fixing a Bug

1. Create bug branch: `git checkout -b fix/bug-description`
2. Write test that reproduces bug
3. Fix the bug
4. Ensure test passes
5. Create PR

### Adding a New Model

1. Write model tests in `internal/models/`
2. Implement model
3. Create database migration
4. Add validation logic
5. Document model fields

## Questions?

- Check existing documentation
- Review `plan.md` for project structure
- Ask in pull request comments
- Contact the development team

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing! 🎉

