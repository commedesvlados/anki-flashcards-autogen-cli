<!--
Split from original README.md, 2024-06-09
-->

# Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## Development Setup

```bash
# Clone repository
git clone https://github.com/commedesvlados/anki-flashcards-autogen-cli.git
cd anki-flashcards-autogen-cli

# Install dependencies
go mod tidy
pip install genanki

# Build and install locally
make build
make install-user

# Run tests
make test
``` 