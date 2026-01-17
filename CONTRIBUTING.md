# Contributing to ibp-geodns-agent

Thank you for your interest in contributing to ibp-geodns-agent!

## Development Setup

1. Clone the repository
2. Ensure you have Go 1.21 or later installed
3. Install dependencies: `go mod download`
4. Build: `make build`

## Code Style

- Follow Go standard formatting (`go fmt`)
- Run linters: `make lint`
- Write tests for new functionality
- Update documentation as needed

## Pull Request Process

1. Create a feature branch from `main`
2. Make your changes
3. Ensure all tests pass: `make test`
4. Update CHANGELOG.md if applicable
5. Submit a pull request with a clear description

## Project Structure

This project follows the same structural patterns as other IBP GeoDNS repositories:

- `src/cmd/agent/` - Main entry point
- `src/` - Source code packages
- `config/` - Configuration files and examples
- `docs/` - Documentation
- `scripts/` - Utility scripts

## Integration with ibp-geodns-libs

This project uses `ibp-geodns-libs` for common functionality. When adding features, prefer using the libs where possible rather than reimplementing.
