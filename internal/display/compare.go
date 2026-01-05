package display

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mheap/listentotaxman-cli/internal/types"
)

// DisplayComparison displays a comparison table for multiple tax calculations
func DisplayComparison(results []types.ComparisonResult, period string, verbose bool) {
	divisor := getPeriodDivisor(period)

	// Calculate table dimensions
	numOptions := len(results)
	fieldColWidth := 20
	valueColWidth := 12

	// Truncate labels if needed (max 11 chars to fit in 12-char column)
	labels := make([]string, numOptions)
	for i, result := range results {
		if len(result.Label) > 11 {
			labels[i] = result.Label[:10] + "…"
		} else {
			labels[i] = result.Label
		}
	}

	// Generate borders
	topBorder := generateBorder(numOptions, fieldColWidth, valueColWidth, "top")
	midBorder := generateBorder(numOptions, fieldColWidth, valueColWidth, "middle")
	sepBorder := generateBorder(numOptions, fieldColWidth, valueColWidth, "separator")
	bottomBorder := generateBorder(numOptions, fieldColWidth, valueColWidth, "bottom")

	// Print header
	fmt.Println()
	fmt.Println(topBorder)

	// Print label row
	fmt.Print("║ ")
	fmt.Printf("%-*s", fieldColWidth, "Field")
	for _, label := range labels {
		fmt.Print(" ║ ")
		fmt.Printf("%-*s", valueColWidth, label)
	}
	fmt.Println(" ║")

	// Print status row if any option has status flags
	hasStatus := false
	statusLines := make([]string, numOptions)
	for i, result := range results {
		statusParts := []string{}
		if result.Request.Married == "y" {
			statusParts = append(statusParts, "M")
		}
		if result.Request.Blind == "y" {
			statusParts = append(statusParts, "B")
		}
		if result.Request.ExNI == "y" {
			statusParts = append(statusParts, "NI")
		}
		if len(statusParts) > 0 {
			hasStatus = true
			statusLines[i] = strings.Join(statusParts, "•")
		}
	}

	if hasStatus {
		fmt.Print("║ ")
		fmt.Printf("%-*s", fieldColWidth, "Status")
		for _, status := range statusLines {
			fmt.Print(" ║ ")
			fmt.Printf("%-*s", valueColWidth, status)
		}
		fmt.Println(" ║")
	}

	fmt.Println(midBorder)

	// Print fields
	if verbose {
		printComparisonFieldVerbose(results, divisor, fieldColWidth, valueColWidth)
	} else {
		printComparisonFieldSummary(results, divisor, fieldColWidth, valueColWidth)
	}

	// Print separator before employer costs
	fmt.Println(sepBorder)

	// Employer costs
	printComparisonRow("Employer's NI", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.EmployersNI })
	printComparisonRow("Pension (HMRC)", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.PensionHMRC })

	// Total cost
	printComparisonRow("Total Cost", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.GrossPay + r.EmployersNI + r.PensionHMRC })

	fmt.Println(bottomBorder)
	fmt.Println()
}

// printComparisonFieldSummary prints summary fields (non-verbose mode)
func printComparisonFieldSummary(results []types.ComparisonResult, divisor float64, fieldColWidth, valueColWidth int) {
	printComparisonRow("Gross Salary", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.GrossPay })
	printComparisonRow("Tax Paid", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.TaxPaid })
	printComparisonRow("National Insurance", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.NationalInsurance })
	printComparisonRow("Student Loan", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.StudentLoanRepayment })
	printComparisonRow("Pension (You)", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.PensionYou })
	printComparisonRow("Net Pay", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.NetPay })
}

// printComparisonFieldVerbose prints all fields (verbose mode)
func printComparisonFieldVerbose(results []types.ComparisonResult, divisor float64, fieldColWidth, valueColWidth int) {
	// Income section
	printComparisonRow("Gross Salary", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.GrossPay })
	printComparisonRow("Additional Gross", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.AdditionalGross })
	printComparisonRow("Tax Free Allowance", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.TaxFreeAllowance })
	printComparisonRow("Taxable Pay", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.TaxablePay })

	// Tax breakdown
	printComparisonRow("Basic Rate Tax", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 {
			if bracket, ok := r.TaxDue["0"]; ok {
				return bracket.Amount
			}
			return 0
		})
	printComparisonRow("Higher Rate Tax", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 {
			if bracket, ok := r.TaxDue["1"]; ok {
				return bracket.Amount
			}
			return 0
		})
	printComparisonRow("Additional Rate Tax", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 {
			if bracket, ok := r.TaxDue["2"]; ok {
				return bracket.Amount
			}
			return 0
		})
	printComparisonRow("Total Tax", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.TaxPaid })

	// Deductions
	printComparisonRow("National Insurance", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.NationalInsurance })
	printComparisonRow("Student Loan", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.StudentLoanRepayment })
	printComparisonRow("Pension (You)", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.PensionYou })
	printComparisonRow("Pension Claimback", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.PensionClaimback })

	// Net pay
	printComparisonRow("Net Pay", results, divisor, fieldColWidth, valueColWidth,
		func(r *types.TaxResponse) float64 { return r.NetPay })
}

