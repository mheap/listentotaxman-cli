package testutil

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mheap/listentotaxman-cli/internal/types"
)

var updateGolden = flag.Bool("update-golden", false, "update golden files")

// CreateTempConfigFile creates a temporary config file for testing
func CreateTempConfigFile(t *testing.T, content string) string {
	t.Helper()

	// Create temp directory
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "listentotaxman")
	err := os.MkdirAll(configDir, 0755) //nolint:gosec // Test files only
	require.NoError(t, err, "failed to create temp config directory")

	// Write config file
	configFile := filepath.Join(configDir, "config.yaml")
	err = os.WriteFile(configFile, []byte(content), 0644) //nolint:gosec // Test files only
	require.NoError(t, err, "failed to write temp config file")

	// Update HOME and USERPROFILE to point to temp directory for config loading
	// HOME is used on Unix/macOS, USERPROFILE is used on Windows
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	require.NoError(t, os.Setenv("HOME", tempDir))
	require.NoError(t, os.Setenv("USERPROFILE", tempDir))
	t.Cleanup(func() {
		_ = os.Setenv("HOME", originalHome)
		_ = os.Setenv("USERPROFILE", originalUserProfile)
	})

	return configDir
}

// CaptureStdout captures stdout during function execution
func CaptureStdout(t *testing.T, f func()) string {
	t.Helper()

	// Save original stdout
	original := os.Stdout

	// Create pipe to capture output
	r, w, err := os.Pipe()
	require.NoError(t, err, "failed to create pipe")

	// Replace stdout
	os.Stdout = w

	// Ensure stdout is restored
	defer func() {
		os.Stdout = original
	}()

	// Run function that writes to stdout
	done := make(chan struct{})
	var buf bytes.Buffer

	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	f()

	// Close writer and wait for reader to finish
	_ = w.Close()
	<-done

	return buf.String()
}

// CreateSampleTaxRequest creates a valid TaxRequest with optional overrides
func CreateSampleTaxRequest(overrides ...func(*types.TaxRequest)) *types.TaxRequest {
	req := &types.TaxRequest{
		Response:  "json",
		Year:      "2024",
		TaxRegion: "uk",
		Age:       "30",
		Pension:   "",
		Time:      "1",
		GrossWage: 50000,
		Plan:      "",
		Extra:     0,
		TaxCode:   "",
	}

	for _, override := range overrides {
		override(req)
	}

	return req
}

// CreateSampleTaxResponse creates a valid TaxResponse with optional overrides
func CreateSampleTaxResponse(overrides ...func(*types.TaxResponse)) *types.TaxResponse {
	resp := &types.TaxResponse{
		TaxYear:          2024,
		TaxablePay:       37430.0,
		GrossPay:         50000.0,
		AdditionalGross:  0.0,
		TaxFreeAllowance: 12570.0,
		TaxPaid:          7486.0,
		TaxDue: map[string]types.TaxBracket{
			"0": {
				Rate:   0.20,
				Amount: 7486.0,
			},
		},
		NationalInsurance:        4218.16,
		NetPay:                   38295.84,
		StudentLoanRepayment:     0.0,
		PensionHMRC:              0.0,
		PensionYou:               0.0,
		PensionClaimback:         0.0,
		EmployersNI:              5220.78,
		TaxFreeMarried:           0.0,
		TaxRegion:                "uk",
		TaxCode:                  "1257L",
		TaxFreeMarriageAllowance: 0.0,
		GrossSacrifice:           0.0,
		ChildcareAmount:          0.0,
	}

	for _, override := range overrides {
		override(resp)
	}

	return resp
}

// SetupViperTest resets viper state for isolated testing
func SetupViperTest(t *testing.T) {
	t.Helper()
	viper.Reset()
	t.Cleanup(func() {
		viper.Reset()
	})
}

// AssertNoError is a convenience wrapper for require.NoError
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.NoError(t, err, msgAndArgs...)
}

// AssertError checks that an error occurred and contains the expected message
func AssertError(t *testing.T, err error, msgContains string) {
	t.Helper()
	require.Error(t, err, "expected an error but got nil")
	assert.Contains(t, err.Error(), msgContains, "error message should contain expected text")
}

// AssertJSONEqual compares two JSON strings for equality
func AssertJSONEqual(t *testing.T, expected, actual string) {
	t.Helper()

	var expectedJSON, actualJSON interface{}

	err := json.Unmarshal([]byte(expected), &expectedJSON)
	require.NoError(t, err, "failed to unmarshal expected JSON")

	err = json.Unmarshal([]byte(actual), &actualJSON)
	require.NoError(t, err, "failed to unmarshal actual JSON")

	assert.Equal(t, expectedJSON, actualJSON, "JSON structures should be equal")
}

// CompareGoldenFile compares actual output with a golden file
func CompareGoldenFile(t *testing.T, name string, actual string) {
	t.Helper()

	goldenPath := filepath.Join(".goldenfiles", name+".txt")

	if *updateGolden {
		// Update golden file
		err := os.WriteFile(goldenPath, []byte(actual), 0644) //nolint:gosec // Test golden files
		require.NoError(t, err, "failed to update golden file")
		t.Logf("Updated golden file: %s", goldenPath)
		return
	}

	// Compare with golden file
	expected, err := os.ReadFile(goldenPath) //nolint:gosec // Test golden files
	require.NoError(t, err, "failed to read golden file: %s (run with -update-golden to create)", goldenPath)

	assert.Equal(t, string(expected), actual, "output should match golden file: %s", goldenPath)
}
