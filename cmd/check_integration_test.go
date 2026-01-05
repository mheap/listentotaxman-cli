package cmd

import (
	"net/http"
	"testing"

	"github.com/mheap/listentotaxman-cli/internal/client"
	"github.com/mheap/listentotaxman-cli/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCheck_BasicIncome(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	t.Setenv("LISTENTOTAXMAN_CONFIG", configPath)

	// Mock API client
	originalClientFactory := checkClientFactory
	t.Cleanup(func() { checkClientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	checkClientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Set flags
	flagIncome = 100000
	flagYear = ""
	flagRegion = ""
	flagAge = ""
	flagPension = ""
	flagStudentLoan = ""
	flagExtra = 0
	flagTaxCode = ""
	flagJSON = false
	flagVerbose = false
	flagPeriod = ""
	flagMarried = false
	flagBlind = false
	flagNoNI = false
	flagPartnerIncome = 0

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCheck(checkCmd, []string{})
		require.NoError(t, err)
	})

	// Verify output
	assert.Contains(t, output, "Gross Salary")
	assert.Contains(t, output, "Net Pay")
	assert.Contains(t, output, "Tax Paid")
}

func TestRunCheck_WithJSONOutput(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	t.Setenv("LISTENTOTAXMAN_CONFIG", configPath)

	// Mock API client
	originalClientFactory := checkClientFactory
	t.Cleanup(func() { checkClientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	checkClientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Set flags
	flagIncome = 100000
	flagJSON = true
	flagVerbose = false
	flagPeriod = "yearly"

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCheck(checkCmd, []string{})
		require.NoError(t, err)
	})

	// Verify JSON output
	assert.Contains(t, output, `"gross_pay":`)
	assert.Contains(t, output, `"net_pay":`)
	assert.Contains(t, output, `"tax_paid":`)
}

func TestRunCheck_WithVerboseOutput(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	t.Setenv("LISTENTOTAXMAN_CONFIG", configPath)

	// Mock API client
	originalClientFactory := checkClientFactory
	t.Cleanup(func() { checkClientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	checkClientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Set flags
	flagIncome = 100000
	flagJSON = false
	flagVerbose = true
	flagPeriod = "yearly"

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCheck(checkCmd, []string{})
		require.NoError(t, err)
	})

	// Verify detailed output contains tax brackets
	assert.Contains(t, output, "Gross Salary")
	assert.Contains(t, output, "Net Pay")
}

func TestRunCheck_WithAllFlags(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	t.Setenv("LISTENTOTAXMAN_CONFIG", configPath)

	// Mock API client
	originalClientFactory := checkClientFactory
	t.Cleanup(func() { checkClientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponseWithStudentLoan)
	checkClientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Set all flags
	flagIncome = 100000
	flagYear = "2024"
	flagRegion = "scotland"
	flagAge = "45"
	flagPension = "5%"
	flagStudentLoan = "plan2"
	flagExtra = 2000
	flagTaxCode = "1257L"
	flagMarried = true
	flagBlind = true
	flagNoNI = false
	flagPartnerIncome = 25000
	flagJSON = false
	flagVerbose = false
	flagPeriod = "yearly"

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCheck(checkCmd, []string{})
		require.NoError(t, err)
	})

	// Verify output
	assert.Contains(t, output, "Gross Salary")
	assert.Contains(t, output, "Net Pay")
}

func TestRunCheck_MissingIncome(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	t.Setenv("LISTENTOTAXMAN_CONFIG", configPath)

	// Set flags with missing income
	flagIncome = 0
	flagYear = ""
	flagRegion = ""
	flagAge = ""
	flagPension = ""
	flagStudentLoan = ""
	flagExtra = 0
	flagTaxCode = ""

	err := runCheck(checkCmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "income must be greater than 0")
}

func TestRunCheck_InvalidPeriod(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	t.Setenv("LISTENTOTAXMAN_CONFIG", configPath)

	// Set flags with invalid period
	flagIncome = 100000
	flagPeriod = "invalid"

	err := runCheck(checkCmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid period")
}

func TestRunCheck_APIError(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	t.Setenv("LISTENTOTAXMAN_CONFIG", configPath)

	// Mock API client with error
	originalClientFactory := checkClientFactory
	t.Cleanup(func() { checkClientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperError(500, "Internal Server Error")
	checkClientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Reset and set flags
	flagIncome = 100000
	flagYear = ""
	flagRegion = ""
	flagAge = ""
	flagPension = ""
	flagStudentLoan = ""
	flagExtra = 0
	flagTaxCode = ""
	flagJSON = false
	flagVerbose = false
	flagPeriod = ""
	flagMarried = false
	flagBlind = false
	flagNoNI = false
	flagPartnerIncome = 0

	err := runCheck(checkCmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to calculate tax")
}

func TestRunCheck_MonthlyPeriod(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	t.Setenv("LISTENTOTAXMAN_CONFIG", configPath)

	// Mock API client
	originalClientFactory := checkClientFactory
	t.Cleanup(func() { checkClientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	checkClientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Set flags
	flagIncome = 100000
	flagPeriod = "monthly"
	flagJSON = false
	flagVerbose = false

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCheck(checkCmd, []string{})
		require.NoError(t, err)
	})

	// Verify output
	assert.Contains(t, output, "Gross Salary")
	assert.Contains(t, output, "Net Pay")
}

func TestRunCheck_WithConfigDefaults(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file with defaults
	configYAML := `defaults:
  region: scotland
  age: "30"
  pension: "3%"
  student_loan: plan1
  year: "2024"
`
	configPath := testutil.CreateTempConfigFile(t, configYAML)
	t.Setenv("LISTENTOTAXMAN_CONFIG", configPath)

	// Mock API client
	originalClientFactory := checkClientFactory
	t.Cleanup(func() { checkClientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponseWithStudentLoan)
	checkClientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Set only income (should use config defaults for others)
	flagIncome = 100000
	flagYear = ""
	flagRegion = ""
	flagAge = ""
	flagPension = ""
	flagStudentLoan = ""
	flagJSON = false
	flagVerbose = false
	flagPeriod = ""

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCheck(checkCmd, []string{})
		require.NoError(t, err)
	})

	// Verify output
	assert.Contains(t, output, "Gross Salary")
	assert.Contains(t, output, "Net Pay")
}
