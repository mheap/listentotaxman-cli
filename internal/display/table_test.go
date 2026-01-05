package display

import (
	"testing"

	"github.com/mheap/listentotaxman-cli/internal/testutil"
	"github.com/mheap/listentotaxman-cli/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestFormatCurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		amount float64
		want   string
	}{
		{"basic amount", 1234.56, "£1,234.56"},
		{"large amount", 1000000.00, "£1,000,000.00"},
		{"zero", 0.00, "£0.00"},
		{"rounding down", 123.456, "£123.46"},
		{"negative", -500.00, "£-500.00"},
		{"rounding up", 999.999, "£1,000.00"},
		{"no decimals", 5000.00, "£5,000.00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatCurrency(tt.amount)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAddThousandSeparators(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"four digits", "1234", "1,234"},
		{"seven digits", "1234567", "1,234,567"},
		{"three digits", "123", "123"},
		{"negative", "-1234", "-1,234"},
		{"single digit", "1", "1"},
		{"large number", "123456789", "123,456,789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := addThousandSeparators(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPeriodDivisor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		period string
		want   float64
	}{
		{"yearly", "yearly", 1.0},
		{"monthly", "monthly", 12.0},
		{"weekly", "weekly", 52.0},
		{"daily", "daily", 365.0},
		{"hourly", "hourly", 2080.0},
		{"invalid", "invalid", 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getPeriodDivisor(tt.period)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPeriodLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		period string
		want   string
	}{
		{"yearly", "yearly", "Yearly"},
		{"monthly", "monthly", "Monthly"},
		{"weekly", "weekly", "Weekly"},
		{"daily", "daily", "Daily"},
		{"hourly", "hourly", "Hourly"},
		{"invalid", "invalid", "Yearly"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getPeriodLabel(tt.period)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDisplaySummary_BasicOutput(t *testing.T) {
	resp := testutil.CreateSampleTaxResponse()
	req := testutil.CreateSampleTaxRequest()

	output := testutil.CaptureStdout(t, func() {
		DisplaySummary(resp, "yearly", req)
	})

	// Verify table structure
	assert.Contains(t, output, "╔")
	assert.Contains(t, output, "╗")
	assert.Contains(t, output, "╚")
	assert.Contains(t, output, "╝")
	assert.Contains(t, output, "Tax Calculation for 2024")
	assert.Contains(t, output, "Gross Salary")
	assert.Contains(t, output, "£50,000.00")
	assert.Contains(t, output, "Tax Paid")
	assert.Contains(t, output, "National Insurance")
	assert.Contains(t, output, "Net Pay")
}

func TestDisplaySummary_WithMarriedStatus(t *testing.T) {
	resp := testutil.CreateSampleTaxResponse()
	req := testutil.CreateSampleTaxRequest(func(r *types.TaxRequest) {
		r.Married = "y"
	})

	output := testutil.CaptureStdout(t, func() {
		DisplaySummary(resp, "yearly", req)
	})

	assert.Contains(t, output, "Married")
}

func TestDisplaySummary_WithAllStatus(t *testing.T) {
	resp := testutil.CreateSampleTaxResponse()
	req := testutil.CreateSampleTaxRequest(func(r *types.TaxRequest) {
		r.Married = "y"
		r.Blind = "y"
		r.ExNI = "y"
	})

	output := testutil.CaptureStdout(t, func() {
		DisplaySummary(resp, "yearly", req)
	})

	assert.Contains(t, output, "Married")
	assert.Contains(t, output, "Blind Allowance")
	assert.Contains(t, output, "NI Exempt")
}

func TestDisplaySummary_NoStatus(t *testing.T) {
	resp := testutil.CreateSampleTaxResponse()
	req := testutil.CreateSampleTaxRequest()

	output := testutil.CaptureStdout(t, func() {
		DisplaySummary(resp, "yearly", req)
	})

	// Should not have status line
	lineCount := len(output)
	// Basic output without status should be shorter
	assert.Greater(t, 1500, lineCount)
}

func TestDisplayDetailed_StatusLine(t *testing.T) {
	resp := testutil.CreateSampleTaxResponse()
	req := testutil.CreateSampleTaxRequest(func(r *types.TaxRequest) {
		r.Married = "y"
		r.Blind = "y"
	})

	output := testutil.CaptureStdout(t, func() {
		DisplayDetailed(resp, "yearly", req)
	})

	assert.Contains(t, output, "Status: Married • Blind Allowance")
}
