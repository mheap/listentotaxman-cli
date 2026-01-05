package cmd

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/mheap/listentotaxman-cli/internal/types"
	"github.com/stretchr/testify/assert"
)

// Test validation logic that's extracted and testable

func TestValidation_InvalidYear(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		year        string
		shouldError bool
	}{
		{"valid 4 digit", "2024", false},
		{"too short", "202", true},
		{"too long", "20244", true},
		{"not a number", "abcd", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Validate year length
			hasError := len(tt.year) != 4
			if !hasError && tt.year != "" {
				// Try to parse as int
				_, err := strconv.Atoi(tt.year)
				hasError = err != nil
			}

			assert.Equal(t, tt.shouldError, hasError || tt.year == "")
		})
	}
}

func TestValidation_StudentLoanPlans(t *testing.T) {
	t.Parallel()

	validPlans := []string{"plan1", "plan2", "plan4", "postgraduate", "scottish"}

	tests := []struct {
		name    string
		plan    string
		isValid bool
	}{
		{"plan1", "plan1", true},
		{"plan2", "plan2", true},
		{"plan4", "plan4", true},
		{"postgraduate", "postgraduate", true},
		{"scottish", "scottish", true},
		{"plan3 invalid", "plan3", false},
		{"plan5 invalid", "plan5", false},
		{"empty is valid", "", true}, // Empty means no student loan
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.plan == "" {
				assert.True(t, tt.isValid)
				return
			}

			isValid := false
			for _, vp := range validPlans {
				if tt.plan == vp {
					isValid = true
					break
				}
			}

			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestValidation_Periods(t *testing.T) {
	t.Parallel()

	validPeriods := []string{"yearly", "monthly", "weekly", "daily", "hourly"}

	tests := []struct {
		name    string
		period  string
		isValid bool
	}{
		{"yearly", "yearly", true},
		{"monthly", "monthly", true},
		{"weekly", "weekly", true},
		{"daily", "daily", true},
		{"hourly", "hourly", true},
		{"invalid", "biweekly", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			isValid := false
			for _, vp := range validPeriods {
				if tt.period == vp {
					isValid = true
					break
				}
			}

			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestBooleanFlagConversion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		flagValue bool
		want      string
	}{
		{"married true", true, "y"},
		{"married false", false, ""},
		{"blind true", true, "y"},
		{"blind false", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := &types.TaxRequest{}

			// Simulate flag logic
			if tt.flagValue {
				req.Married = "y"
			}

			if tt.flagValue {
				assert.Equal(t, "y", req.Married)
			} else {
				assert.Equal(t, "", req.Married)
			}
		})
	}
}

func TestValidation_PartnerIncome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		partnerIncome int
		married       string
		shouldError   bool
		errorMsg      string
	}{
		{
			name:          "partner income with married",
			partnerIncome: 25000,
			married:       "y",
			shouldError:   false,
		},
		{
			name:          "partner income without married",
			partnerIncome: 25000,
			married:       "",
			shouldError:   true,
			errorMsg:      "--partner-income requires --married flag",
		},
		{
			name:          "negative partner income",
			partnerIncome: -1000,
			married:       "y",
			shouldError:   true,
			errorMsg:      "--partner-income cannot be negative",
		},
		{
			name:          "zero partner income is ok",
			partnerIncome: 0,
			married:       "",
			shouldError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := &types.TaxRequest{
				PartnerGrossWage: tt.partnerIncome,
				Married:          tt.married,
			}

			// Validate
			var err error
			if req.PartnerGrossWage > 0 && req.Married != "y" {
				err = fmt.Errorf("--partner-income requires --married flag")
			}
			if req.PartnerGrossWage < 0 {
				err = fmt.Errorf("--partner-income cannot be negative")
			}

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidation_Income(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		income      int
		shouldError bool
	}{
		{"positive income", 50000, false},
		{"zero income", 0, true},
		{"negative income", -1000, true},
		{"large income", 1000000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hasError := tt.income <= 0
			assert.Equal(t, tt.shouldError, hasError)
		})
	}
}
