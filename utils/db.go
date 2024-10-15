package utils

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type TransactionDB struct {
	ID       int     `db:"id"`
	Amount   float64 `db:"amount"`
	CardName string  `db:"card_name"`
	Date     string  `db:"date"`
}

func InitDB(db *sqlx.DB) {
	schema := `
		CREATE TABLE IF NOT EXISTS cc_transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount REAL,
			card_name TEXT,
			date TEXT
		);
		CREATE TABLE IF NOT EXISTS fetched_date_ranges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			start_date TEXT,
			end_date TEXT
		);
	`
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalln(err)
	}
}

func SaveTransactionToDB(db *sqlx.DB, cardName string, amount float64, timestamp int64) {
	query := `
		INSERT INTO cc_transactions (card_name, amount, date)
		VALUES (?, ?, ?);
	`

	t := time.Unix(timestamp/1000, 0)
	formattedDate := t.Format("2006-01-02")

	_, err := db.Exec(query, cardName, amount, formattedDate)
	if err != nil {
		log.Fatalln("Error inserting transaction into DB:", err)
	}
}

func SaveDateRangeToDB(db *sqlx.DB, startDate, endDate *time.Time) {
	query := `
		INSERT INTO fetched_date_ranges (start_date, end_date)
		VALUES (?, ?);
	`

	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	_, err := db.Exec(query, startDateStr, endDateStr)
	if err != nil {
		log.Fatalln("Error inserting date range:", err)
	}
}

func GetTransactions(db *sqlx.DB, startDate, endDate time.Time) []TransactionDB {
	var transactions []TransactionDB

	query := `
		SELECT id, amount, card_name, date
		FROM cc_transactions
		WHERE date >= ? AND date <= ?
	`

	err := db.Select(&transactions, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		log.Fatalln("Error retrieving transactions: ", err)
	}

	return transactions
}
