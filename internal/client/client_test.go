package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/mheap/listentotaxman-cli/internal/testutil"
	"github.com/mheap/listentotaxman-cli/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateTax_Success(t *testing.T) {
	t.Parallel()

	// Create mock HTTP client
	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	client := &Client{
		httpClient: &http.Client{
			Transport: mockRT,
		},
	}

	// Create request
	req := testutil.CreateSampleTaxRequest()

	// Execute
	resp, err := client.CalculateTax(req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2024, resp.TaxYear)
	assert.Equal(t, 50000.0, resp.GrossPay)
	assert.Equal(t, 7486.0, resp.TaxPaid)
	assert.Equal(t, 4218.16, resp.NationalInsurance)
	assert.Equal(t, 38295.84, resp.NetPay)
	assert.Equal(t, "uk", resp.TaxRegion)
	assert.Equal(t, "1257L", resp.TaxCode)
}

func TestCalculateTax_APIError(t *testing.T) {
	t.Parallel()

	// Create mock HTTP client that returns 500
	mockRT := testutil.NewMockRoundTripperError(500, testutil.SampleAPIResponse500)
	client := &Client{
		httpClient: &http.Client{
			Transport: mockRT,
		},
	}

	// Create request
	req := testutil.CreateSampleTaxRequest()

	// Execute
	resp, err := client.CalculateTax(req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "API returned status 500")
	assert.Contains(t, err.Error(), "Internal server error")
}

func TestCalculateTax_InvalidJSON(t *testing.T) {
	t.Parallel()

	// Create mock HTTP client with malformed JSON
	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponseInvalidJSON)
	client := &Client{
		httpClient: &http.Client{
			Transport: mockRT,
		},
	}

	// Create request
	req := testutil.CreateSampleTaxRequest()

	// Execute
	resp, err := client.CalculateTax(req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to parse response")
}

func TestCalculateTax_NetworkError(t *testing.T) {
	t.Parallel()

	// Create mock HTTP client that returns network error
	networkErr := errors.New("connection refused")
	mockRT := testutil.NewMockRoundTripperNetworkError(networkErr)
	client := &Client{
		httpClient: &http.Client{
			Transport: mockRT,
		},
	}

	// Create request
	req := testutil.CreateSampleTaxRequest()

	// Execute
	resp, err := client.CalculateTax(req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to execute request")
	assert.Contains(t, err.Error(), "connection refused")
}

func TestCalculateTax_SetsRequiredFields(t *testing.T) {
	t.Parallel()

	// Create mock HTTP client
	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	client := &Client{
		httpClient: &http.Client{
			Transport: mockRT,
		},
	}

	// Create request without Response and Time fields
	req := &types.TaxRequest{
		Year:      "2024",
		TaxRegion: "uk",
		Age:       "30",
		GrossWage: 50000,
	}

	// Execute
	_, err := client.CalculateTax(req)

	// Assert
	require.NoError(t, err)

	// Get request body and verify fields were set
	body, err := mockRT.GetRequestBody()
	require.NoError(t, err)

	var sentReq types.TaxRequest
	err = json.Unmarshal([]byte(body), &sentReq)
	require.NoError(t, err)

	assert.Equal(t, "json", sentReq.Response, "Response field should be set to 'json'")
	assert.Equal(t, "1", sentReq.Time, "Time field should be set to '1'")
}

func TestCalculateTax_RequestMarshaling(t *testing.T) {
	t.Parallel()

	// Create mock HTTP client
	mockRT := testutil.NewMockRoundTripperSuccess(200, testutil.SampleAPIResponse200)
	client := &Client{
		httpClient: &http.Client{
			Transport: mockRT,
		},
	}

	// Create request with all fields
	req := &types.TaxRequest{
		Year:             "2024",
		TaxRegion:        "scotland",
		Age:              "35",
		Pension:          "5%",
		GrossWage:        60000,
		Plan:             "plan2",
		Extra:            1000,
		TaxCode:          "1257L",
		Married:          "y",
		Blind:            "y",
		ExNI:             "y",
		PartnerGrossWage: 25000,
	}

	// Execute
	_, err := client.CalculateTax(req)

	// Assert
	require.NoError(t, err)

	// Get request body and verify all fields
	body, err := mockRT.GetRequestBody()
	require.NoError(t, err)

	var sentReq types.TaxRequest
	err = json.Unmarshal([]byte(body), &sentReq)
	require.NoError(t, err)

	assert.Equal(t, "2024", sentReq.Year)
	assert.Equal(t, "scotland", sentReq.TaxRegion)
	assert.Equal(t, "35", sentReq.Age)
	assert.Equal(t, "5%", sentReq.Pension)
	assert.Equal(t, 60000, sentReq.GrossWage)
	assert.Equal(t, "plan2", sentReq.Plan)
	assert.Equal(t, 1000, sentReq.Extra)
	assert.Equal(t, "1257L", sentReq.TaxCode)
	assert.Equal(t, "y", sentReq.Married)
	assert.Equal(t, "y", sentReq.Blind)
	assert.Equal(t, "y", sentReq.ExNI)
	assert.Equal(t, 25000, sentReq.PartnerGrossWage)
}
