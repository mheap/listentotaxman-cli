package cmd

import (
	"testing"
	"time"

	"github.com/mheap/listentotaxman-cli/internal/testutil"
	"github.com/mheap/listentotaxman-cli/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultYear(t *testing.T) {
	// Save original and restore after all tests
	originalTimeNow := timeNowFunc
	defer func() { timeNowFunc = originalTimeNow }()

	tests := []struct {
		name     string
		mockTime time.Time
		want     string
	}{
		{
			name:     "Before April 5th",
			mockTime: time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC),
			want:     "2025",
		},
		{
			name:     "On April 5th at midnight",
			mockTime: time.Date(2026, time.April, 5, 0, 0, 0, 0, time.UTC),
			want:     "2025",
		},
		{
			name:     "On April 5th at 11:59 PM - after midnight counts as after",
			mockTime: time.Date(2026, time.April, 5, 23, 59, 59, 0, time.UTC),
			want:     "2026", // After midnight on April 5th
		},
		{
			name:     "After April 5th - April 6th (new tax year)",
			mockTime: time.Date(2026, time.April, 6, 0, 0, 0, 0, time.UTC),
			want:     "2026",
		},
		{
			name:     "December",
			mockTime: time.Date(2026, time.December, 31, 0, 0, 0, 0, time.UTC),
			want:     "2026",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock time for this specific test
			timeNowFunc = func() time.Time {
				return tt.mockTime
			}

			got := getDefaultYear()
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
		{"invalid defaults to yearly", "invalid", 1.0},
		{"empty defaults to yearly", "", 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getPeriodDivisor(tt.period)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNormalizeRegion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		region string
		want   string
	}{
		{"england becomes uk", "england", "uk"},
		{"uk stays uk", "uk", "uk"},
		{"scotland stays scotland", "scotland", "scotland"},
		{"wales stays wales", "wales", "wales"},
		{"northern-ireland stays", "northern-ireland", "northern-ireland"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := normalizeRegion(tt.region)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAdjustResponseForPeriod_Yearly(t *testing.T) {
	t.Parallel()

	resp := testutil.CreateSampleTaxResponse()
	original := *resp

	adjusted := adjustResponseForPeriod(resp, "yearly")

	// Should return the same response (pointer may differ, but values same)
	assert.Equal(t, original.GrossPay, adjusted.GrossPay)
	assert.Equal(t, original.TaxPaid, adjusted.TaxPaid)
	assert.Equal(t, original.NationalInsurance, adjusted.NationalInsurance)
	assert.Equal(t, original.NetPay, adjusted.NetPay)
}

func TestAdjustResponseForPeriod_Monthly(t *testing.T) {
	t.Parallel()

	resp := testutil.CreateSampleTaxResponse()
	original := *resp

	adjusted := adjustResponseForPeriod(resp, "monthly")

	// All monetary values should be divided by 12
	assert.InDelta(t, original.GrossPay/12.0, adjusted.GrossPay, 0.01)
	assert.InDelta(t, original.TaxablePay/12.0, adjusted.TaxablePay, 0.01)
	assert.InDelta(t, original.TaxPaid/12.0, adjusted.TaxPaid, 0.01)
	assert.InDelta(t, original.NationalInsurance/12.0, adjusted.NationalInsurance, 0.01)
	assert.InDelta(t, original.NetPay/12.0, adjusted.NetPay, 0.01)
	assert.InDelta(t, original.TaxFreeAllowance/12.0, adjusted.TaxFreeAllowance, 0.01)
}

func TestAdjustResponseForPeriod_Weekly(t *testing.T) {
	t.Parallel()

	resp := testutil.CreateSampleTaxResponse()
	original := *resp

	adjusted := adjustResponseForPeriod(resp, "weekly")

	// Spot check - divided by 52
	assert.InDelta(t, original.GrossPay/52.0, adjusted.GrossPay, 0.01)
	assert.InDelta(t, original.NetPay/52.0, adjusted.NetPay, 0.01)
}

func TestAdjustResponseForPeriod_Hourly(t *testing.T) {
	t.Parallel()

	resp := testutil.CreateSampleTaxResponse()
	original := *resp

	adjusted := adjustResponseForPeriod(resp, "hourly")

	// Divided by 2080
	assert.InDelta(t, original.GrossPay/2080.0, adjusted.GrossPay, 0.01)
	assert.InDelta(t, original.NetPay/2080.0, adjusted.NetPay, 0.01)
}

func TestAdjustResponseForPeriod_TaxBrackets(t *testing.T) {
	t.Parallel()

	resp := testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
		r.TaxDue = map[string]types.TaxBracket{
			"0": {Rate: 0.20, Amount: 7486.0},
			"1": {Rate: 0.40, Amount: 2000.0},
			"2": {Rate: 0.45, Amount: 500.0},
		}
	})

	adjusted := adjustResponseForPeriod(resp, "monthly")

	// Tax brackets should be adjusted
	assert.InDelta(t, 7486.0/12.0, adjusted.TaxDue["0"].Amount, 0.01)
	assert.InDelta(t, 2000.0/12.0, adjusted.TaxDue["1"].Amount, 0.01)
	assert.InDelta(t, 500.0/12.0, adjusted.TaxDue["2"].Amount, 0.01)

	// Rates should be unchanged
	assert.Equal(t, 0.20, adjusted.TaxDue["0"].Rate)
	assert.Equal(t, 0.40, adjusted.TaxDue["1"].Rate)
	assert.Equal(t, 0.45, adjusted.TaxDue["2"].Rate)
}

func TestAdjustResponseForPeriod_NestedPrevious(t *testing.T) {
	t.Parallel()

	previousResp := testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
		r.TaxYear = 2023
		r.GrossPay = 45000.0
	})

	resp := testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
		r.Previous = previousResp
	})

	adjusted := adjustResponseForPeriod(resp, "monthly")

	// Verify previous response was also adjusted
	assert.NotNil(t, adjusted.Previous)
	assert.InDelta(t, 45000.0/12.0, adjusted.Previous.GrossPay, 0.01)
}

