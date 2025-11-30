# CI/CD Pipeline Documentation

This directory contains GitHub Actions workflows for continuous integration and continuous deployment.

## Workflows

### CI Workflow (`ci.yml`)

The main CI pipeline that runs on every push and pull request to `main` and `develop` branches.

#### Jobs

1. **Test and Coverage**
   - Runs all unit tests with race detection
   - Generates coverage report
   - **Enforces minimum 90% code coverage**
   - Uploads coverage to Codecov
   - Uploads coverage report as artifact

2. **Lint**
   - Runs `golangci-lint` to check code quality
   - Enforces code style and best practices

3. **Build**
   - Compiles the application
   - Ensures the code builds successfully
   - Uploads binary as artifact

#### Coverage Requirements

The CI pipeline enforces a **minimum coverage threshold of 90%**. If coverage falls below this threshold, the pipeline will fail.

#### Triggers

- **Push**: Triggers on push to `main` or `develop` branches
- **Pull Request**: Triggers on PRs targeting `main` or `develop` branches

## Running Tests Locally

### Quick Test
```bash
make test
```

### Test with Coverage Check (90% threshold)
```bash
make coverage
```

### Generate HTML Coverage Report
```bash
make coverage-html
open coverage.html
```

### Show Coverage by Function
```bash
make coverage-func
```

### Manual Coverage Check
```bash
./scripts/coverage.sh
```

## Coverage Configuration

### Threshold
The coverage threshold is set to **90%** in:
- `.github/workflows/ci.yml` (line 63)
- `scripts/coverage.sh` (line 3)

### Excluded from Coverage
- Test files (`*_test.go`)
- Generated files (`*.gen.go`, `*.pb.go`)
- Vendor directory

## Linting

The project uses `golangci-lint` with the following enabled linters:
- bodyclose
- errcheck
- gosec
- gosimple
- govet
- ineffassign
- staticcheck
- stylecheck
- unused
- And more (see `.golangci.yml`)

### Run Linter Locally
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
make lint
```

## Codecov Integration

Coverage reports are automatically uploaded to Codecov on every CI run.

### Setup Codecov (First Time)

1. Go to [codecov.io](https://codecov.io)
2. Sign in with GitHub
3. Add your repository
4. Copy the upload token (optional, works without token for public repos)
5. Add `CODECOV_TOKEN` to GitHub Secrets if needed

### Codecov Badge

Add this to your README.md:
```markdown
[![codecov](https://codecov.io/gh/YOUR_USERNAME/usermes-backend/branch/main/graph/badge.svg)](https://codecov.io/gh/YOUR_USERNAME/usermes-backend)
```

## Artifacts

The CI pipeline uploads the following artifacts:

1. **coverage-report**: HTML coverage report
   - Available for 7 days
   - Can be downloaded from the Actions tab

2. **usermes-backend**: Compiled binary
   - Available for 7 days
   - Ready to deploy

## Troubleshooting

### Pipeline Fails with "Coverage Below Threshold"

If your PR fails because coverage is below 90%:

1. Run `make coverage` locally to see which packages need more tests
2. Add unit tests for uncovered code
3. Focus on:
   - Domain layer (entities, value objects)
   - Application layer (use cases)
   - Critical business logic

### Linter Errors

If the linter job fails:

1. Run `make lint` locally
2. Fix reported issues
3. Run `make fmt` to auto-format code
4. Commit and push changes

### Build Fails

If the build job fails:

1. Ensure code compiles locally: `go build ./...`
2. Check for missing dependencies: `go mod tidy`
3. Verify Go version matches CI (1.22.1)

## Best Practices

### Writing Tests

1. **Test the behavior, not the implementation**
   - Focus on what the code does, not how it does it

2. **Follow AAA pattern**
   - Arrange: Set up test data
   - Act: Execute the code being tested
   - Assert: Verify the results

3. **Use table-driven tests** for multiple scenarios:
   ```go
   tests := []struct {
       name    string
       input   string
       want    string
       wantErr bool
   }{
       {"valid input", "test", "TEST", false},
       {"empty input", "", "", true},
   }
   ```

4. **Mock external dependencies**
   - Use interfaces for testability
   - Mock repositories, external APIs, etc.

5. **Test edge cases**
   - Empty inputs
   - Nil values
   - Boundary conditions
   - Error scenarios

### Coverage Tips

- Aim for 90%+ coverage in:
  - Domain layer (entities, value objects)
  - Application layer (use cases)
  
- Don't obsess over 100% coverage in:
  - Infrastructure layer (HTTP handlers, DTOs)
  - Configuration files
  - Main function

- Focus on testing **business logic** over **plumbing code**

## Updating the Pipeline

### Changing Coverage Threshold

Update in both files:

1. `.github/workflows/ci.yml`:
   ```yaml
   THRESHOLD=90.0  # Change this value
   ```

2. `scripts/coverage.sh`:
   ```bash
   THRESHOLD=90.0  # Change this value
   ```

### Adding New Linters

Edit `.golangci.yml` and add to the `linters.enable` section.

### Adding New Test Steps

Edit `.github/workflows/ci.yml` and add steps under the `test` job.

## Performance

### CI Run Time

- **Average run time**: ~3-5 minutes
- **Test job**: ~2 minutes
- **Lint job**: ~1 minute
- **Build job**: ~30 seconds

### Optimization

To speed up CI:
- Use caching for Go modules (already enabled)
- Run jobs in parallel (already configured)
- Skip tests for documentation-only changes

## Security

### Secrets

The pipeline uses the following secrets:
- `CODECOV_TOKEN` (optional for public repos)

### Dependency Security

Consider adding:
```yaml
- name: Run gosec security scanner
  uses: securego/gosec@master
```

## Badges

Add these badges to your README.md:

```markdown
[![CI](https://github.com/YOUR_USERNAME/usermes-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/usermes-backend/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/YOUR_USERNAME/usermes-backend/branch/main/graph/badge.svg)](https://codecov.io/gh/YOUR_USERNAME/usermes-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/usermes-backend)](https://goreportcard.com/report/github.com/YOUR_USERNAME/usermes-backend)
```

## Support

For issues with the CI pipeline, please open an issue with:
- The workflow run URL
- Error message
- Steps to reproduce locally