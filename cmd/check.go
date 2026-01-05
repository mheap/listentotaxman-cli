// Package cmd implements the CLI commands for the listentotaxman application.
package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/mheap/listentotaxman-cli/internal/client"
	"github.com/mheap/listentotaxman-cli/internal/config"
	"github.com/mheap/listentotaxman-cli/internal/display"
	"github.com/mheap/listentotaxman-cli/internal/types"
)

var (
	flagYear          string
	flagRegion        string
	flagAge           string
	flagPension       string
	flagIncome        int
	flagStudentLoan   string
	flagExtra         int
	flagTaxCode       string
	flagJSON          bool
	flagVerbose       bool
	flagPeriod        string
	flagMarried       bool
	flagBlind         bool
	flagNoNI          bool
	flagPartnerIncome int
)

const (
	periodYearly = "yearly"
)

// timeNowFunc allows time mocking in tests
var timeNowFunc = time.Now

// clientFactory allows API client mocking in tests
var checkClientFactory = func() *client.Client {
	return client.New()
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check tax calculation for a given salary",
	Long:  `Calculate UK tax and national insurance for a given salary and parameters.`,
	RunE:  runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Define flags
	checkCmd.Flags().StringVar(&flagYear, "year", "", "Tax year (defaults to current tax year)")
	checkCmd.Flags().StringVar(&flagRegion, "region", "", "Tax region (default: uk)")
	checkCmd.Flags().StringVar(&flagAge, "age", "", "Age (default: 0)")
	checkCmd.Flags().StringVar(&flagPension, "pension", "", "Pension contribution (e.g., 3% or 3000)")
	checkCmd.Flags().IntVar(&flagIncome, "income", 0, "Gross annual salary (required)")
	checkCmd.Flags().StringVar(&flagStudentLoan, "student-loan", "", "Student loan plan (e.g., plan1, plan2, plan4, postgraduate, scottish)")
	checkCmd.Flags().IntVar(&flagExtra, "extra", 0, "Extra income/deductions")
	checkCmd.Flags().StringVar(&flagTaxCode, "tax-code", "", "Tax code (e.g., 1257L, K12)")
	checkCmd.Flags().BoolVar(&flagMarried, "married", false, "Married status (enables marriage allowance)")
	checkCmd.Flags().BoolVar(&flagBlind, "blind", false, "Blind person's allowance")
	checkCmd.Flags().BoolVar(&flagNoNI, "no-ni", false, "Exempt from National Insurance")
	checkCmd.Flags().IntVar(&flagPartnerIncome, "partner-income", 0, "Partner's gross wage (requires --married)")
	checkCmd.Flags().BoolVar(&flagJSON, "json", false, "Output as JSON")
	checkCmd.Flags().BoolVar(&flagVerbose, "verbose", false, "Show detailed breakdown")
	checkCmd.Flags().StringVar(&flagPeriod, "period", "", "Display period (yearly, monthly, weekly, daily, hourly) (default: yearly)")

	// Mark required flags
	_ = checkCmd.MarkFlagRequired("income")
}

// getDefaultYear returns the default tax year based on current date
// If today is after April 5th, use current year. Otherwise use previous year.
func getDefaultYear() string {
	now := timeNowFunc()
	year := now.Year()

	// Create April 5th of current year
	aprilFifth := time.Date(year, time.April, 5, 0, 0, 0, 0, time.UTC)

	// If we're after April 5th, use current year. Otherwise use previous year.
	if now.After(aprilFifth) {
		return strconv.Itoa(year)
	}
	return strconv.Itoa(year - 1)
}

// getPeriodDivisor returns the divisor for a given period
func getPeriodDivisor(period string) float64 {
	switch period {
	case periodYearly:
		return 1.0
	case "monthly":
		return 12.0
	case "weekly":
		return 52.0
	case "daily":
		return 365.0
	case "hourly":
		return 2080.0
	default:
		return 1.0
	}
}

