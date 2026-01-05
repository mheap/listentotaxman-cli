package display

import (
	"fmt"
	"strings"

	"github.com/mheap/listentotaxman-cli/internal/types"
)

// formatCurrency formats a float as currency with the £ symbol and thousand separators
func formatCurrency(amount float64) string {
	// Format to 2 decimal places
	str := fmt.Sprintf("%.2f", amount)

	// Split into integer and decimal parts
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := parts[1]

	// Add thousand separators to integer part
	intPart = addThousandSeparators(intPart)

	return fmt.Sprintf("£%s.%s", intPart, decPart)
}

// addThousandSeparators adds commas as thousand separators to a number string
func addThousandSeparators(s string) string {
	// Handle negative numbers
	negative := false
	if strings.HasPrefix(s, "-") {
		negative = true
		s = s[1:]
	}

	// Add commas from right to left
	n := len(s)
	if n <= 3 {
		if negative {
			return "-" + s
		}
		return s
	}

	var result strings.Builder
	for i, digit := range s {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}

	if negative {
		return "-" + result.String()
	}
	return result.String()
}

// getPeriodDivisor returns the divisor for a given period
func getPeriodDivisor(period string) float64 {
	switch period {
	case "yearly":
		return 1.0
	case "monthly":
		return 12.0
	case "weekly":
		return 52.0
	case "daily":
		return 365.0
	case "hourly":
		return 2080.0
	default:
		return 1.0
	}
}

// getPeriodLabel returns human-readable label for period
func getPeriodLabel(period string) string {
	switch period {
	case "yearly":
		return "Yearly"
	case "monthly":
		return "Monthly"
	case "weekly":
		return "Weekly"
	case "daily":
		return "Daily"
	case "hourly":
		return "Hourly"
	default:
		return "Yearly"
	}
}

// Summary displays the tax calculation as a summary table (Option A)
func Summary(resp *types.TaxResponse, period string, req *types.TaxRequest) {
	divisor := getPeriodDivisor(period)
	periodLabel := getPeriodLabel(period)

	// Build status line if any status flags are active
	statusParts := []string{}
	if req.Married == "y" {
		statusParts = append(statusParts, "Married")
	}
	if req.Blind == "y" {
		statusParts = append(statusParts, "Blind Allowance")
	}
	if req.ExNI == "y" {
		statusParts = append(statusParts, "NI Exempt")
	}
	statusLine := ""
	if len(statusParts) > 0 {
		statusLine = " • " + strings.Join(statusParts, " • ")
	}

	// Calculate padding for header alignment (total width = 45 chars inside borders)
	headerText := fmt.Sprintf("Tax Calculation for %d (%s) - %s", resp.TaxYear, resp.TaxRegion, periodLabel)
	padding := 45 - len(headerText)

	fmt.Printf("\n╔══════════════════════════════════════════════╗\n")
	fmt.Printf("║ %s%*s║\n", headerText, padding, "")
	if statusLine != "" {
		statusPadding := 45 - len(statusLine)
		fmt.Printf("║ %s%*s║\n", statusLine, statusPadding, "")
	}
	fmt.Printf("╠══════════════════════════════════════════════╣\n")

	// Main income and deductions - use right-aligned currency with proper width
	fmt.Printf("║ %-25s %18s ║\n", "Gross Salary", formatCurrency(resp.GrossPay/divisor))
	fmt.Printf("║ %-25s %18s ║\n", "Taxable Pay", formatCurrency(resp.TaxablePay/divisor))
	fmt.Printf("║ %-25s %18s ║\n", "Tax Paid", formatCurrency(resp.TaxPaid/divisor))
	fmt.Printf("║ %-25s %18s ║\n", "National Insurance", formatCurrency(resp.NationalInsurance/divisor))

	if resp.StudentLoanRepayment > 0 {
		fmt.Printf("║ %-25s %18s ║\n", "Student Loan", formatCurrency(resp.StudentLoanRepayment/divisor))
	}

	fmt.Printf("║ %-25s %18s ║\n", "Pension (You)", formatCurrency(resp.PensionYou/divisor))
	fmt.Printf("║ %-25s %18s ║\n", "Net Pay", formatCurrency(resp.NetPay/divisor))

	fmt.Printf("╠══════════════════════════════════════════════╣\n")

	// Employer costs
	totalCost := resp.GrossPay + resp.EmployersNI + resp.PensionHMRC
	fmt.Printf("║ %-25s %18s ║\n", "Employer's NI", formatCurrency(resp.EmployersNI/divisor))
	fmt.Printf("║ %-25s %18s ║\n", "Pension (HMRC)", formatCurrency(resp.PensionHMRC/divisor))
	fmt.Printf("║ %-25s %18s ║\n", "Total Cost", formatCurrency(totalCost/divisor))

	fmt.Printf("╚══════════════════════════════════════════════╝\n\n")
}

