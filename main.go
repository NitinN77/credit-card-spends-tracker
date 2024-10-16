package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/NitinN77/credit-card-spends-tracker/utils"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func main() {

	args := os.Args[1:]
	command := args[0]

	db, err := sqlx.Connect("sqlite", "file:spends-tracker.db")
	if err != nil {
		log.Fatalln(err)
	}
	utils.InitDB(db)

	if command == "--fetch" {
		startDate := args[1]
		endDate := args[2]

		parsedStartDate := utils.DateStringToTime(startDate)
		parsedEndDate := utils.DateStringToTime(endDate)

		fetchTransactions(parsedStartDate, parsedEndDate, db)
	} else if command == "--alias" {
		alias := args[1]
		merchant := args[2]

		utils.StoreAlias(db, alias, merchant)
		fmt.Printf("Stored alias %s for merchant %s\n", alias, merchant)
	} else if command == "--filter" {
		merchant := args[1]
		var transactions []utils.TransactionDB

		if len(args) == 4 {
			startDate := args[2]
			endDate := args[3]

			parsedStartDate := utils.DateStringToTime(startDate)
			parsedEndDate := utils.DateStringToTime(endDate)

			transactions = utils.GetTransactionsByMerchantWithinDateRange(db, merchant, parsedStartDate, parsedEndDate)
		} else {
			transactions = utils.GetTransactionsByMerchant(db, merchant)
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 5, 3, ' ', 0)
		for _, txn := range transactions {
			fmt.Fprintf(w, "%s\t%s\t%s\tRs.%.2f\t%s\n", txn.Date, txn.CardName, txn.Last4, txn.Amount, txn.Merchant)
		}
		w.Flush()

	}
}
