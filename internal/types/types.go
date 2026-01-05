// Package types defines data structures used across the application.
package types

// TaxRequest represents the API request structure
type TaxRequest struct {
	Response         string `json:"response"`
	Year             string `json:"year"`
	TaxRegion        string `json:"taxregion"`
	Age              string `json:"age"`
	Pension          string `json:"pension"`
	Time             string `json:"time"`
	GrossWage        int    `json:"grosswage"`
	Plan             string `json:"plan,omitempty"`
	Extra            int    `json:"extra,omitempty"`
	TaxCode          string `json:"taxcode,omitempty"`
	Married          string `json:"married,omitempty"`
	Blind            string `json:"blind,omitempty"`
	ExNI             string `json:"exNI,omitempty"`
	PartnerGrossWage int    `json:"partnerGrossWage,omitempty"`
}

// TaxBracket represents tax at a specific rate
type TaxBracket struct {
	Rate   float64 `json:"rate"`
	Amount float64 `json:"amount"`
}

// TaxResponse represents the API response structure
type TaxResponse struct {
	TaxYear                  int                   `json:"tax_year"`
	TaxablePay               float64               `json:"taxable_pay"`
	GrossPay                 float64               `json:"gross_pay"`
	AdditionalGross          float64               `json:"additional_gross"`
	TaxFreeAllowance         float64               `json:"tax_free_allowance"`
	TaxPaid                  float64               `json:"tax_paid"`
	TaxDue                   map[string]TaxBracket `json:"tax_due"`
	NationalInsurance        float64               `json:"national_insurance"`
	NetPay                   float64               `json:"net_pay"`
	StudentLoanRepayment     float64               `json:"student_loan_repayment"`
	PensionHMRC              float64               `json:"pension_hmrc"`
	PensionYou               float64               `json:"pension_you"`
	PensionClaimback         float64               `json:"pension_claimback"`
	EmployersNI              float64               `json:"employers_ni"`
	TaxFreeMarried           float64               `json:"tax_free_married"`
	TaxRegion                string                `json:"tax_region"`
	TaxCode                  string                `json:"tax_code"`
	TaxFreeMarriageAllowance float64               `json:"tax_free_marriage_allowance"`
	GrossSacrifice           float64               `json:"gross_sacrifice"`
	ChildcarePre2011         interface{}           `json:"childcare_pre2011"`
	Debug                    interface{}           `json:"debug"`
	ChildcareAmount          float64               `json:"childcare_amount"`
	Previous                 *TaxResponse          `json:"previous,omitempty"`
}

// ComparisonResult represents one option's calculation result with its label
type ComparisonResult struct {
	Label    string
	Request  *TaxRequest
	Response *TaxResponse
}
