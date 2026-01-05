package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mheap/listentotaxman-cli/internal/client"
	"github.com/mheap/listentotaxman-cli/internal/testutil"
)

func TestRunCompare_TwoOptions(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	_ = os.Setenv("LISTENTOTAXMAN_CONFIG", configPath)
	t.Cleanup(func() { _ = os.Unsetenv("LISTENTOTAXMAN_CONFIG") })

	// Mock API client
	originalClientFactory := clientFactory
	t.Cleanup(func() { clientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	clientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Mock os.Args to simulate command-line args
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{
		"listentotaxman",
		"compare",
		"--option", "Job 1", "--income", "100000",
		"--option", "Job 2", "--income", "120000",
	}

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCompare(compareCmd, []string{})
		require.NoError(t, err)
	})

	// Verify output contains both job labels
	assert.Contains(t, output, "Job 1")
	assert.Contains(t, output, "Job 2")
	assert.Contains(t, output, "Net Pay")
}

func TestRunCompare_WithJSONFlag(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	_ = os.Setenv("LISTENTOTAXMAN_CONFIG", configPath)
	t.Cleanup(func() { _ = os.Unsetenv("LISTENTOTAXMAN_CONFIG") })

	// Mock API client
	originalClientFactory := clientFactory
	t.Cleanup(func() { clientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	clientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Mock os.Args
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{
		"listentotaxman",
		"compare",
		"--json",
		"--option", "Job 1", "--income", "100000",
		"--option", "Job 2", "--income", "120000",
	}

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCompare(compareCmd, []string{})
		require.NoError(t, err)
	})

	// Verify JSON output - labels are used as keys in the JSON structure
	assert.Contains(t, output, `"Job 1":`)
	assert.Contains(t, output, `"Job 2":`)
	assert.Contains(t, output, `"net_pay":`)
}

func TestRunCompare_WithPeriodFlag(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	_ = os.Setenv("LISTENTOTAXMAN_CONFIG", configPath)
	t.Cleanup(func() { _ = os.Unsetenv("LISTENTOTAXMAN_CONFIG") })

	// Mock API client
	originalClientFactory := clientFactory
	t.Cleanup(func() { clientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	clientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Mock os.Args
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{
		"listentotaxman",
		"compare",
		"--period", "monthly",
		"--option", "Job 1", "--income", "100000",
		"--option", "Job 2", "--income", "120000",
	}

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCompare(compareCmd, []string{})
		require.NoError(t, err)
	})

	// Verify output
	assert.Contains(t, output, "Job 1")
	assert.Contains(t, output, "Job 2")
}

func TestRunCompare_LessThanTwoOptions(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	_ = os.Setenv("LISTENTOTAXMAN_CONFIG", configPath)
	t.Cleanup(func() { _ = os.Unsetenv("LISTENTOTAXMAN_CONFIG") })

	// Mock os.Args
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{
		"listentotaxman",
		"compare",
		"--option", "Job 1", "--income", "100000",
	}

	err := runCompare(compareCmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 2 options required")
}

func TestRunCompare_MoreThanFourOptions(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	_ = os.Setenv("LISTENTOTAXMAN_CONFIG", configPath)
	t.Cleanup(func() { _ = os.Unsetenv("LISTENTOTAXMAN_CONFIG") })

	// Mock os.Args
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{
		"listentotaxman",
		"compare",
		"--option", "Job 1", "--income", "100000",
		"--option", "Job 2", "--income", "110000",
		"--option", "Job 3", "--income", "120000",
		"--option", "Job 4", "--income", "130000",
		"--option", "Job 5", "--income", "140000",
	}

	err := runCompare(compareCmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum 4 options supported")
}

func TestRunCompare_InvalidPeriod(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	_ = os.Setenv("LISTENTOTAXMAN_CONFIG", configPath)
	t.Cleanup(func() { _ = os.Unsetenv("LISTENTOTAXMAN_CONFIG") })

	// Mock os.Args
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{
		"listentotaxman",
		"compare",
		"--period", "invalid",
		"--option", "Job 1", "--income", "100000",
		"--option", "Job 2", "--income", "120000",
	}

	err := runCompare(compareCmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid period")
}

func TestRunCompare_ValidationError(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	_ = os.Setenv("LISTENTOTAXMAN_CONFIG", configPath)
	t.Cleanup(func() { _ = os.Unsetenv("LISTENTOTAXMAN_CONFIG") })

	// Mock os.Args with invalid income (zero)
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{
		"listentotaxman",
		"compare",
		"--option", "Job 1", "--income", "0",
		"--option", "Job 2", "--income", "120000",
	}

	err := runCompare(compareCmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "income must be greater than 0")
}

func TestRunCompare_APIError(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	_ = os.Setenv("LISTENTOTAXMAN_CONFIG", configPath)
	t.Cleanup(func() { _ = os.Unsetenv("LISTENTOTAXMAN_CONFIG") })

	// Mock API client with error
	originalClientFactory := clientFactory
	t.Cleanup(func() { clientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperError(500, "Internal Server Error")
	clientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Mock os.Args
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{
		"listentotaxman",
		"compare",
		"--option", "Job 1", "--income", "100000",
		"--option", "Job 2", "--income", "120000",
	}

	err := runCompare(compareCmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to calculate tax")
}

func TestRunCompare_HelpFlag(t *testing.T) {
	// Mock os.Args with --help
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{
		"listentotaxman",
		"compare",
		"--help",
	}

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCompare(compareCmd, []string{})
		// Help doesn't return an error
		require.NoError(t, err)
	})

	// Verify help text is displayed
	assert.Contains(t, output, "Compare tax calculations")
	assert.Contains(t, output, "Usage:")
}

func TestRunCompare_FourOptionsSuccess(t *testing.T) {
	// Setup viper test
	testutil.SetupViperTest(t)

	// Create config file
	configPath := testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)
	_ = os.Setenv("LISTENTOTAXMAN_CONFIG", configPath)
	t.Cleanup(func() { _ = os.Unsetenv("LISTENTOTAXMAN_CONFIG") })

	// Mock API client
	originalClientFactory := clientFactory
	t.Cleanup(func() { clientFactory = originalClientFactory })

	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	clientFactory = func() *client.Client {
		return client.NewWithHTTPClient(&http.Client{Transport: mockRT})
	}

	// Mock os.Args
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{
		"listentotaxman",
		"compare",
		"--option", "Job 1", "--income", "100000",
		"--option", "Job 2", "--income", "110000",
		"--option", "Job 3", "--income", "120000",
		"--option", "Job 4", "--income", "130000",
	}

	// Capture stdout
	output := testutil.CaptureStdout(t, func() {
		err := runCompare(compareCmd, []string{})
		require.NoError(t, err)
	})

	// Verify all four jobs are in output
	assert.Contains(t, output, "Job 1")
	assert.Contains(t, output, "Job 2")
	assert.Contains(t, output, "Job 3")
	assert.Contains(t, output, "Job 4")
}
