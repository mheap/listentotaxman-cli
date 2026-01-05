package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mheap/listentotaxman-cli/internal/client"
	"github.com/mheap/listentotaxman-cli/internal/config"
	"github.com/mheap/listentotaxman-cli/internal/display"
	"github.com/mheap/listentotaxman-cli/internal/types"
	"github.com/spf13/cobra"
)

// ComparisonOption holds one option's label and tax request parameters
type ComparisonOption struct {
	Label   string
	Request *types.TaxRequest
}

var compareCmd = &cobra.Command{
	Use:   "compare --option LABEL --income AMOUNT [flags] --option LABEL --income AMOUNT [flags]...",
	Short: "Compare tax calculations across multiple scenarios",
	Long: `Compare tax calculations across different job offers, salary levels, or pension contributions.

Each --option group represents one scenario and supports all flags from the 'check' command.
Minimum 2 options required, maximum 4 options supported.

Global Flags (apply to all options):
  --period PERIOD   Display period (yearly, monthly, weekly, daily, hourly)
  --json            Output as JSON comparison object
  --verbose         Show detailed breakdown including tax brackets

Per-Option Flags (use after each --option):
  --income INT         Gross annual salary (required)
  --year YEAR          Tax year (defaults to current tax year)
  --region REGION      Tax region (default: "uk", alias: "england")
  --age AGE            Age (default: "0")
  --pension VALUE      Pension contribution (e.g., "3%" or "3000")
  --student-loan PLAN  Student loan plan (plan1, plan2, plan4, postgraduate, scottish)
  --extra INT          Extra income/deductions
  --tax-code CODE      Tax code (e.g., "1257L")
  --married            Married status (enables marriage allowance)
  --blind              Blind person's allowance
  --no-ni              Exempt from National Insurance
  --partner-income INT Partner's gross wage (requires --married)`,
	Example: `  # Compare two job offers
  listentotaxman compare \
    --option "Current Job" --income 100000 --pension 3% \
    --option "New Offer" --income 120000 --pension 5%

  # Compare monthly take-home across salary levels
  listentotaxman compare --period monthly \
    --option "Low" --income 80000 \
    --option "Mid" --income 100000 \
    --option "High" --income 120000

  # Compare regions
  listentotaxman compare \
    --option "England" --income 100000 --region uk \
    --option "Scotland" --income 100000 --region scotland

  # Compare marriage allowance impact
  listentotaxman compare \
    --option "Single" --income 100000 \
    --option "Married" --income 100000 --married --partner-income 25000

  # Detailed comparison with verbose mode
  listentotaxman compare --verbose \
    --option "Job 1" --income 100000 \
    --option "Job 2" --income 120000`,
	RunE:                  runCompare,
	DisableFlagParsing:    true, // We parse flags manually
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	// Check for help flag early
	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			return cmd.Help()
		}
	}

	// Load config file
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Parse global flags and comparison options from args
	globalFlags, options, err := parseComparisonArgs(os.Args, cfg)
	if err != nil {
		return err
	}

	// Validate number of options
	if len(options) < 2 {
		return fmt.Errorf("at least 2 options required for comparison (use --option to define each scenario)")
	}
	if len(options) > 4 {
		return fmt.Errorf("maximum 4 options supported for comparison (found %d)", len(options))
	}

	// Validate each option
	for i := range options {
		if err := validateOption(&options[i], cfg); err != nil {
			return err
		}
	}

	// Get period (global flag > config > default)
	period := "yearly"
	if periodFlag, ok := globalFlags["period"]; ok && periodFlag != "" {
		period = periodFlag
	} else if cfg.Defaults.Period != "" {
		period = cfg.Defaults.Period
	}

	// Validate period
	validPeriods := []string{"yearly", "monthly", "weekly", "daily", "hourly"}
	isValid := false
	for _, vp := range validPeriods {
		if period == vp {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid period: %s (must be one of: yearly, monthly, weekly, daily, hourly)", period)
	}

	// Call API for each option
	apiClient := client.New()
	results := make([]types.ComparisonResult, len(options))

	for i, opt := range options {
		resp, err := apiClient.CalculateTax(opt.Request)
		if err != nil {
			return fmt.Errorf("failed to calculate tax for option '%s': %w", opt.Label, err)
		}
		results[i] = types.ComparisonResult{
			Label:    opt.Label,
			Request:  opt.Request,
			Response: resp,
		}
	}

	// Display results
	jsonFlag := globalFlags["json"] == "true"
	verboseFlag := globalFlags["verbose"] == "true"

	if jsonFlag {
		display.DisplayComparisonJSON(results, period)
	} else {
		display.DisplayComparison(results, period, verboseFlag)
	}

	return nil
}

