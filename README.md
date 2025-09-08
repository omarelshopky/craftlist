# CraftList

A configurable wordlist generator that creates targeted password lists based on organizational intelligence. Transforms seed words through variations and combines them with contextual elements for comprehensive security testing wordlists.


## Features

### Target-Specific Intelligence

- Custom word integration for company names, products, locations
- SSID-based patterns from wireless networks
- Smart pairing of company terms with common password elements

### Advanced Word Variations

- L33t speak transformations (a → @, e → 3, i → 1, o → 0, s → $)
- Multiple case variations
- Word segmentation and recombination
- Hybrid combinations of custom and common terms

### Flexible Pattern System

- Year integration with configurable ranges
- Sequential number patterns
- Multiple separator handling
- Complex multi-element patterns

### High Performance

- Concurrent processing
- Memory streaming for large wordlists
- Real-time filtering and deduplication

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
craftlist -w words.ls [-s ssids.ls] [-c config.json] [-o passwords.ls] [-max-length 8] [-max-length 64] [--max-year 2025] [--min-year 1990]
```

## Quick Start

Follow these steps to generate password lists using CraftList:

1. Copy the content from `examples/config.json` and modify it according to your specific requirements. The configuration file controls the number and types of passwords generated.

2. Create the following input files:

- `words.ls`: Contains company names, abbreviations, and other custom terms
- `ssids.ls`: Contains wireless network names (SSIDs)

3. Run the following command to generate your password list:

```bash
craftlist -c config.json -w words.ls -s ssids.ls -o passwords.ls
```

## Patterns

With these placeholders, you can create flexible password patterns like:

- `<CUSTOM><SEP><YEAR>`
- `<COMMON><SHORTYEAR><SEP><NUM>`

### Patterns Placeholders

CraftList supports the following placeholders in your password patterns:

- `<CUSTOM>`: Inserts custom word variations from the file specified with the --words flag
- `<COMMON>`: Inserts common word variations based on the list defined in your config file
- `<SSID>`: Inserts SSID variations from the file specified with the --ssids flag
- `<SEP>`: Inserts separators based on the list defined in your config file
- `<YEAR>`: Inserts full year based on the range defined in flags or config file (e.g., 2025)
- `<SHORTYEAR>`: Inserts two-digit year based on the range defined in flags or config file (e.g., 25)
- `<NUM>`: Inserts numbers based on the list defined in your config file

### Special Numeric Notation

You can use `d` characters to generate digit ranges:

- `ddd` generates all numbers from 000 to 999
- `5d` generates all numbers from 50 to 59
- `d` generates all single digits from 0 to 9

## Project Structure

```
craftlist/
├── cmd/craftlist/          # Application entry point
├── internal/               # Private application code
│   ├── config/             # Configuration management
│   ├── generator/          # Password generation logic
├── pkg/                    # Public packages
│   └── wordlist/           # Wordlist management
├── examples/               # Examples
│   └── config.json         # Config file sample
├── go.mod                  # Go module definition
├── Makefile                # Build automation
└── README.md               # Documentation
```

## Development

```bash
# Install dependencies
make deps

# Lint code
make lint

# Build for multiple platforms
make release
```