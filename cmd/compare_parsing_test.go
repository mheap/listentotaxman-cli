package cmd

import (
	"testing"
	"time"

	"github.com/mheap/listentotaxman-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseComparisonArgs_BasicTwoOptions(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	args := []string{
		"listentotaxman",
		"compare",
		"--option", "Job 1", "--income", "100000",
		"--option", "Job 2", "--income", "120000",
	}

	globalFlags, options, err := parseComparisonArgs(args, cfg)
	require.NoError(t, err)
	assert.Len(t, options, 2)

	// Check first option
	assert.Equal(t, "Job 1", options[0].Label)
	assert.Equal(t, 100000, options[0].Request.GrossWage)

	// Check second option
	assert.Equal(t, "Job 2", options[1].Label)
	assert.Equal(t, 120000, options[1].Request.GrossWage)

	// No global flags
	assert.Empty(t, globalFlags["period"])
	assert.Empty(t, globalFlags["json"])
	assert.Empty(t, globalFlags["verbose"])
}

func TestParseComparisonArgs_WithGlobalFlags(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	args := []string{
		"listentotaxman",
		"compare",
		"--period", "monthly",
		"--json",
		"--verbose",
		"--option", "Job 1", "--income", "100000",
		"--option", "Job 2", "--income", "120000",
	}

	globalFlags, options, err := parseComparisonArgs(args, cfg)
	require.NoError(t, err)
	assert.Len(t, options, 2)

	// Check global flags
	assert.Equal(t, "monthly", globalFlags["period"])
	assert.Equal(t, "true", globalFlags["json"])
	assert.Equal(t, "true", globalFlags["verbose"])
}

func TestParseComparisonArgs_FourOptions(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	args := []string{
		"listentotaxman",
		"compare",
		"--option", "Low", "--income", "80000",
		"--option", "Mid", "--income", "100000",
		"--option", "High", "--income", "120000",
		"--option", "Higher", "--income", "140000",
	}

	globalFlags, options, err := parseComparisonArgs(args, cfg)
	require.NoError(t, err)
	assert.Len(t, options, 4)
	assert.Equal(t, "Low", options[0].Label)
	assert.Equal(t, "Mid", options[1].Label)
	assert.Equal(t, "High", options[2].Label)
	assert.Equal(t, "Higher", options[3].Label)
	assert.Empty(t, globalFlags["period"])
}

func TestParseComparisonArgs_NoOptions(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	args := []string{
		"listentotaxman",
		"compare",
		"--period", "monthly",
	}

	_, _, err := parseComparisonArgs(args, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no options specified")
}

func TestParseComparisonArgs_NoCompareCommand(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	args := []string{
		"listentotaxman",
		"check",
		"--income", "100000",
	}

	_, _, err := parseComparisonArgs(args, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "compare command not found")
}

func TestParseOptionChunk_BasicOption(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	chunk := []string{
		"--option", "Job 1", "--income", "100000",
	}

	opt, err := parseOptionChunk(chunk, cfg)
	require.NoError(t, err)
	assert.Equal(t, "Job 1", opt.Label)
	assert.Equal(t, 100000, opt.Request.GrossWage)
	assert.Equal(t, "uk", opt.Request.TaxRegion) // default
	assert.Equal(t, "0", opt.Request.Age)        // default
}

func TestParseOptionChunk_WithAllFlags(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	chunk := []string{
		"--option", "Complex Job",
		"--income", "100000",
		"--year", "2024",
		"--region", "scotland",
		"--age", "45",
		"--pension", "5%",
		"--student-loan", "plan2",
		"--extra", "2000",
		"--tax-code", "1257L",
		"--married",
		"--blind",
		"--no-ni",
		"--partner-income", "25000",
	}

	opt, err := parseOptionChunk(chunk, cfg)
	require.NoError(t, err)
	assert.Equal(t, "Complex Job", opt.Label)
	assert.Equal(t, 100000, opt.Request.GrossWage)
	assert.Equal(t, "2024", opt.Request.Year)
	assert.Equal(t, "scotland", opt.Request.TaxRegion)
	assert.Equal(t, "45", opt.Request.Age)
	assert.Equal(t, "5%", opt.Request.Pension)
	assert.Equal(t, "plan2", opt.Request.Plan)
	assert.Equal(t, 2000, opt.Request.Extra)
	assert.Equal(t, "1257L", opt.Request.TaxCode)
	assert.Equal(t, "y", opt.Request.Married)
	assert.Equal(t, "y", opt.Request.Blind)
	assert.Equal(t, "y", opt.Request.ExNI)
	assert.Equal(t, 25000, opt.Request.PartnerGrossWage)
}

func TestParseOptionChunk_NoLabel(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	chunk := []string{
		"--option",
	}

	_, err := parseOptionChunk(chunk, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires a label")
}

func TestParseOptionChunk_SkipsGlobalFlags(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	chunk := []string{
		"--option", "Job 1",
		"--income", "100000",
		"--period", "monthly", // Should be skipped
		"--json",    // Should be skipped
		"--verbose", // Should be skipped
	}

	opt, err := parseOptionChunk(chunk, cfg)
	require.NoError(t, err)
	assert.Equal(t, "Job 1", opt.Label)
	assert.Equal(t, 100000, opt.Request.GrossWage)
}

func TestBuildTaxRequest_WithConfigDefaults(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{
			Year:        "2024",
			Region:      "scotland",
			Age:         "30",
			Pension:     "3%",
			StudentLoan: "plan1",
			TaxCode:     "1100L",
			Extra:       1000,
			Married:     true,
			Blind:       true,
			NoNI:        true,
		},
	}

	flags := map[string]string{
		"income": "100000",
	}

	req, err := buildTaxRequest(flags, cfg)
	require.NoError(t, err)

	// Should use config defaults
	assert.Equal(t, 100000, req.GrossWage)
	assert.Equal(t, "2024", req.Year)
	assert.Equal(t, "scotland", req.TaxRegion)
	assert.Equal(t, "30", req.Age)
	assert.Equal(t, "3%", req.Pension)
	assert.Equal(t, "plan1", req.Plan)
	assert.Equal(t, "1100L", req.TaxCode)
	assert.Equal(t, 1000, req.Extra)
	assert.Equal(t, "y", req.Married)
	assert.Equal(t, "y", req.Blind)
	assert.Equal(t, "y", req.ExNI)
}

func TestBuildTaxRequest_FlagsOverrideConfig(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{
			Year:        "2024",
			Region:      "scotland",
			Age:         "30",
			Pension:     "3%",
			StudentLoan: "plan1",
		},
	}

	flags := map[string]string{
		"income":       "100000",
		"year":         "2025",
		"region":       "uk",
		"age":          "45",
		"pension":      "5%",
		"student-loan": "plan2",
	}

	req, err := buildTaxRequest(flags, cfg)
	require.NoError(t, err)

	// Flags should override config
	assert.Equal(t, "2025", req.Year)
	assert.Equal(t, "uk", req.TaxRegion)
	assert.Equal(t, "45", req.Age)
	assert.Equal(t, "5%", req.Pension)
	assert.Equal(t, "plan2", req.Plan)
}

