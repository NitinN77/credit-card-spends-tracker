package extractors

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/NitinN77/credit-card-spends-tracker/global"
)

type HDFCCardTxn struct {
	CardName string
	Last4    string
	Amount   float64
	Merchant string
}

func ExtractHDFCCard(snippet string, hdfcCardDetails []global.CardDetails) (bool, HDFCCardTxn) {
	if !checkIfCreditCardEmail(snippet) {
		return false, HDFCCardTxn{}
	}

	hdfcCCTxnRegex := regexp.MustCompile(`Rs\s(\d+\.\d{2})`)

	if !strings.Contains(snippet, "HDFC") {
		return false, HDFCCardTxn{}
	}

	match := hdfcCCTxnRegex.FindStringSubmatch(snippet)

	if match != nil {
		amount, _ := strconv.ParseFloat(match[1], 64)

		for _, hdfcCard := range hdfcCardDetails {
			if strings.Contains(snippet, "HDFC Bank Credit Card ending "+hdfcCard.Last4) {
				merchant := extractStringBetween(snippet, "at", "on")
				return true, HDFCCardTxn{hdfcCard.Name, hdfcCard.Last4, amount, merchant}
			}
		}
		return false, HDFCCardTxn{}
	} else {
		return false, HDFCCardTxn{}
	}
}
