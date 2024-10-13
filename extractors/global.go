package extractors

import "strings"

func checkIfCreditCardEmail(snippet string) bool {
	if strings.Contains(snippet, "Credit Card") || strings.Contains(snippet, "credit card") {
		return true
	} else {
		return false
	}
}