func TestBuildTaxRequest_MissingIncome(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	flags := map[string]string{
		"year": "2024",
	}

	_, err := buildTaxRequest(flags, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "income is required")
}

func TestBuildTaxRequest_InvalidIncome(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	flags := map[string]string{
		"income": "not-a-number",
	}

	_, err := buildTaxRequest(flags, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "income must be a valid number")
}

func TestBuildTaxRequest_InvalidExtra(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	flags := map[string]string{
		"income": "100000",
		"extra":  "not-a-number",
	}

	_, err := buildTaxRequest(flags, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "extra must be a valid number")
}

func TestBuildTaxRequest_InvalidPartnerIncome(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	flags := map[string]string{
		"income":         "100000",
		"partner-income": "not-a-number",
	}

	_, err := buildTaxRequest(flags, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "partner-income must be a valid number")
}

func TestBuildTaxRequest_BooleanFlags(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	flags := map[string]string{
		"income":  "100000",
		"married": "true",
		"blind":   "true",
		"no-ni":   "true",
	}

	req, err := buildTaxRequest(flags, cfg)
	require.NoError(t, err)
	assert.Equal(t, "y", req.Married)
	assert.Equal(t, "y", req.Blind)
	assert.Equal(t, "y", req.ExNI)
}

func TestBuildTaxRequest_RegionNormalization(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	flags := map[string]string{
		"income": "100000",
		"region": "england",
	}

	req, err := buildTaxRequest(flags, cfg)
	require.NoError(t, err)
	// "england" should be normalized to "uk"
	assert.Equal(t, "uk", req.TaxRegion)
}

func TestBuildTaxRequest_DefaultYearUsesSmartDefault(t *testing.T) {
	// Mock timeNowFunc to return a fixed date
	originalTimeNow := timeNowFunc
	defer func() { timeNowFunc = originalTimeNow }()

	// Before April 5 -> current year
	timeNowFunc = func() time.Time {
		return time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)
	}

	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	flags := map[string]string{
		"income": "100000",
	}

	req, err := buildTaxRequest(flags, cfg)
	require.NoError(t, err)
	assert.Equal(t, "2025", req.Year) // Should use 2025 (2026 - 1)
}

func TestBuildTaxRequest_PartnerIncomeFromConfig(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{
			PartnerIncome: 30000,
		},
	}

	flags := map[string]string{
		"income": "100000",
	}

	req, err := buildTaxRequest(flags, cfg)
	require.NoError(t, err)
	assert.Equal(t, 30000, req.PartnerGrossWage)
}

func TestBuildTaxRequest_StaticFields(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	flags := map[string]string{
		"income": "100000",
	}

	req, err := buildTaxRequest(flags, cfg)
	require.NoError(t, err)

	// These fields should always be set
	assert.Equal(t, "json", req.Response)
	assert.Equal(t, "1", req.Time)
}