// printComparisonRow prints a single comparison row
func printComparisonRow(fieldName string, results []types.ComparisonResult, divisor float64, fieldColWidth, valueColWidth int, extractor func(*types.TaxResponse) float64) {
	fmt.Print("║ ")
	fmt.Printf("%-*s", fieldColWidth, fieldName)

	for _, result := range results {
		value := extractor(result.Response) / divisor
		fmt.Print(" ║ ")
		fmt.Printf("%*s", valueColWidth, formatCurrency(value))
	}

	fmt.Println(" ║")
}

// generateBorder generates a table border based on the number of options
func generateBorder(numOptions, fieldColWidth, valueColWidth int, borderType string) string {
	var left, mid, right, horiz, vert string

	switch borderType {
	case "top":
		left, mid, right, horiz, vert = "╔", "╦", "╗", "═", "═"
	case "middle":
		left, mid, right, horiz, vert = "╠", "╬", "╣", "═", "═"
	case "separator":
		left, mid, right, horiz, vert = "╠", "╬", "╣", "═", "═"
	case "bottom":
		left, mid, right, horiz, vert = "╚", "╩", "╝", "═", "═"
	}

	// Build border: left + field column + (mid + value column) * numOptions + right
	var sb strings.Builder
	sb.WriteString(left)
	sb.WriteString(strings.Repeat(horiz, fieldColWidth+2)) // +2 for padding

	for i := 0; i < numOptions; i++ {
		sb.WriteString(mid)
		sb.WriteString(strings.Repeat(vert, valueColWidth+2)) // +2 for padding
	}

	sb.WriteString(right)
	return sb.String()
}

// DisplayComparisonJSON displays comparison results as a JSON comparison object
func DisplayComparisonJSON(results []types.ComparisonResult, period string) {
	divisor := getPeriodDivisor(period)

	// Build comparison object structure
	output := map[string]interface{}{
		"period":     period,
		"comparison": buildComparisonFields(results, divisor),
		"metadata":   buildMetadata(results),
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))
}

// buildComparisonFields builds the comparison object with all fields
func buildComparisonFields(results []types.ComparisonResult, divisor float64) map[string]map[string]float64 {
	fields := map[string]map[string]float64{
		"gross_pay":                   extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.GrossPay }),
		"taxable_pay":                 extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.TaxablePay }),
		"additional_gross":            extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.AdditionalGross }),
		"tax_free_allowance":          extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.TaxFreeAllowance }),
		"tax_paid":                    extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.TaxPaid }),
		"national_insurance":          extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.NationalInsurance }),
		"student_loan":                extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.StudentLoanRepayment }),
		"pension_you":                 extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.PensionYou }),
		"pension_hmrc":                extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.PensionHMRC }),
		"pension_claimback":           extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.PensionClaimback }),
		"net_pay":                     extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.NetPay }),
		"employers_ni":                extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.EmployersNI }),
		"total_cost":                  extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.GrossPay + r.EmployersNI + r.PensionHMRC }),
		"gross_sacrifice":             extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.GrossSacrifice }),
		"childcare_amount":            extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.ChildcareAmount }),
		"tax_free_married":            extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.TaxFreeMarried }),
		"tax_free_marriage_allowance": extractField(results, divisor, func(r *types.TaxResponse) float64 { return r.TaxFreeMarriageAllowance }),
	}

	// Add tax brackets
	fields["basic_rate_tax"] = extractField(results, divisor, func(r *types.TaxResponse) float64 {
		if bracket, ok := r.TaxDue["0"]; ok {
			return bracket.Amount
		}
		return 0
	})
	fields["higher_rate_tax"] = extractField(results, divisor, func(r *types.TaxResponse) float64 {
		if bracket, ok := r.TaxDue["1"]; ok {
			return bracket.Amount
		}
		return 0
	})
	fields["additional_rate_tax"] = extractField(results, divisor, func(r *types.TaxResponse) float64 {
		if bracket, ok := r.TaxDue["2"]; ok {
			return bracket.Amount
		}
		return 0
	})

	return fields
}

// buildMetadata builds metadata section with tax year, region, code per option
func buildMetadata(results []types.ComparisonResult) map[string]map[string]interface{} {
	metadata := make(map[string]map[string]interface{})

	for _, result := range results {
		metadata[result.Label] = map[string]interface{}{
			"tax_year":   result.Response.TaxYear,
			"tax_region": result.Response.TaxRegion,
			"tax_code":   result.Response.TaxCode,
		}
	}

	return metadata
}

// extractField extracts a specific field from all results and applies divisor
func extractField(results []types.ComparisonResult, divisor float64, extractor func(*types.TaxResponse) float64) map[string]float64 {
	field := make(map[string]float64)

	for _, result := range results {
		field[result.Label] = extractor(result.Response) / divisor
	}

	return field
}