// adjustResponseForPeriod creates a copy of the response with values adjusted for the period
func adjustResponseForPeriod(resp *types.TaxResponse, period string) *types.TaxResponse {
	if period == periodYearly {
		return resp
	}

	divisor := getPeriodDivisor(period)

	// Create a copy
	adjusted := *resp

	// Adjust all monetary values
	adjusted.TaxablePay /= divisor
	adjusted.GrossPay /= divisor
	adjusted.AdditionalGross /= divisor
	adjusted.TaxFreeAllowance /= divisor
	adjusted.TaxPaid /= divisor
	adjusted.NationalInsurance /= divisor
	adjusted.NetPay /= divisor
	adjusted.StudentLoanRepayment /= divisor
	adjusted.PensionHMRC /= divisor
	adjusted.PensionYou /= divisor
	adjusted.PensionClaimback /= divisor
	adjusted.EmployersNI /= divisor
	adjusted.TaxFreeMarried /= divisor
	adjusted.TaxFreeMarriageAllowance /= divisor
	adjusted.GrossSacrifice /= divisor
	adjusted.ChildcareAmount /= divisor

	// Adjust tax brackets
	for key, bracket := range adjusted.TaxDue {
		bracket.Amount /= divisor
		adjusted.TaxDue[key] = bracket
	}

	// Also adjust previous year if present
	if adjusted.Previous != nil {
		adjusted.Previous = adjustResponseForPeriod(adjusted.Previous, period)
	}

	return &adjusted
}

func runCheck(cmd *cobra.Command, _ []string) error {
	// Load config file
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Build and validate request
	req, err := buildCheckTaxRequest(cmd, cfg)
	if err != nil {
		return err
	}

	// Get and validate period
	period, err := getPeriod(cfg)
	if err != nil {
		return err
	}

	// Call API
	apiClient := checkClientFactory()
	resp, err := apiClient.CalculateTax(req)
	if err != nil {
		return fmt.Errorf("failed to calculate tax: %w", err)
	}

	// Display result
	return displayCheckResult(resp, period, req)
}

// buildCheckTaxRequest builds and validates a TaxRequest from flags and config
func buildCheckTaxRequest(cmd *cobra.Command, cfg *config.Config) (*types.TaxRequest, error) {
	req := &types.TaxRequest{}

	// Apply flag and config values
	applyCheckRequestDefaults(cmd, cfg, req)

	// Validate the request
	if err := validateCheckRequest(req); err != nil {
		return nil, err
	}

	return req, nil
}

// applyCheckRequestDefaults applies flag and config values to the request
func applyCheckRequestDefaults(cmd *cobra.Command, cfg *config.Config, req *types.TaxRequest) {
	// Apply string fields
	applyCheckStringFields(cfg, req)

	// Apply integer and numeric fields
	applyCheckNumericFields(cmd, cfg, req)

	// Apply boolean fields
	applyCheckBooleanFields(cmd, cfg, req)

	// Normalise region (england -> uk)
	req.TaxRegion = normalizeRegion(req.TaxRegion)
}

// applyCheckStringFields applies string flag and config values
func applyCheckStringFields(cfg *config.Config, req *types.TaxRequest) {
	// Year: flag > config > smart default
	if flagYear != "" {
		req.Year = flagYear
	} else if cfg.Defaults.Year != "" {
		req.Year = cfg.Defaults.Year
	} else {
		req.Year = getDefaultYear()
	}

	// Region: flag > config > default "uk"
	if flagRegion != "" {
		req.TaxRegion = flagRegion
	} else if cfg.Defaults.Region != "" {
		req.TaxRegion = cfg.Defaults.Region
	} else {
		req.TaxRegion = "uk"
	}

	// Age: flag > config > default "0"
	if flagAge != "" {
		req.Age = flagAge
	} else if cfg.Defaults.Age != "" {
		req.Age = cfg.Defaults.Age
	} else {
		req.Age = "0"
	}

	// Pension: flag > config > empty
	if flagPension != "" {
		req.Pension = flagPension
	} else if cfg.Defaults.Pension != "" {
		req.Pension = cfg.Defaults.Pension
	}

	// Student loan: flag > config > empty
	if flagStudentLoan != "" {
		req.Plan = flagStudentLoan
	} else if cfg.Defaults.StudentLoan != "" {
		req.Plan = cfg.Defaults.StudentLoan
	}

	// Tax code: flag > config > empty
	if flagTaxCode != "" {
		req.TaxCode = flagTaxCode
	} else if cfg.Defaults.TaxCode != "" {
		req.TaxCode = cfg.Defaults.TaxCode
	}
}