// parseComparisonArgs parses command-line args into global flags and comparison options
func parseComparisonArgs(allArgs []string, cfg *config.Config) (map[string]string, []ComparisonOption, error) {
	// Find the "compare" command position
	compareIdx := -1
	for i, arg := range allArgs {
		if arg == "compare" {
			compareIdx = i
			break
		}
	}
	if compareIdx == -1 {
		return nil, nil, fmt.Errorf("compare command not found in args")
	}

	// Get args after "compare"
	args := allArgs[compareIdx+1:]

	// Parse global flags first
	globalFlags := make(map[string]string)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--period" && i+1 < len(args) {
			globalFlags["period"] = args[i+1]
			i++ // Skip value
		} else if arg == "--json" {
			globalFlags["json"] = "true"
		} else if arg == "--verbose" {
			globalFlags["verbose"] = "true"
		}
	}

	// Find all --option positions
	optionIndices := []int{}
	for i, arg := range args {
		if arg == "--option" {
			optionIndices = append(optionIndices, i)
		}
	}

	if len(optionIndices) == 0 {
		return nil, nil, fmt.Errorf("no options specified (use --option to define each scenario)")
	}

	// Split args into chunks between --option flags
	options := []ComparisonOption{}
	for i, startIdx := range optionIndices {
		// Determine end of this chunk
		endIdx := len(args)
		if i+1 < len(optionIndices) {
			endIdx = optionIndices[i+1]
		}

		chunk := args[startIdx:endIdx]

		// Parse this chunk into an option
		opt, err := parseOptionChunk(chunk, cfg)
		if err != nil {
			return nil, nil, err
		}

		options = append(options, opt)
	}

	return globalFlags, options, nil
}

// parseOptionChunk parses a single option chunk (--option label --flag value ...)
func parseOptionChunk(chunk []string, cfg *config.Config) (ComparisonOption, error) {
	if len(chunk) < 2 {
		return ComparisonOption{}, fmt.Errorf("--option requires a label")
	}

	// First arg is "--option", second is the label
	label := chunk[1]

	// Parse remaining args as flag-value pairs
	flags := make(map[string]string)
	for i := 2; i < len(chunk); i++ {
		arg := chunk[i]

		// Skip global flags (they're handled separately)
		if arg == "--period" || arg == "--json" || arg == "--verbose" {
			if arg == "--period" && i+1 < len(chunk) {
				i++ // Skip the value too
			}
			continue
		}

		// Check if this is a flag
		if strings.HasPrefix(arg, "--") {
			flagName := strings.TrimPrefix(arg, "--")

			// Check if next arg is the value
			if i+1 < len(chunk) && !strings.HasPrefix(chunk[i+1], "--") {
				flags[flagName] = chunk[i+1]
				i++ // Skip the value
			} else {
				// Boolean flag (no value)
				flags[flagName] = "true"
			}
		}
	}

	// Build TaxRequest with config defaults and flag overrides
	req, err := buildTaxRequest(flags, cfg)
	if err != nil {
		return ComparisonOption{}, fmt.Errorf("option '%s': %w", label, err)
	}

	return ComparisonOption{
		Label:   label,
		Request: req,
	}, nil
}

