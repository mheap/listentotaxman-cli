package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mheap/listentotaxman-cli/internal/client"
	"github.com/mheap/listentotaxman-cli/internal/config"
	"github.com/mheap/listentotaxman-cli/internal/display"
	"github.com/mheap/listentotaxman-cli/internal/types"
)

const (
	flagValueTrue  = "true"
	periodFlagName = "--period"
)

// ComparisonOption holds one option's label and tax request parameters
type ComparisonOption struct {
	Label   string
	Request *types.TaxRequest
}

// clientFactory is a function that creates a new API client (can be mocked in tests)
var clientFactory = func() *client.Client {
	return client.New()
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

func runCompare(cmd *cobra.Command, _ []string) error {
	// Check for help flag early
	if isHelpRequested() {
		return cmd.Help()
	}

	// Load config file
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Parse and validate options
	globalFlags, options, err := parseAndValidateOptions(cfg)
	if err != nil {
		return err
	}

	// Get and validate period
	period, err := getComparePeriod(globalFlags, cfg)
	if err != nil {
		return err
	}

	// Calculate tax for all options
	results, err := calculateTaxForOptions(options)
	if err != nil {
		return err
	}

	// Display results
	displayCompareResults(results, period, globalFlags)

	return nil
}

// isHelpRequested checks if the user requested help
func isHelpRequested() bool {
	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

// parseAndValidateOptions parses and validates comparison options
func parseAndValidateOptions(cfg *config.Config) (map[string]string, []ComparisonOption, error) {
	// Parse global flags and comparison options from args
	globalFlags, options, err := parseComparisonArgs(os.Args, cfg)
	if err != nil {
		return nil, nil, err
	}

	// Validate number of options
	if len(options) < 2 {
		return nil, nil, fmt.Errorf("at least 2 options required for comparison (use --option to define each scenario)")
	}
	if len(options) > 4 {
		return nil, nil, fmt.Errorf("maximum 4 options supported for comparison (found %d)", len(options))
	}

	// Validate each option
	for i := range options {
		if err := validateOption(&options[i], cfg); err != nil {
			return nil, nil, err
		}
	}

	return globalFlags, options, nil
}

// getComparePeriod gets and validates the period for comparison
func getComparePeriod(globalFlags map[string]string, cfg *config.Config) (string, error) {
	// Get period (global flag > config > default)
	period := periodYearly
	if periodFlag, ok := globalFlags["period"]; ok && periodFlag != "" {
		period = periodFlag
	} else if cfg.Defaults.Period != "" {
		period = cfg.Defaults.Period
	}

	// Validate period
	if err := validatePeriod(period); err != nil {
		return "", err
	}

	return period, nil
}

// calculateTaxForOptions calculates tax for all comparison options
func calculateTaxForOptions(options []ComparisonOption) ([]types.ComparisonResult, error) {
	apiClient := clientFactory()
	results := make([]types.ComparisonResult, len(options))

	for i, opt := range options {
		resp, err := apiClient.CalculateTax(opt.Request)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate tax for option '%s': %w", opt.Label, err)
		}
		results[i] = types.ComparisonResult{
			Label:    opt.Label,
			Request:  opt.Request,
			Response: resp,
		}
	}

	return results, nil
}

// displayCompareResults displays the comparison results
func displayCompareResults(results []types.ComparisonResult, period string, globalFlags map[string]string) {
	jsonFlag := globalFlags["json"] == flagValueTrue
	verboseFlag := globalFlags["verbose"] == flagValueTrue

	if jsonFlag {
		display.ComparisonJSON(results, period)
	} else {
		display.Comparison(results, period, verboseFlag)
	}
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
		if arg == periodFlagName && i+1 < len(args) {
			globalFlags["period"] = args[i+1]
			i++ // Skip value
		} else if arg == "--json" {
			globalFlags["json"] = flagValueTrue
		} else if arg == "--verbose" {
			globalFlags["verbose"] = flagValueTrue
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
		if arg == periodFlagName || arg == "--json" || arg == "--verbose" {
			if arg == periodFlagName && i+1 < len(chunk) {
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
				flags[flagName] = flagValueTrue
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

	// Apply flag and config values to request
	if err := applyCompareTaxRequestDefaults(flags, cfg, req); err != nil {
		return nil, err
	}

	return req, nil
}

// applyCompareTaxRequestDefaults applies flag and config defaults to a tax request
func applyCompareTaxRequestDefaults(flags map[string]string, cfg *config.Config, req *types.TaxRequest) error {
	// Apply string fields
	applyStringField(flags, cfg, req, "year", &req.Year, cfg.Defaults.Year, getDefaultYear())
	applyStringField(flags, cfg, req, "region", &req.TaxRegion, cfg.Defaults.Region, "uk")
	applyStringField(flags, cfg, req, "age", &req.Age, cfg.Defaults.Age, "0")
	applyStringField(flags, cfg, req, "pension", &req.Pension, cfg.Defaults.Pension, "")
	applyStringField(flags, cfg, req, "student-loan", &req.Plan, cfg.Defaults.StudentLoan, "")
	applyStringField(flags, cfg, req, "tax-code", &req.TaxCode, cfg.Defaults.TaxCode, "")

	// Apply integer fields with error handling
	if err := applyIntField(flags, cfg, req, "extra", &req.Extra, cfg.Defaults.Extra); err != nil {
		return err
	}

	// Income is required - check if it exists in flags
	if _, hasIncome := flags["income"]; !hasIncome {
		return fmt.Errorf("income is required")
	}
	if err := applyIntField(flags, cfg, req, "income", &req.GrossWage, 0); err != nil {
		return err
	}

	if err := applyIntField(flags, cfg, req, "partner-income", &req.PartnerGrossWage, cfg.Defaults.PartnerIncome); err != nil {
		return err
	}

	// Apply boolean fields
	applyBoolField(flags, cfg, &req.Married, "married", cfg.Defaults.Married)
	applyBoolField(flags, cfg, &req.Blind, "blind", cfg.Defaults.Blind)
	applyBoolField(flags, cfg, &req.ExNI, "no-ni", cfg.Defaults.NoNI)

	// Normalise region (england -> uk)
	req.TaxRegion = normalizeRegion(req.TaxRegion)

	return nil
}

// applyStringField applies a string field from flags or config defaults
func applyStringField(flags map[string]string, _ *config.Config, _ *types.TaxRequest, flagName string, target *string, configDefault, hardDefault string) {
	if val, ok := flags[flagName]; ok {
		*target = val
	} else if configDefault != "" {
		*target = configDefault
	} else if hardDefault != "" {
		*target = hardDefault
	}
}

// applyIntField applies an integer field from flags or config defaults
func applyIntField(flags map[string]string, _ *config.Config, _ *types.TaxRequest, flagName string, target *int, configDefault int) error {
	if val, ok := flags[flagName]; ok {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("%s must be a valid number: %s", flagName, val)
		}
		*target = intVal
	} else if configDefault != 0 {
		*target = configDefault
	}
	return nil
}

// applyBoolField applies a boolean field from flags or config defaults
func applyBoolField(flags map[string]string, _ *config.Config, target *string, flagName string, configDefault bool) {
	if val, ok := flags[flagName]; ok && val == flagValueTrue {
		*target = "y"
	} else if configDefault {
		*target = "y"
	}
}

// validateOption validates a single comparison option
func validateOption(opt *ComparisonOption, _ *config.Config) error {
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
