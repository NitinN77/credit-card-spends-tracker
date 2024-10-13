package main

import (
	"fmt"
	"log"

	"time"

	"github.com/NitinN77/spends-tracker/extractors"
	"github.com/NitinN77/spends-tracker/gmail"
	"github.com/NitinN77/spends-tracker/utils"
)

func main() {

	srv := gmail.GetGmailService()
	appConfig := utils.GetAppConfig()

	today := time.Now()
	dateNDaysAgo := utils.GetDateNDaysAgo(today, appConfig.FetchDaysCount)
	sourceEmailsStr := utils.GenerateFromEmailsQuery(appConfig.SourceEmails)

	emailList, err := srv.Users.Messages.List(appConfig.UserEmail).Q(fmt.Sprintf("after:%s {%s}", dateNDaysAgo, sourceEmailsStr)).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve email list: %v", err)
	}

	cardTotals := make(map[string]float64)

	for _, email := range emailList.Messages {
		emailData, err := srv.Users.Messages.Get(appConfig.UserEmail, email.Id).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve email: %v", err)
		}

		snippet := emailData.Snippet

		isHDFCCard, hdfcTxn := extractors.ExtractHDFCCard(snippet, appConfig.HDFCCardDetails)

		if isHDFCCard {
			if cardTotal, exists := cardTotals[hdfcTxn.CardName]; exists {
				cardTotals[hdfcTxn.CardName] = cardTotal + hdfcTxn.Amount
			} else {
				cardTotals[hdfcTxn.CardName] = hdfcTxn.Amount
			}
			continue
		}

		isAxisCard, axisTxn := extractors.ExtractAxisCard(snippet, appConfig.AxisCardDetails)

		if isAxisCard {
			if cardTotal, exists := cardTotals[axisTxn.CardName]; exists {
				cardTotals[axisTxn.CardName] = cardTotal + axisTxn.Amount
			} else {
				cardTotals[axisTxn.CardName] = axisTxn.Amount
			}
			continue
		}

	}

	for cardName, totalSpent := range cardTotals {
		fmt.Printf("Spends with %s: %.2f\n", cardName, totalSpent)
	}
}