// Detailed displays the tax calculation with detailed breakdown (Option B)
func Detailed(resp *types.TaxResponse, period string, req *types.TaxRequest) {
	divisor := getPeriodDivisor(period)
	periodLabel := getPeriodLabel(period)

	fmt.Printf("Tax Year: %d (%s) - %s\n", resp.TaxYear, resp.TaxRegion, periodLabel)
	if resp.TaxCode != "" {
		fmt.Printf("Tax Code: %s\n", resp.TaxCode)
	}

	// Show status flags if any are active
	statusParts := []string{}
	if req.Married == "y" {
		statusParts = append(statusParts, "Married")
	}
	if req.Blind == "y" {
		statusParts = append(statusParts, "Blind Allowance")
	}
	if req.ExNI == "y" {
		statusParts = append(statusParts, "NI Exempt")
	}
	if len(statusParts) > 0 {
		fmt.Printf("Status: %s\n", strings.Join(statusParts, " • "))
	}

	fmt.Println()

	// Income section
	fmt.Println("Income:")
	fmt.Printf("  Gross Salary:        %15s\n", formatCurrency(resp.GrossPay/divisor))
	if resp.AdditionalGross > 0 {
		fmt.Printf("  Additional Gross:    %15s\n", formatCurrency(resp.AdditionalGross/divisor))
	}
	fmt.Printf("  Tax Free Allowance:  %15s\n", formatCurrency(resp.TaxFreeAllowance/divisor))
	fmt.Printf("  Taxable Pay:         %15s\n", formatCurrency(resp.TaxablePay/divisor))
	fmt.Println()

	// Tax breakdown
	fmt.Println("Tax Breakdown:")

	// Sort tax brackets by rate (0 = basic, 1 = higher, 2 = additional)
	if bracket, ok := resp.TaxDue["0"]; ok && bracket.Amount > 0 {
		ratePercent := bracket.Rate * 100
		fmt.Printf("  Basic Rate (%.0f%%):    %15s\n", ratePercent, formatCurrency(bracket.Amount/divisor))
	}
	if bracket, ok := resp.TaxDue["1"]; ok && bracket.Amount > 0 {
		ratePercent := bracket.Rate * 100
		fmt.Printf("  Higher Rate (%.0f%%):   %15s\n", ratePercent, formatCurrency(bracket.Amount/divisor))
	}
	if bracket, ok := resp.TaxDue["2"]; ok && bracket.Amount > 0 {
		ratePercent := bracket.Rate * 100
		fmt.Printf("  Additional (%.0f%%):    %15s\n", ratePercent, formatCurrency(bracket.Amount/divisor))
	}
	fmt.Printf("  Total Tax:           %15s\n", formatCurrency(resp.TaxPaid/divisor))
	fmt.Println()

	// Deductions
	fmt.Println("Deductions:")
	fmt.Printf("  National Insurance:  %15s\n", formatCurrency(resp.NationalInsurance/divisor))
	if resp.StudentLoanRepayment > 0 {
		fmt.Printf("  Student Loan:        %15s\n", formatCurrency(resp.StudentLoanRepayment/divisor))
	}
	fmt.Printf("  Pension (You):       %15s\n", formatCurrency(resp.PensionYou/divisor))

	totalDeductions := resp.TaxPaid + resp.NationalInsurance + resp.StudentLoanRepayment + resp.PensionYou
	fmt.Printf("  Total Deductions:    %15s\n", formatCurrency(totalDeductions/divisor))
	fmt.Println()

	// Net pay
	fmt.Printf("Net Pay:               %15s\n", formatCurrency(resp.NetPay/divisor))
	fmt.Println()

	// Employer costs
	fmt.Println("Employer Costs:")
	fmt.Printf("  Employer's NI:       %15s\n", formatCurrency(resp.EmployersNI/divisor))
	fmt.Printf("  Pension (HMRC):      %15s\n", formatCurrency(resp.PensionHMRC/divisor))

	totalCost := resp.GrossPay + resp.EmployersNI + resp.PensionHMRC
	fmt.Printf("  Total Cost:          %15s\n", formatCurrency(totalCost/divisor))
}
