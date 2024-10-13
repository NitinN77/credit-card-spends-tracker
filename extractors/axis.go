package extractors

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/NitinN77/credit-card-spends-tracker/global"
)

type AxisCardTxn struct {
	CardName string
	Amount   float64
}

func ExtractAxisCard(snippet string, axisCardDetails []global.CardDetails) (bool, AxisCardTxn) {

	if !checkIfCreditCardEmail(snippet) {
		return false, AxisCardTxn{}
	}

	axisCCTxnRegex := regexp.MustCompile(`INR\s(\d+(?:\.\d+)?)(\sat)?`)

	if !strings.Contains(snippet, "total credit limit") {
		return false, AxisCardTxn{}
	}

	updatedSnippet := strings.ReplaceAll(snippet, "is INR", "")

	match := axisCCTxnRegex.FindStringSubmatch(updatedSnippet)

	if match != nil {
		amount, _ := strconv.ParseFloat(match[1], 64)

		for _, axisCard := range axisCardDetails {
			if strings.Contains(updatedSnippet, "XX"+axisCard.Last4) {
				return true, AxisCardTxn{axisCard.Name, amount}
			}
		}
		return false, AxisCardTxn{}
	} else {
		return false, AxisCardTxn{}
	}

}
