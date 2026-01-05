package cmd

import (
	"testing"

	"github.com/mheap/listentotaxman-cli/internal/config"
	"github.com/mheap/listentotaxman-cli/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateOption_ValidOption(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Valid Job",
		Request: &types.TaxRequest{
			Year:      "2024",
			GrossWage: 100000,
		},
	}

	err := validateOption(&opt, cfg)
	assert.NoError(t, err)
}

func TestValidateOption_InvalidYearLength(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Job 1",
		Request: &types.TaxRequest{
			Year:      "24", // Too short
			GrossWage: 100000,
		},
	}

	err := validateOption(&opt, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "year must be a 4-digit number")
	assert.Contains(t, err.Error(), "Job 1")
}

func TestValidateOption_InvalidYearNotNumber(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Job 1",
		Request: &types.TaxRequest{
			Year:      "abcd",
			GrossWage: 100000,
		},
	}

	err := validateOption(&opt, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "year must be a valid number")
	assert.Contains(t, err.Error(), "Job 1")
}

func TestValidateOption_ZeroIncome(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Job 1",
		Request: &types.TaxRequest{
			Year:      "2024",
			GrossWage: 0,
		},
	}

	err := validateOption(&opt, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "income must be greater than 0")
	assert.Contains(t, err.Error(), "Job 1")
}

func TestValidateOption_NegativeIncome(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Job 1",
		Request: &types.TaxRequest{
			Year:      "2024",
			GrossWage: -100000,
		},
	}

	err := validateOption(&opt, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "income must be greater than 0")
}

func TestValidateOption_PartnerIncomeWithoutMarried(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Job 1",
		Request: &types.TaxRequest{
			Year:             "2024",
			GrossWage:        100000,
			PartnerGrossWage: 25000,
			Married:          "", // Not married
		},
	}

	err := validateOption(&opt, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--partner-income requires --married flag")
	assert.Contains(t, err.Error(), "Job 1")
	assert.Contains(t, err.Error(), "25000")
}

func TestValidateOption_PartnerIncomeWithMarried(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Job 1",
		Request: &types.TaxRequest{
			Year:             "2024",
			GrossWage:        100000,
			PartnerGrossWage: 25000,
			Married:          "y",
		},
	}

	err := validateOption(&opt, cfg)
	assert.NoError(t, err)
}

func TestValidateOption_NegativePartnerIncome(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Job 1",
		Request: &types.TaxRequest{
			Year:             "2024",
			GrossWage:        100000,
			PartnerGrossWage: -5000,
			Married:          "y",
		},
	}

	err := validateOption(&opt, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--partner-income cannot be negative")
	assert.Contains(t, err.Error(), "Job 1")
}

func TestValidateOption_ValidStudentLoanPlans(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	validPlans := []string{"plan1", "plan2", "plan4", "postgraduate", "scottish"}

	for _, plan := range validPlans {
		t.Run(plan, func(t *testing.T) {
			opt := ComparisonOption{
				Label: "Job 1",
				Request: &types.TaxRequest{
					Year:      "2024",
					GrossWage: 100000,
					Plan:      plan,
				},
			}

			err := validateOption(&opt, cfg)
			assert.NoError(t, err)
		})
	}
}

func TestValidateOption_InvalidStudentLoanPlan(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Job 1",
		Request: &types.TaxRequest{
			Year:      "2024",
			GrossWage: 100000,
			Plan:      "plan3", // Invalid
		},
	}

	err := validateOption(&opt, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid student loan plan")
	assert.Contains(t, err.Error(), "plan3")
	assert.Contains(t, err.Error(), "Job 1")
	assert.Contains(t, err.Error(), "plan1, plan2, plan4, postgraduate, scottish")
}

func TestValidateOption_EmptyStudentLoan(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Job 1",
		Request: &types.TaxRequest{
			Year:      "2024",
			GrossWage: 100000,
			Plan:      "", // Empty is valid (no student loan)
		},
	}

	err := validateOption(&opt, cfg)
	assert.NoError(t, err)
}

func TestValidateOption_AllFieldsValid(t *testing.T) {
	cfg := &config.Config{
		Defaults: config.Defaults{},
	}

	opt := ComparisonOption{
		Label: "Complex Job",
		Request: &types.TaxRequest{
			Year:             "2024",
			GrossWage:        100000,
			TaxRegion:        "scotland",
			Age:              "45",
			Pension:          "5%",
			Plan:             "plan2",
			Extra:            2000,
			TaxCode:          "1257L",
			Married:          "y",
			Blind:            "y",
			ExNI:             "y",
			PartnerGrossWage: 25000,
		},
	}

	err := validateOption(&opt, cfg)
	assert.NoError(t, err)
}
