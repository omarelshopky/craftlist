# CraftList

A tool for generating customized wordlists tailored to a company's specific details.


## Features

- Target-specific password generation
- Multiple case variations
- L33t speak substitutions
- Year-based combinations
- Common word patterns
- Fast, concurrent generation
- Clean output formatting

## Installation

### From Source
```bash
git clone https://github.com/omarelshopky/craftlist
cd craftlist
make build
```

### Using Go
```bash
go install github.com/omarelshopky/craftlist/cmd/craftlist@latest
```

## Usage

```bash
# Build and run
make run

# Or run directly
./bin/craftlist
```

### Example Input
```
Company names: Target Corp, Target
Abbreviations: TC, TCLLC
SSIDs: TC-WiFi, Corporate-Net
```

## Project Structure

```
craftlist/
├── cmd/craftlist/          # Application entry point
├── internal/               # Private application code
│   ├── config/             # Configuration management
│   ├── generator/          # Password generation logic
│   └── output/             # Output handling
├── pkg/                    # Public packages
│   └── wordlist/           # Wordlist management
├── go.mod                  # Go module definition
├── Makefile                # Build automation
└── README.md               # Documentation
```

## Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Lint code
make lint

# Build for multiple platforms
make release
```