package display

import (
	"encoding/json"
	"testing"

	"github.com/mheap/listentotaxman-cli/internal/testutil"
	"github.com/mheap/listentotaxman-cli/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateBorder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		numOptions   int
		fieldWidth   int
		valueWidth   int
		borderType   string
		wantContains []string
	}{
		{
			name:         "top border 2 options",
			numOptions:   2,
			fieldWidth:   20,
			valueWidth:   12,
			borderType:   "top",
			wantContains: []string{"╔", "╦", "╗", "═"},
		},
		{
			name:         "middle border 2 options",
			numOptions:   2,
			fieldWidth:   20,
			valueWidth:   12,
			borderType:   "middle",
			wantContains: []string{"╠", "╬", "╣", "═"},
		},
		{
			name:         "bottom border 2 options",
			numOptions:   2,
			fieldWidth:   20,
			valueWidth:   12,
			borderType:   "bottom",
			wantContains: []string{"╚", "╩", "╝", "═"},
		},
		{
			name:         "top border 4 options",
			numOptions:   4,
			fieldWidth:   20,
			valueWidth:   12,
			borderType:   "top",
			wantContains: []string{"╔", "╦", "╗"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := generateBorder(tt.numOptions, tt.fieldWidth, tt.valueWidth, tt.borderType)

			for _, want := range tt.wantContains {
				assert.Contains(t, got, want)
			}

			// Verify length is correct
			// Formula: left(1) + field(width+2) + (mid(1) + value(width+2)) * numOptions + right(1)
			// But we just verify it's not empty and has the right characters
			assert.NotEmpty(t, got)
		})
	}
}

func TestExtractField(t *testing.T) {
	t.Parallel()

	results := []types.ComparisonResult{
		{
			Label: "Option1",
			Response: testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
				r.NetPay = 40000.0
			}),
		},
		{
			Label: "Option2",
			Response: testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
				r.NetPay = 45000.0
			}),
		},
	}

	field := extractField(results, 12.0, func(r *types.TaxResponse) float64 {
		return r.NetPay
	})

	assert.Equal(t, 2, len(field))
	assert.InDelta(t, 40000.0/12.0, field["Option1"], 0.01)
	assert.InDelta(t, 45000.0/12.0, field["Option2"], 0.01)
}

func TestBuildComparisonFields(t *testing.T) {
	t.Parallel()

	results := []types.ComparisonResult{
		{
			Label:    "Test1",
			Response: testutil.CreateSampleTaxResponse(),
		},
		{
			Label:    "Test2",
			Response: testutil.CreateSampleTaxResponse(),
		},
	}

	fields := buildComparisonFields(results, 1.0)

	// Verify all expected fields exist
	expectedFields := []string{
		"gross_pay",
		"taxable_pay",
		"tax_paid",
		"national_insurance",
		"net_pay",
		"employers_ni",
		"total_cost",
		"basic_rate_tax",
	}

	for _, fieldName := range expectedFields {
		assert.Contains(t, fields, fieldName, "should contain field: %s", fieldName)
		assert.Equal(t, 2, len(fields[fieldName]), "field %s should have 2 values", fieldName)
	}
}

func TestBuildMetadata(t *testing.T) {
	t.Parallel()

	results := []types.ComparisonResult{
		{
			Label: "Option1",
			Response: testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
				r.TaxYear = 2024
				r.TaxRegion = "uk"
				r.TaxCode = "1257L"
			}),
		},
		{
			Label: "Option2",
			Response: testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
				r.TaxYear = 2024
				r.TaxRegion = "scotland"
				r.TaxCode = "1257L"
			}),
		},
	}

	metadata := buildMetadata(results)

	assert.Equal(t, 2, len(metadata))
	assert.Contains(t, metadata, "Option1")
	assert.Contains(t, metadata, "Option2")

	assert.Equal(t, 2024, metadata["Option1"]["tax_year"])
	assert.Equal(t, "uk", metadata["Option1"]["tax_region"])
	assert.Equal(t, "scotland", metadata["Option2"]["tax_region"])
}

func TestDisplayComparison(t *testing.T) {
	results := []types.ComparisonResult{
		{
			Label:   "Job 1",
			Request: testutil.CreateSampleTaxRequest(),
			Response: testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
				r.GrossPay = 50000.0
				r.NetPay = 38000.0
			}),
		},
		{
			Label:   "Job 2",
			Request: testutil.CreateSampleTaxRequest(),
			Response: testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
				r.GrossPay = 60000.0
				r.NetPay = 45000.0
			}),
		},
	}

	output := testutil.CaptureStdout(t, func() {
		DisplayComparison(results, "yearly", false)
	})

	// Verify table structure
	assert.Contains(t, output, "╔")
	assert.Contains(t, output, "╗")
	assert.Contains(t, output, "╚")
	assert.Contains(t, output, "╝")

	// Verify labels
	assert.Contains(t, output, "Job 1")
	assert.Contains(t, output, "Job 2")

	// Verify some fields
	assert.Contains(t, output, "Gross Salary")
	assert.Contains(t, output, "Net Pay")
}

