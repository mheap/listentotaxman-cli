package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mheap/listentotaxman-cli/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_NoConfigFile(t *testing.T) {
	testutil.SetupViperTest(t)

	// Set HOME to a temp directory with no config file
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	})

	// Execute
	cfg, err := Load()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Verify default values
	assert.Equal(t, "uk", cfg.Defaults.Region)
	assert.Equal(t, "0", cfg.Defaults.Age)
	assert.Equal(t, "", cfg.Defaults.Pension)
	assert.Equal(t, "", cfg.Defaults.StudentLoan)
	assert.Equal(t, "", cfg.Defaults.TaxCode)
	assert.Equal(t, 0, cfg.Defaults.Extra)
	assert.Equal(t, "", cfg.Defaults.Year)
	assert.Equal(t, "yearly", cfg.Defaults.Period)
	assert.Equal(t, 0, cfg.Defaults.Income)
	assert.False(t, cfg.Defaults.Married)
	assert.False(t, cfg.Defaults.Blind)
	assert.False(t, cfg.Defaults.NoNI)
	assert.Equal(t, 0, cfg.Defaults.PartnerIncome)
}

func TestLoad_ValidConfigFile(t *testing.T) {
	testutil.SetupViperTest(t)

	// Create temp config
	testutil.CreateTempConfigFile(t, testutil.ValidConfigYAML)

	// Execute
	cfg, err := Load()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Verify all fields loaded correctly
	assert.Equal(t, "scotland", cfg.Defaults.Region)
	assert.Equal(t, "2024", cfg.Defaults.Year)
	assert.Equal(t, "30", cfg.Defaults.Age)
	assert.Equal(t, "5%", cfg.Defaults.Pension)
	assert.Equal(t, "plan2", cfg.Defaults.StudentLoan)
	assert.Equal(t, "1257L", cfg.Defaults.TaxCode)
	assert.Equal(t, 1000, cfg.Defaults.Extra)
	assert.Equal(t, "monthly", cfg.Defaults.Period)
	assert.Equal(t, 50000, cfg.Defaults.Income)
	assert.True(t, cfg.Defaults.Married)
	assert.False(t, cfg.Defaults.Blind)
	assert.False(t, cfg.Defaults.NoNI)
	assert.Equal(t, 25000, cfg.Defaults.PartnerIncome)
}

func TestLoad_PartialConfig(t *testing.T) {
	testutil.SetupViperTest(t)

	// Create temp config with only some fields
	testutil.CreateTempConfigFile(t, testutil.PartialConfigYAML)

	// Execute
	cfg, err := Load()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Verify specified fields loaded
	assert.Equal(t, "uk", cfg.Defaults.Region)
	assert.Equal(t, "2025", cfg.Defaults.Year)
	assert.Equal(t, "weekly", cfg.Defaults.Period)

	// Verify unspecified fields use defaults
	assert.Equal(t, "0", cfg.Defaults.Age)
	assert.Equal(t, "", cfg.Defaults.Pension)
	assert.Equal(t, 0, cfg.Defaults.Extra)
	assert.False(t, cfg.Defaults.Married)
}

func TestLoad_InvalidYAML(t *testing.T) {
	testutil.SetupViperTest(t)

	// Create temp config with invalid YAML
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "listentotaxman")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configFile := filepath.Join(configDir, "config.yaml")
	err = os.WriteFile(configFile, []byte(testutil.InvalidConfigYAML), 0644)
	require.NoError(t, err)

	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	})

	// Execute
	cfg, err := Load()

	// Assert - invalid YAML should return error
	require.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoad_AllDefaultFields(t *testing.T) {
	testutil.SetupViperTest(t)

	// Set HOME to a temp directory with no config file
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	})

	tests := []struct {
		name     string
		field    string
		expected interface{}
	}{
		{"region default", "region", "uk"},
		{"age default", "age", "0"},
		{"pension default", "pension", ""},
		{"student-loan default", "student-loan", ""},
		{"tax-code default", "tax-code", ""},
		{"extra default", "extra", 0},
		{"year default", "year", ""},
		{"period default", "period", "yearly"},
		{"income default", "income", 0},
		{"married default", "married", false},
		{"blind default", "blind", false},
		{"no-ni default", "no-ni", false},
		{"partner-income default", "partner-income", 0},
	}

	cfg, err := Load()
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.field {
			case "region":
				assert.Equal(t, tt.expected, cfg.Defaults.Region)
			case "age":
				assert.Equal(t, tt.expected, cfg.Defaults.Age)
			case "pension":
				assert.Equal(t, tt.expected, cfg.Defaults.Pension)
			case "student-loan":
				assert.Equal(t, tt.expected, cfg.Defaults.StudentLoan)
			case "tax-code":
				assert.Equal(t, tt.expected, cfg.Defaults.TaxCode)
			case "extra":
				assert.Equal(t, tt.expected, cfg.Defaults.Extra)
			case "year":
				assert.Equal(t, tt.expected, cfg.Defaults.Year)
			case "period":
				assert.Equal(t, tt.expected, cfg.Defaults.Period)
			case "income":
				assert.Equal(t, tt.expected, cfg.Defaults.Income)
			case "married":
				assert.Equal(t, tt.expected, cfg.Defaults.Married)
			case "blind":
				assert.Equal(t, tt.expected, cfg.Defaults.Blind)
			case "no-ni":
				assert.Equal(t, tt.expected, cfg.Defaults.NoNI)
			case "partner-income":
				assert.Equal(t, tt.expected, cfg.Defaults.PartnerIncome)
			}
		})
	}
}

func TestGetString_ExistsInConfig(t *testing.T) {
	testutil.SetupViperTest(t)

	// Create temp config
	testutil.CreateTempConfigFile(t, testutil.PartialConfigYAML)

	// Load config
	_, err := Load()
	require.NoError(t, err)

	// Test GetString with existing key
	result := GetString("defaults.region", "fallback")
	assert.Equal(t, "uk", result)
}

func TestGetString_NotInConfig(t *testing.T) {
	testutil.SetupViperTest(t)

	// Set HOME to a temp directory with no config file
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	})

	// Load config (no file exists)
	_, err := Load()
	require.NoError(t, err)

	// Test GetString with non-existent key
	result := GetString("nonexistent.key", "fallback")
	assert.Equal(t, "fallback", result)
}

func TestGetInt_WithFallback(t *testing.T) {
	tests := []struct {
		name     string
		config   string
		key      string
		fallback int
		expected int
	}{
		{
			name:     "exists in config",
			config:   testutil.ValidConfigYAML,
			key:      "defaults.extra",
			fallback: 999,
			expected: 1000,
		},
		{
			name:     "not in config uses default from viper",
			config:   testutil.PartialConfigYAML,
			key:      "defaults.extra",
			fallback: 999,
			expected: 0, // viper default is 0, not fallback
		},
		{
			name:     "key not set uses fallback",
			config:   testutil.PartialConfigYAML,
			key:      "nonexistent.key",
			fallback: 999,
			expected: 999, // Should use fallback when key doesn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.SetupViperTest(t)
			testutil.CreateTempConfigFile(t, tt.config)

			_, err := Load()
			require.NoError(t, err)

			result := GetInt(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}
