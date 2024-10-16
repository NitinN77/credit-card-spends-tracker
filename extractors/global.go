package extractors

import "strings"

func checkIfCreditCardEmail(snippet string) bool {
	if strings.Contains(snippet, "Credit Card") || strings.Contains(snippet, "credit card") {
		return true
	} else {
		return false
	}
}

func extractStringBetween(str, start, end string) string {
	startIndex := strings.Index(str, start)
	if startIndex == -1 {
		return ""
	}

	startIndex += len(start)

	endIndex := strings.Index(str[startIndex:], end)
	if endIndex == -1 {
		return ""
	}

	result := str[startIndex : startIndex+endIndex]
	return strings.TrimSpace(result)
}
