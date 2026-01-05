# listentotaxman-cli

A command-line interface for calculating UK tax and national insurance using the [listentotaxman.com](https://listentotaxman.com/) API.

## Installation

### Homebrew (macOS/Linux)

The easiest way to install on macOS or Linux:

```bash
brew install mheap/tap/listentotaxman
```

To upgrade:
```bash
brew upgrade mheap/tap/listentotaxman
```

### Docker

Run using Docker without installing:

```bash
# Pull the latest image
docker pull mheap/listentotaxman:latest

# Run a command
docker run --rm mheap/listentotaxman:latest check --income 100000

# With config file support (mount your local config)
docker run --rm \
  -v ~/.config/listentotaxman:/home/appuser/.config/listentotaxman \
  mheap/listentotaxman:latest check --income 100000

# Create an alias for convenience
alias listentotaxman='docker run --rm -v ~/.config/listentotaxman:/home/appuser/.config/listentotaxman mheap/listentotaxman:latest'
```

**Available tags:**
- `latest` - Latest stable release
- `0.1.0` - Specific version
- `0.1` - Latest patch version of 0.1.x

### Pre-built Binaries

Download pre-built binaries from [GitHub Releases](https://github.com/mheap/listentotaxman-cli/releases):

**Linux (x86_64):**
```bash
curl -L https://github.com/mheap/listentotaxman-cli/releases/latest/download/listentotaxman_linux_amd64.tar.gz | tar xz
sudo mv listentotaxman /usr/local/bin/
```

**Linux (ARM64):**
```bash
curl -L https://github.com/mheap/listentotaxman-cli/releases/latest/download/listentotaxman_linux_arm64.tar.gz | tar xz
sudo mv listentotaxman /usr/local/bin/
```

**macOS (Intel):**
```bash
curl -L https://github.com/mheap/listentotaxman-cli/releases/latest/download/listentotaxman_darwin_amd64.tar.gz | tar xz
sudo mv listentotaxman /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
curl -L https://github.com/mheap/listentotaxman-cli/releases/latest/download/listentotaxman_darwin_arm64.tar.gz | tar xz
sudo mv listentotaxman /usr/local/bin/
```

**Windows (x86_64):**
1. Download `listentotaxman_windows_amd64.zip` from [releases](https://github.com/mheap/listentotaxman-cli/releases)
2. Extract the archive
3. Add the directory to your PATH

**Windows (ARM64):**
1. Download `listentotaxman_windows_arm64.zip` from [releases](https://github.com/mheap/listentotaxman-cli/releases)
2. Extract the archive
3. Add the directory to your PATH

### Go Install

If you have Go 1.22+ installed:

```bash
go install github.com/mheap/listentotaxman-cli@latest
```

### Building from Source

```bash
git clone https://github.com/mheap/listentotaxman-cli
cd listentotaxman-cli
go build -o listentotaxman
```

### Shell Completions

After installation via Homebrew, completions are automatically installed.

For manual installations, generate completions:

```bash
# Bash
listentotaxman completion bash > /etc/bash_completion.d/listentotaxman

# Zsh
listentotaxman completion zsh > "${fpath[1]}/_listentotaxman"

# Fish
listentotaxman completion fish > ~/.config/fish/completions/listentotaxman.fish

# PowerShell
listentotaxman completion powershell > listentotaxman.ps1
```

### Verify Installation

```bash
listentotaxman version
```

## Usage

### Basic Example

Calculate tax for a salary:

```bash
listentotaxman check --income 100000
```

### Full Example

Calculate tax with all parameters:

```bash
listentotaxman check \
  --year 2025 \
  --region uk \
  --age 0 \
  --pension 3% \
  --income 100000 \
  --student-loan postgraduate \
  --extra 999 \
  --tax-code K12 \
  --married \
  --partner-income 25000
```

### Available Commands

#### `check` - Calculate Tax

Calculate UK tax and national insurance for a given salary.

**Flags:**

- `--income` (required) - Gross annual salary in pounds
- `--year` - Tax year (defaults to current tax year, based on April 5th cutoff)
- `--region` - Tax region (default: "uk", alias: "england")
  - Options: `uk`, `england` (alias for uk), `scotland`, `wales`, `ni`
- `--age` - Age for age-related calculations (default: "0")
- `--pension` - Pension contribution (e.g., "3%" or "3000")
- `--student-loan` - Student loan plan
  - Options: `plan1`, `plan2`, `plan4`, `postgraduate`, `scottish`
- `--extra` - Extra income or deductions
- `--tax-code` - Tax code (e.g., "1257L", "K12")
- `--married` - Married status (enables marriage allowance calculations)
- `--blind` - Blind person's allowance
- `--no-ni` - Exempt from National Insurance (e.g., working past state pension age)
- `--partner-income` - Partner's gross wage (requires `--married` flag)
- `--period` - Display period: yearly, monthly, weekly, daily, or hourly (default: "yearly")
- `--json` - Output as JSON instead of formatted table
- `--verbose` - Show detailed breakdown of tax calculation

**Examples:**

Basic calculation:

```bash
listentotaxman check --income 100000
```

With marriage allowance:

```bash
listentotaxman check --income 100000 --married --partner-income 25000
```

With blind person's allowance:

```bash
listentotaxman check --income 80000 --blind
```

Exempt from NI (e.g., past state pension age):

```bash
listentotaxman check --income 50000 --no-ni
```

#### `compare` - Compare Multiple Scenarios

Compare tax calculations across different job offers, salary levels, pension contributions, or tax years side-by-side.

**Usage:**

```bash
listentotaxman compare \
  --option "Label 1" --income AMOUNT [--pension X% --year YYYY ...] \
  --option "Label 2" --income AMOUNT [--pension X% --year YYYY ...] \
  [--period PERIOD] [--json] [--verbose]
```

**Per-Option Flags:**

Each `--option` group supports all flags from the `check` command:

- `--income` (required) - Gross annual salary in pounds
- `--year` - Tax year (defaults to current tax year)
- `--region` - Tax region (default: "uk", alias: "england")
- `--age` - Age for age-related calculations (default: "0")
- `--pension` - Pension contribution (e.g., "3%" or "3000")
- `--student-loan` - Student loan plan (plan1, plan2, plan4, postgraduate, scottish)
- `--extra` - Extra income or deductions
- `--tax-code` - Tax code (e.g., "1257L", "K12")
- `--married` - Married status
- `--blind` - Blind person's allowance
- `--no-ni` - Exempt from National Insurance
- `--partner-income` - Partner's gross wage (requires `--married`)

**Global Flags:**

These apply to all options in the comparison:

- `--period` - Display period: yearly, monthly, weekly, daily, or hourly (default: "yearly")
- `--json` - Output as JSON comparison object
- `--verbose` - Show detailed breakdown including tax brackets

**Requirements:**

- Minimum 2 options required
- Maximum 4 options supported
- Each option must have a unique label and `--income`

**Examples:**

Compare two job offers:

```bash
listentotaxman compare \
  --option "Current Job" --income 100000 --pension 3% \
  --option "New Offer" --income 120000 --pension 5%
```

Compare monthly take-home across salary levels:

```bash
listentotaxman compare \
  --period monthly \
  --option "Low" --income 80000 \
  --option "Mid" --income 100000 \
  --option "High" --income 120000
```

Compare tax years (to see rate changes):

```bash
listentotaxman compare \
  --option "2024" --income 100000 --year 2024 \
  --option "2025" --income 100000 --year 2025
```

Compare regions:

```bash
listentotaxman compare \
  --option "England" --income 100000 --region uk \
  --option "Scotland" --income 100000 --region scotland
```

Compare marriage allowance impact:

```bash
listentotaxman compare \
  --option "Single" --income 100000 \
  --option "Married" --income 100000 --married --partner-income 25000
```

Detailed comparison:

```bash
listentotaxman compare \
  --verbose \
  --option "Job 1" --income 100000 \
  --option "Job 2" --income 120000
```

JSON output for scripting:

```bash
listentotaxman compare \
  --json \
  --option "Job 1" --income 100000 \
  --option "Job 2" --income 120000 | jq '.comparison.net_pay'
```

**Output Format:**

The default output is a side-by-side comparison table. Status indicators (M=Married, B=Blind, NI=NI Exempt) are shown when applicable:

```
╔══════════════════════╦══════════════╦══════════════╗
║ Field                ║ Job 1        ║ Job 2        ║
║ Status               ║              ║ M            ║
╠══════════════════════╬══════════════╬══════════════╣
║ Gross Salary         ║   £100000.00 ║   £120000.00 ║
║ Tax Paid             ║    £27428.40 ║    £39428.40 ║
║ National Insurance   ║     £4010.60 ║     £4410.60 ║
║ Student Loan         ║        £0.00 ║        £0.00 ║
║ Pension (You)        ║        £0.00 ║        £0.00 ║
║ Net Pay              ║    £68561.00 ║    £76161.00 ║
╠══════════════════════╬══════════════╬══════════════╣
║ Employer's NI        ║    £14250.00 ║    £17250.00 ║
║ Pension (HMRC)       ║        £0.00 ║        £0.00 ║
║ Total Cost           ║   £114250.00 ║   £137250.00 ║
╚══════════════════════╩══════════════╩══════════════╝
```

With `--json`, outputs a comparison object:

```json
{
  "period": "yearly",
  "comparison": {
    "gross_pay": {
      "Job 1": 100000,
      "Job 2": 120000
    },
    "net_pay": {
      "Job 1": 68561,
      "Job 2": 76161
    },
    ...
  },
  "metadata": {
    "Job 1": {
      "tax_year": 2025,
      "tax_region": "uk",
      "tax_code": "1257L"
    },
    ...
  }
}
```

**Use Cases:**

- Compare multiple job offers to see which has better take-home pay
- Evaluate different pension contribution levels
- Compare tax burden across UK regions (England, Scotland, Wales, Northern Ireland)
- Analyze impact of student loan repayments on different salaries
- Compare tax years to see how rate changes affect your take-home
- Understand monthly/weekly income differences between salary options

#### `version` - Show Version

Display the CLI version information:

```bash
listentotaxman version
```

## Output Formats

### Default (Summary Table)

```
╔══════════════════════════════════════════════╗
║ Tax Calculation for 2025 (UK) - Yearly       ║
╠══════════════════════════════════════════════╣
║ Gross Salary                     £100000.00 ║
║ Taxable Pay                       £84421.00 ║
║ Tax Paid                          £26228.40 ║
║ National Insurance                 £4010.60 ║
║ Pension (You)                      £3000.00 ║
║ Net Pay                           £66761.00 ║
╠══════════════════════════════════════════════╣
║ Employer's NI                     £14250.00 ║
║ Pension (HMRC)                     £2000.00 ║
║ Total Cost                       £116250.00 ║
╚══════════════════════════════════════════════╝
```

### Verbose Output

Use the `--verbose` flag for a detailed breakdown:

```bash
listentotaxman check --income 100000 --pension 3% --verbose
```

```
Tax Year: 2025 (UK) - Yearly

Income:
  Gross Salary:             £100000.00
  Tax Free Allowance:        £12579.00
  Taxable Pay:               £84421.00

Tax Breakdown:
  Basic Rate (20%):           £7540.00
  Higher Rate (40%):         £18688.40
  Total Tax:                 £26228.40

Deductions:
  National Insurance:         £4010.60
  Pension (You):              £3000.00
  Total Deductions:          £33239.00

Net Pay:                     £66761.00

Employer Costs:
  Employer's NI:             £14250.00
  Pension (HMRC):             £2000.00
  Total Cost:               £116250.00
```

When status flags are active, they appear in the output:

```bash
listentotaxman check --income 100000 --married --blind --verbose
```

```
Tax Year: 2025 (UK) - Yearly
Status: Married • Blind Allowance

Income:
  ...
```

### JSON Output

Use the `--json` flag for machine-readable output:

```bash
listentotaxman check --income 100000 --pension 3% --json
```

## Time Periods

You can view tax calculations in different time periods using the `--period` flag. This divides all yearly values by the appropriate divisor, making it easy to understand your take-home pay on a monthly, weekly, daily, or hourly basis.

### Period Options

- `yearly` (default) - Annual figures (no division)
- `monthly` - Divide by 12
- `weekly` - Divide by 52
- `daily` - Divide by 365
- `hourly` - Divide by 2080 (52 weeks × 40 hours)

### Examples

**Monthly breakdown:**

```bash
listentotaxman check --income 100000 --period monthly
```

Output:

```
╔══════════════════════════════════════════════╗
║ Tax Calculation for 2025 (UK) - Monthly      ║
╠══════════════════════════════════════════════╣
║ Gross Salary                       £8333.33 ║
║ Taxable Pay                        £7035.08 ║
║ Tax Paid                           £2185.70 ║
║ National Insurance                  £334.22 ║
║ Pension (You)                       £250.00 ║
║ Net Pay                            £5563.42 ║
╠══════════════════════════════════════════════╣
║ Employer's NI                      £1187.50 ║
║ Pension (HMRC)                      £166.67 ║
║ Total Cost                         £9687.50 ║
╚══════════════════════════════════════════════╝
```

**Hourly rate (for contractors):**

```bash
listentotaxman check --income 75000 --period hourly
```

This shows your effective hourly rate after all deductions (divides by 2080 hours).

**Weekly take-home:**

```bash
listentotaxman check --income 52000 --period weekly
```

Shows weekly income and deductions.

**Daily rate:**

```bash
listentotaxman check --income 36500 --period daily
```

### Use Cases

- **Contractors & Freelancers**: Use `--period hourly` or `--period daily` to see your effective rate after tax
- **Monthly Budgeting**: Use `--period monthly` to understand your monthly take-home pay
- **Weekly Planning**: Use `--period weekly` to calculate weekly disposable income
- **Rate Negotiation**: Compare different salary offers on an hourly basis

### Period with Other Flags

The period flag works seamlessly with all other options:

```bash
# Detailed monthly breakdown
listentotaxman check --income 80000 --pension 5% --period monthly --verbose

# Hourly rate with student loan
listentotaxman check --income 100000 --student-loan postgraduate --period hourly

# JSON output with weekly values
listentotaxman check --income 52000 --period weekly --json | jq '.net_pay'
```

## Configuration File

You can create a configuration file at `~/.config/listentotaxman/config.yaml` to set default values:

```yaml
defaults:
  income: 0 # Set a default income if desired
  region: uk
  age: "0"
  pension: "5%"
  student-loan: ""
  tax-code: ""
  extra: 0
  year: "" # Leave empty to use smart default based on current date
  period: yearly # Options: yearly, monthly, weekly, daily, hourly
  married: false
  blind: false
  no-ni: false
  partner-income: 0
```

**Configuration Precedence:**

CLI flags override config file values, which override built-in defaults:

```
CLI Flags > Config File > Built-in Defaults
```

### Creating a Config File

```bash
mkdir -p ~/.config/listentotaxman
cat > ~/.config/listentotaxman/config.yaml << EOF
defaults:
  region: uk
  pension: "5%"
  age: "0"
  married: false
  blind: false
EOF
```

## Year Default Logic

If no year is specified (via flag or config), the CLI uses smart defaults:

- **After April 5th**: Uses current year (e.g., on April 6th, 2026 → uses 2026)
- **Before April 5th**: Uses previous year (e.g., on January 5th, 2026 → uses 2025)

This matches the UK tax year, which runs from April 6th to April 5th.

## Examples

### Simple Salary Check

```bash
listentotaxman check --income 50000
```

### With Student Loan

```bash
listentotaxman check --income 60000 --student-loan plan2
```

### With Pension Contribution

```bash
# Percentage-based
listentotaxman check --income 75000 --pension 5%

# Fixed amount
listentotaxman check --income 75000 --pension 5000
```

### With Tax Code

```bash
listentotaxman check --income 80000 --tax-code 1257L
```

### With Marriage Allowance

```bash
# Married with partner income
listentotaxman check --income 100000 --married --partner-income 25000

# Just married status (no partner income specified)
listentotaxman check --income 100000 --married
```

### With Blind Person's Allowance

```bash
listentotaxman check --income 60000 --blind
```

### Exempt from National Insurance

```bash
# For example, if working past state pension age
listentotaxman check --income 50000 --no-ni
```

### Multiple Parameters

```bash
listentotaxman check \
  --income 120000 \
  --pension 8% \
  --student-loan postgraduate \
  --tax-code K12 \
  --married \
  --partner-income 30000 \
  --verbose
```

### Contractor Rate Calculation

Calculate your effective hourly rate after all deductions:

```bash
listentotaxman check --income 80000 --pension 5% --period hourly
```

This divides all yearly values by 2080 hours (52 weeks × 40 hours), showing you exactly what you earn per hour after tax, NI, and pension contributions.

### Monthly Budgeting

See your monthly take-home pay:

```bash
listentotaxman check --income 60000 --period monthly
```

Perfect for understanding your monthly income for budgeting purposes.

### Weekly Income

Calculate weekly income (useful for weekly budgets):

```bash
listentotaxman check --income 52000 --period weekly
```

### JSON Output for Scripting

```bash
# Get net pay using jq
listentotaxman check --income 100000 --json | jq '.net_pay'

# Compare different salaries
for salary in 50000 75000 100000; do
  echo "Salary: £$salary"
  listentotaxman check --income $salary --json | jq '.net_pay'
done
```

## Testing

### Running Tests

The project maintains comprehensive test coverage with 200 test cases.

```bash
# Run all tests (quiet mode)
make test

# Run with verbose output
make test-verbose

# Run with race detector
make test-race

# Generate coverage report
make test-coverage
# Opens coverage.html in your browser

# Run CI checks (includes coverage threshold)
make test-ci
```

### Test Coverage

The project maintains **≥90% test coverage** (91.3%) across core packages. Coverage excludes the `internal/testutil` package as it contains test helpers.

View coverage report:
```bash
make test-coverage
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

Current coverage by package:
- `internal/config`: 94.3% (Configuration loading)
- `internal/display`: 94.0% (Display formatting)
- `cmd`: 88.9% (Command logic, validation, integration)
- `internal/client`: 79.2% (HTTP client, API interactions)

### Updating Golden Files

Display output tests use golden files for regression testing:

```bash
# Update all golden files after intentional output changes
make test-update-golden

# Review changes before committing
git diff .goldenfiles/
```

### Test Structure

```
cmd/                          - Command tests (validation, parsing, integration)
  check_test.go              - Validation logic tests
  check_logic_test.go        - Calculation and date logic tests
  check_integration_test.go  - End-to-end check command tests
  compare_parsing_test.go    - Argument parsing tests
  compare_validation_test.go - Input validation tests
  compare_integration_test.go - End-to-end compare command tests
internal/client/              - API client tests (HTTP mocking)
internal/config/              - Configuration loading tests
internal/display/             - Display formatting tests
  table_test.go              - Table display tests
  compare_test.go            - Comparison display tests
internal/testutil/            - Shared test utilities and mocks
  testutil.go                - Test helpers (excluded from coverage)
  mocks.go                   - HTTP mocking
  fixtures.go                - Test data
```

### Writing Tests

Tests use table-driven patterns and run in parallel where safe:

```go
func TestExample(t *testing.T) {
    t.Parallel() // Safe for stateless tests
    
    tests := []struct {
        name string
        input int
        want string
    }{
        {"zero", 0, "£0.00"},
        {"positive", 1234, "£1,234.00"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := formatCurrency(tt.input)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

## Building with Version Info

To build with version information:

```bash
VERSION="1.0.0"
GIT_COMMIT=$(git rev-parse HEAD)
BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build -ldflags "\
  -X main.Version=$VERSION \
  -X main.GitCommit=$GIT_COMMIT \
  -X main.BuildDate=$BUILD_DATE" \
  -o listentotaxman
```

## API Reference

This CLI uses the [listentotaxman.com](https://listentotaxman.com/) API. All calculations are performed by their service.

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Author

Michael Heap (@mheap)
