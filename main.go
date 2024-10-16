package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

		parsedStartDate, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}

		parsedEndDate, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}

		fetchTransactions(parsedStartDate, parsedEndDate, db)
	} else if command == "--alias" {
		alias := args[1]
		merchant := args[2]

		utils.StoreAlias(db, alias, merchant)
		fmt.Printf("Stored alias %s for merchant %s\n", alias, merchant)
	}
}