// applyCheckNumericFields applies numeric flag and config values
func applyCheckNumericFields(cmd *cobra.Command, cfg *config.Config, req *types.TaxRequest) {
	// Extra: flag > config > 0
	if cmd.Flags().Changed("extra") {
		req.Extra = flagExtra
	} else if cfg.Defaults.Extra != 0 {
		req.Extra = cfg.Defaults.Extra
	}

	// Income is always from flag (required)
	req.GrossWage = flagIncome

	// Partner income: flag > config > 0
	if cmd.Flags().Changed("partner-income") {
		req.PartnerGrossWage = flagPartnerIncome
	} else if cfg.Defaults.PartnerIncome != 0 {
		req.PartnerGrossWage = cfg.Defaults.PartnerIncome
	}
}

// applyCheckBooleanFields applies boolean flag and config values
func applyCheckBooleanFields(cmd *cobra.Command, cfg *config.Config, req *types.TaxRequest) {
	// Married: flag > config > default false
	if cmd.Flags().Changed("married") {
		if flagMarried {
			req.Married = "y"
		}
	} else if cfg.Defaults.Married {
		req.Married = "y"
	}

	// Blind: flag > config > default false
	if cmd.Flags().Changed("blind") {
		if flagBlind {
			req.Blind = "y"
		}
	} else if cfg.Defaults.Blind {
		req.Blind = "y"
	}

	// No NI: flag > config > default false
	if cmd.Flags().Changed("no-ni") {
		if flagNoNI {
			req.ExNI = "y"
		}
	} else if cfg.Defaults.NoNI {
		req.ExNI = "y"
	}
}

// validateCheckRequest validates the tax request
func validateCheckRequest(req *types.TaxRequest) error {
	// Validate year is a 4-digit number
	if len(req.Year) != 4 {
		return fmt.Errorf("year must be a 4-digit number, got: %s", req.Year)
	}
	if _, yearErr := strconv.Atoi(req.Year); yearErr != nil {
		return fmt.Errorf("year must be a valid number: %s", req.Year)
	}

	// Validate income is positive
	if req.GrossWage <= 0 {
		return fmt.Errorf("income must be greater than 0")
	}

	// Validate partner income requires married flag
	if req.PartnerGrossWage > 0 && req.Married != "y" {
		return fmt.Errorf("--partner-income requires --married flag\nHint: Use --married --partner-income %d", req.PartnerGrossWage)
	}

	// Validate partner income is not negative
	if req.PartnerGrossWage < 0 {
		return fmt.Errorf("--partner-income cannot be negative")
	}

	// Validate student loan plan
	if req.Plan != "" {
		if err := validateStudentLoanPlan(req.Plan); err != nil {
			return err
		}
	}

	return nil
}

// validateStudentLoanPlan validates the student loan plan
func validateStudentLoanPlan(plan string) error {
	validPlans := []string{"plan1", "plan2", "plan4", "postgraduate", "scottish"}
	for _, vp := range validPlans {
		if plan == vp {
			return nil
		}
	}
	return fmt.Errorf("invalid student loan plan: %s (must be one of: plan1, plan2, plan4, postgraduate, scottish)", plan)
}

// getPeriod gets and validates the period
func getPeriod(cfg *config.Config) (string, error) {
	// Period: flag > config > default periodYearly
	period := periodYearly
	if flagPeriod != "" {
		period = flagPeriod
	} else if cfg.Defaults.Period != "" {
		period = cfg.Defaults.Period
	}

	// Validate period
	if err := validatePeriod(period); err != nil {
		return "", err
	}

	return period, nil
}

// validatePeriod validates the period value
func validatePeriod(period string) error {
	validPeriods := []string{periodYearly, "monthly", "weekly", "daily", "hourly"}
	for _, vp := range validPeriods {
		if period == vp {
			return nil
		}
	}
	return fmt.Errorf("invalid period: %s (must be one of: yearly, monthly, weekly, daily, hourly)", period)
}

// displayCheckResult displays the tax calculation result
func displayCheckResult(resp *types.TaxResponse, period string, req *types.TaxRequest) error {
	if flagJSON {
		// Output as JSON - adjust response for period
		adjustedResp := adjustResponseForPeriod(resp, period)
		jsonData, err := json.MarshalIndent(adjustedResp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else if flagVerbose {
		// Display detailed breakdown
		display.Detailed(resp, period, req)
	} else {
		// Display summary table
		display.Summary(resp, period, req)
	}

	return nil
}

// normalizeRegion converts region aliases to their canonical form
func normalizeRegion(region string) string {
	if region == "england" {
		return "uk"
	}
	return region
}