// buildTaxRequest builds a TaxRequest from flags with config defaults
func buildTaxRequest(flags map[string]string, cfg *config.Config) (*types.TaxRequest, error) {
	req := &types.TaxRequest{
		Response: "json",
		Time:     "1",
	}

	// Year: flag > config > smart default
	if val, ok := flags["year"]; ok {
		req.Year = val
	} else if cfg.Defaults.Year != "" {
		req.Year = cfg.Defaults.Year
	} else {
		req.Year = getDefaultYear()
	}

	// Region: flag > config > default "uk"
	if val, ok := flags["region"]; ok {
		req.TaxRegion = val
	} else if cfg.Defaults.Region != "" {
		req.TaxRegion = cfg.Defaults.Region
	} else {
		req.TaxRegion = "uk"
	}

	// Age: flag > config > default "0"
	if val, ok := flags["age"]; ok {
		req.Age = val
	} else if cfg.Defaults.Age != "" {
		req.Age = cfg.Defaults.Age
	} else {
		req.Age = "0"
	}

	// Pension: flag > config > empty
	if val, ok := flags["pension"]; ok {
		req.Pension = val
	} else if cfg.Defaults.Pension != "" {
		req.Pension = cfg.Defaults.Pension
	}

	// Student loan: flag > config > empty
	if val, ok := flags["student-loan"]; ok {
		req.Plan = val
	} else if cfg.Defaults.StudentLoan != "" {
		req.Plan = cfg.Defaults.StudentLoan
	}

	// Tax code: flag > config > empty
	if val, ok := flags["tax-code"]; ok {
		req.TaxCode = val
	} else if cfg.Defaults.TaxCode != "" {
		req.TaxCode = cfg.Defaults.TaxCode
	}

	// Extra: flag > config > 0
	if val, ok := flags["extra"]; ok {
		extraInt, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("extra must be a valid number: %s", val)
		}
		req.Extra = extraInt
	} else if cfg.Defaults.Extra != 0 {
		req.Extra = cfg.Defaults.Extra
	}

	// Income: flag only (required per option)
	if val, ok := flags["income"]; ok {
		incomeInt, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("income must be a valid number: %s", val)
		}
		req.GrossWage = incomeInt
	} else {
		return nil, fmt.Errorf("income is required")
	}

	// Married: flag > config > default false
	if val, ok := flags["married"]; ok && val == "true" {
		req.Married = "y"
	} else if cfg.Defaults.Married {
		req.Married = "y"
	}

	// Blind: flag > config > default false
	if val, ok := flags["blind"]; ok && val == "true" {
		req.Blind = "y"
	} else if cfg.Defaults.Blind {
		req.Blind = "y"
	}

	// No NI: flag > config > default false
	if val, ok := flags["no-ni"]; ok && val == "true" {
		req.ExNI = "y"
	} else if cfg.Defaults.NoNI {
		req.ExNI = "y"
	}

	// Partner income: flag > config > 0
	if val, ok := flags["partner-income"]; ok {
		partnerInt, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("partner-income must be a valid number: %s", val)
		}
		req.PartnerGrossWage = partnerInt
	} else if cfg.Defaults.PartnerIncome != 0 {
		req.PartnerGrossWage = cfg.Defaults.PartnerIncome
	}

	// Normalize region (england -> uk)
	req.TaxRegion = normalizeRegion(req.TaxRegion)

	return req, nil
}

// validateOption validates a single comparison option
func validateOption(opt *ComparisonOption, cfg *config.Config) error {
	req := opt.Request

	// Validate year is a 4-digit number
	if len(req.Year) != 4 {
		return fmt.Errorf("option '%s': year must be a 4-digit number, got: %s", opt.Label, req.Year)
	}
	if _, err := strconv.Atoi(req.Year); err != nil {
		return fmt.Errorf("option '%s': year must be a valid number: %s", opt.Label, req.Year)
	}

	// Validate income is positive
	if req.GrossWage <= 0 {
		return fmt.Errorf("option '%s': income must be greater than 0", opt.Label)
	}

	// Validate partner income requires married flag
	if req.PartnerGrossWage > 0 && req.Married != "y" {
		return fmt.Errorf("option '%s': --partner-income requires --married flag\nHint: Use --married --partner-income %d", opt.Label, req.PartnerGrossWage)
	}

	// Validate partner income is not negative
	if req.PartnerGrossWage < 0 {
		return fmt.Errorf("option '%s': --partner-income cannot be negative", opt.Label)
	}

	// Validate student loan plan
	if req.Plan != "" {
		validPlans := []string{"plan1", "plan2", "plan4", "postgraduate", "scottish"}
		isValid := false
		for _, vp := range validPlans {
			if req.Plan == vp {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("option '%s': invalid student loan plan: %s (must be one of: plan1, plan2, plan4, postgraduate, scottish)", opt.Label, req.Plan)
		}
	}

	return nil
}