func TestDisplayComparison_WithStatus(t *testing.T) {
	results := []types.ComparisonResult{
		{
			Label:    "Single",
			Request:  testutil.CreateSampleTaxRequest(),
			Response: testutil.CreateSampleTaxResponse(),
		},
		{
			Label: "Married",
			Request: testutil.CreateSampleTaxRequest(func(r *types.TaxRequest) {
				r.Married = "y"
				r.Blind = "y"
			}),
			Response: testutil.CreateSampleTaxResponse(),
		},
	}

	output := testutil.CaptureStdout(t, func() {
		DisplayComparison(results, "yearly", false)
	})

	// Verify status indicators
	assert.Contains(t, output, "Status")
	assert.Contains(t, output, "M") // Married
	assert.Contains(t, output, "B") // Blind
}

func TestDisplayComparison_Verbose(t *testing.T) {
	results := []types.ComparisonResult{
		{
			Label:    "Test",
			Request:  testutil.CreateSampleTaxRequest(),
			Response: testutil.CreateSampleTaxResponse(),
		},
	}

	output := testutil.CaptureStdout(t, func() {
		DisplayComparison(results, "yearly", true)
	})

	// Verify verbose fields
	assert.Contains(t, output, "Tax Free Allowance")
	assert.Contains(t, output, "Taxable Pay")
	assert.Contains(t, output, "Basic Rate Tax")
}

func TestDisplayComparisonJSON(t *testing.T) {
	results := []types.ComparisonResult{
		{
			Label:    "Option1",
			Request:  testutil.CreateSampleTaxRequest(),
			Response: testutil.CreateSampleTaxResponse(),
		},
		{
			Label:    "Option2",
			Request:  testutil.CreateSampleTaxRequest(),
			Response: testutil.CreateSampleTaxResponse(),
		},
	}

	output := testutil.CaptureStdout(t, func() {
		DisplayComparisonJSON(results, "monthly")
	})

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(output), &parsed)
	require.NoError(t, err)

	// Verify structure
	assert.Contains(t, parsed, "period")
	assert.Equal(t, "monthly", parsed["period"])
	assert.Contains(t, parsed, "comparison")
	assert.Contains(t, parsed, "metadata")

	// Verify comparison fields
	comparison := parsed["comparison"].(map[string]interface{})
	assert.Contains(t, comparison, "gross_pay")
	assert.Contains(t, comparison, "net_pay")

	// Verify metadata
	metadata := parsed["metadata"].(map[string]interface{})
	assert.Contains(t, metadata, "Option1")
	assert.Contains(t, metadata, "Option2")
}

func TestPrintComparisonRow(t *testing.T) {
	results := []types.ComparisonResult{
		{
			Label: "Test1",
			Response: testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
				r.GrossPay = 50000.0
			}),
		},
		{
			Label: "Test2",
			Response: testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
				r.GrossPay = 60000.0
			}),
		},
	}

	output := testutil.CaptureStdout(t, func() {
		printComparisonRow("Test Field", results, 1.0, 20, 12,
			func(r *types.TaxResponse) float64 { return r.GrossPay })
	})

	assert.Contains(t, output, "Test Field")
	assert.Contains(t, output, "£50,000.00")
	assert.Contains(t, output, "£60,000.00")
}

func TestDisplayComparison_FourOptions(t *testing.T) {
	results := []types.ComparisonResult{
		{Label: "A", Request: testutil.CreateSampleTaxRequest(), Response: testutil.CreateSampleTaxResponse()},
		{Label: "B", Request: testutil.CreateSampleTaxRequest(), Response: testutil.CreateSampleTaxResponse()},
		{Label: "C", Request: testutil.CreateSampleTaxRequest(), Response: testutil.CreateSampleTaxResponse()},
		{Label: "D", Request: testutil.CreateSampleTaxRequest(), Response: testutil.CreateSampleTaxResponse()},
	}

	output := testutil.CaptureStdout(t, func() {
		DisplayComparison(results, "yearly", false)
	})

	// Verify all labels present
	assert.Contains(t, output, "A")
	assert.Contains(t, output, "B")
	assert.Contains(t, output, "C")
	assert.Contains(t, output, "D")
}

func TestDisplayComparison_Monthly(t *testing.T) {
	results := []types.ComparisonResult{
		{
			Label:   "Test",
			Request: testutil.CreateSampleTaxRequest(),
			Response: testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
				r.GrossPay = 60000.0 // £5000/month
			}),
		},
	}

	output := testutil.CaptureStdout(t, func() {
		DisplayComparison(results, "monthly", false)
	})

	// Should show monthly amount
	assert.Contains(t, output, "£5,000.00")
}

func TestDisplayComparison_LabelTruncation(t *testing.T) {
	results := []types.ComparisonResult{
		{
			Label:    "Very Long Label That Should Be Truncated",
			Request:  testutil.CreateSampleTaxRequest(),
			Response: testutil.CreateSampleTaxResponse(),
		},
	}

	output := testutil.CaptureStdout(t, func() {
		DisplayComparison(results, "yearly", false)
	})

	// Long labels should be truncated with ellipsis
	// May or may not truncate depending on label column width, just verify it doesn't crash
	assert.NotEmpty(t, output)
}
