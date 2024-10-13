package extractors

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/NitinN77/spends-tracker/global"
)

type HDFCCardTxn struct {
	CardName string
	Amount   float64
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
				return true, HDFCCardTxn{hdfcCard.Name, amount}
			}
		}
		return false, HDFCCardTxn{}
	} else {
		return false, HDFCCardTxn{}
	}
}
