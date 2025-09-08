# CraftList

A tool for generating customized wordlists tailored to a company's specific details.


## Features

- Target-specific password generation
- Multiple case variations
- L33t speak substitutions
- Year-based combinations
- Common word patterns
- Fast, concurrent generation

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
./bin/craftlist -w words.ls [-s ssids.ls] [-c config.json] [-max-length 8] [-max-length 64] [--max-year 2025] [--min-year 1990]
```


## Patterns

available placeholders: 
- `<CUSTOM>`: Custom word variation based on the file passed to `--words` flag
- `<COMMON>`: Common word variation based on the list defined in the config file
- `<SSID>`: SSID variation based on the file passed to `--ssids` flag
- `<SEP>`: Separator based on the list defined in the config file
- `<YEAR>`: Year based on the range defined using flags or config file. ex. 2025
- `<SHORTYEAR>`: Year based on the range defined using flags or config file. ex. 25
- `<NUM>`: Number based on the list defined in the config file. Note that the items can contains `d` to utilize all the digits between 0-9. ex. `ddd` generates num 000 to 999.


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