func TestAdjustResponseForPeriod_AllFields(t *testing.T) {
	t.Parallel()

	resp := testutil.CreateSampleTaxResponse(func(r *types.TaxResponse) {
		r.TaxablePay = 40000.0
		r.GrossPay = 50000.0
		r.AdditionalGross = 1000.0
		r.TaxFreeAllowance = 12570.0
		r.TaxPaid = 7500.0
		r.NationalInsurance = 4200.0
		r.NetPay = 38300.0
		r.StudentLoanRepayment = 450.0
		r.PensionHMRC = 1000.0
		r.PensionYou = 2000.0
		r.PensionClaimback = 200.0
		r.EmployersNI = 5200.0
		r.TaxFreeMarried = 1260.0
		r.TaxFreeMarriageAllowance = 1260.0
		r.GrossSacrifice = 500.0
		r.ChildcareAmount = 100.0
	})

	adjusted := adjustResponseForPeriod(resp, "monthly")
	divisor := 12.0

	// Verify all monetary fields are adjusted
	assert.InDelta(t, 40000.0/divisor, adjusted.TaxablePay, 0.01)
	assert.InDelta(t, 50000.0/divisor, adjusted.GrossPay, 0.01)
	assert.InDelta(t, 1000.0/divisor, adjusted.AdditionalGross, 0.01)
	assert.InDelta(t, 12570.0/divisor, adjusted.TaxFreeAllowance, 0.01)
	assert.InDelta(t, 7500.0/divisor, adjusted.TaxPaid, 0.01)
	assert.InDelta(t, 4200.0/divisor, adjusted.NationalInsurance, 0.01)
	assert.InDelta(t, 38300.0/divisor, adjusted.NetPay, 0.01)
	assert.InDelta(t, 450.0/divisor, adjusted.StudentLoanRepayment, 0.01)
	assert.InDelta(t, 1000.0/divisor, adjusted.PensionHMRC, 0.01)
	assert.InDelta(t, 2000.0/divisor, adjusted.PensionYou, 0.01)
	assert.InDelta(t, 200.0/divisor, adjusted.PensionClaimback, 0.01)
	assert.InDelta(t, 5200.0/divisor, adjusted.EmployersNI, 0.01)
	assert.InDelta(t, 1260.0/divisor, adjusted.TaxFreeMarried, 0.01)
	assert.InDelta(t, 1260.0/divisor, adjusted.TaxFreeMarriageAllowance, 0.01)
	assert.InDelta(t, 500.0/divisor, adjusted.GrossSacrifice, 0.01)
	assert.InDelta(t, 100.0/divisor, adjusted.ChildcareAmount, 0.01)
}
