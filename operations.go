package main

import (
	"fmt"
	"log"

	"time"

	"github.com/NitinN77/credit-card-spends-tracker/extractors"
	"github.com/NitinN77/credit-card-spends-tracker/gmail"
	"github.com/NitinN77/credit-card-spends-tracker/utils"
	"github.com/jmoiron/sqlx"
	_gmail "google.golang.org/api/gmail/v1"
)

func fetchTransactions(startDate, endDate time.Time, db *sqlx.DB) {

	srv := gmail.GetGmailService()
	appConfig := utils.GetAppConfig()

	sourceEmailsStr := utils.GenerateFromEmailsQuery(appConfig.SourceEmails)

	var emailList []_gmail.Message

	fetchedRanges := utils.GetOverlappingRanges(startDate, endDate, db)
	missingRanges := utils.CalculateMissingRanges(startDate, endDate, fetchedRanges)

	for _, missingRange := range missingRanges {
		fetchedEmailList, err := srv.Users.Messages.List(appConfig.UserEmail).
			Q(fmt.Sprintf(
				"after:%s before:%s {%s}",
				missingRange.StartDate.Format("2006-01-02"),
				missingRange.EndDate.Format("2006-01-02"),
				sourceEmailsStr,
			)).
			Do()

		if err != nil {
			log.Fatalf("Unable to fetch email list: %v", err)
		}
		for _, email := range fetchedEmailList.Messages {
			emailList = append(emailList, *email)
		}
	}

	for _, email := range emailList {
		emailData, err := srv.Users.Messages.Get(appConfig.UserEmail, email.Id).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve email: %v", err)
		}

		snippet := emailData.Snippet

		isHDFCCard, hdfcTxn := extractors.ExtractHDFCCard(snippet, appConfig.HDFCCardDetails)

		if isHDFCCard {
			utils.SaveTransactionToDB(db, hdfcTxn.CardName, hdfcTxn.Amount, emailData.InternalDate)
			continue
		}

		isAxisCard, axisTxn := extractors.ExtractAxisCard(snippet, appConfig.AxisCardDetails)

		if isAxisCard {
			utils.SaveTransactionToDB(db, axisTxn.CardName, axisTxn.Amount, emailData.InternalDate)
			continue
		}
	}

	for _, missingRange := range missingRanges {
		utils.SaveDateRangeToDB(db, &missingRange.StartDate, &missingRange.EndDate)
	}

	fetchedTransactions := utils.GetTransactions(db, startDate, endDate)

	cardTotals := make(map[string]float64)

	for _, txn := range fetchedTransactions {
		if cardTotal, exists := cardTotals[txn.CardName]; exists {
			cardTotals[txn.CardName] = cardTotal + txn.Amount
		} else {
			cardTotals[txn.CardName] = txn.Amount
		}
	}

	for cardName, totalSpent := range cardTotals {
		fmt.Printf("Spends with %s: %.2f\n", cardName, totalSpent)
	}
}
