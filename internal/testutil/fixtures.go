package testutil

// ValidConfigYAML is a complete valid config file for testing
const ValidConfigYAML = `defaults:
  region: scotland
  year: "2024"
  age: "30"
  pension: "5%"
  student-loan: plan2
  tax-code: 1257L
  extra: 1000
  period: monthly
  income: 50000
  married: true
  blind: false
  no-ni: false
  partner-income: 25000
`

// PartialConfigYAML is a config with only some fields set
const PartialConfigYAML = `defaults:
  region: uk
  year: "2025"
  period: weekly
`

// InvalidConfigYAML is malformed YAML for error testing
const InvalidConfigYAML = `defaults:
  region: uk
  year: "2025"
    extra indentation: invalid
  pension 5%
`

// SampleAPIResponse200 is a valid JSON response from the listentotaxman API
const SampleAPIResponse200 = `{
  "tax_year": 2024,
  "taxable_pay": 37430.0,
  "gross_pay": 50000.0,
  "additional_gross": 0.0,
  "tax_free_allowance": 12570.0,
  "tax_paid": 7486.0,
  "tax_due": {
    "0": {
      "rate": 0.20,
      "amount": 7486.0
    }
  },
  "national_insurance": 4218.16,
  "net_pay": 38295.84,
  "student_loan_repayment": 0.0,
  "pension_hmrc": 0.0,
  "pension_you": 0.0,
  "pension_claimback": 0.0,
  "employers_ni": 5220.78,
  "tax_free_married": 0.0,
  "tax_region": "uk",
  "tax_code": "1257L",
  "tax_free_marriage_allowance": 0.0,
  "gross_sacrifice": 0.0,
  "childcare_pre2011": null,
  "debug": null,
  "childcare_amount": 0.0
}`

// SampleAPIResponse500 is an error response
const SampleAPIResponse500 = `{
  "error": "Internal server error",
  "message": "Failed to calculate tax"
}`

// SampleAPIResponseInvalidJSON is malformed JSON for parsing error tests
const SampleAPIResponseInvalidJSON = `{
  "tax_year": 2024,
  "gross_pay": 50000.0,
  "invalid json here...
}`

// SampleAPIResponseWithStudentLoan includes student loan repayment
const SampleAPIResponseWithStudentLoan = `{
  "tax_year": 2024,
  "taxable_pay": 37430.0,
  "gross_pay": 50000.0,
  "additional_gross": 0.0,
  "tax_free_allowance": 12570.0,
  "tax_paid": 7486.0,
  "tax_due": {
    "0": {
      "rate": 0.20,
      "amount": 7486.0
    }
  },
  "national_insurance": 4218.16,
  "net_pay": 37845.84,
  "student_loan_repayment": 450.0,
  "pension_hmrc": 0.0,
  "pension_you": 0.0,
  "pension_claimback": 0.0,
  "employers_ni": 5220.78,
  "tax_free_married": 0.0,
  "tax_region": "uk",
  "tax_code": "1257L",
  "tax_free_marriage_allowance": 0.0,
  "gross_sacrifice": 0.0,
  "childcare_pre2011": null,
  "debug": null,
  "childcare_amount": 0.0
}`

// SampleAPIResponseWithAllStatus includes married, blind allowance, and NI exempt flags
const SampleAPIResponseWithAllStatus = `{
  "tax_year": 2024,
  "taxable_pay": 35430.0,
  "gross_pay": 50000.0,
  "additional_gross": 0.0,
  "tax_free_allowance": 14570.0,
  "tax_paid": 6486.0,
  "tax_due": {
    "0": {
      "rate": 0.20,
      "amount": 6486.0
    }
  },
  "national_insurance": 0.0,
  "net_pay": 43514.0,
  "student_loan_repayment": 0.0,
  "pension_hmrc": 0.0,
  "pension_you": 0.0,
  "pension_claimback": 0.0,
  "employers_ni": 0.0,
  "tax_free_married": 1260.0,
  "tax_region": "uk",
  "tax_code": "1457L",
  "tax_free_marriage_allowance": 1260.0,
  "gross_sacrifice": 0.0,
  "childcare_pre2011": null,
  "debug": null,
  "childcare_amount": 0.0
}`
