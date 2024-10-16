package utils

import (
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type TransactionDB struct {
	ID       int     `db:"id"`
	Amount   float64 `db:"amount"`
	CardName string  `db:"card_name"`
	Last4    string  `db:"last_4"`
	Merchant string  `db:"merchant"`
	Date     string  `db:"date"`
}

type AliasDB struct {
	ID       int    `db:"id"`
	Alias    string `db:"alias"`
	Merchant string `db:"merchant"`
}

func InitDB(db *sqlx.DB) {
	schema, err := os.ReadFile("sql/init_tables.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalln(err)
	}

	var count int
	query := `SELECT COUNT(1) FROM merchant_aliases LIMIT 1;`

	err = db.Get(&count, query)
	if err != nil {
		log.Fatalf("Error fetching merchant table count from DB: %v", err)
	}

	if count == 0 {
		schema, err := os.ReadFile("sql/init_aliases.sql")
		if err != nil {
			log.Fatalf("Error reading SQL file: %v", err)
		}

		_, err = db.Exec(string(schema))
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func SaveTransactionToDB(db *sqlx.DB, cardName string, last4 string, amount float64, merchant string, timestamp int64) {
	query := `
		INSERT INTO cc_transactions (card_name, last_4, amount, merchant, date)
		VALUES (?, ?, ?, ?, ?);
	`

	t := time.Unix(timestamp/1000, 0)
	formattedDate := t.Format("2006-01-02")

	_, err := db.Exec(query, cardName, last4, amount, merchant, formattedDate)
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
		SELECT id, amount, card_name, last_4, merchant, date
		FROM cc_transactions
		WHERE date >= ? AND date <= ?
	`

	err := db.Select(&transactions, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		log.Fatalln("Error retrieving transactions: ", err)
	}

	return transactions
}

func GetMerchantAliases(db *sqlx.DB) map[string]string {
	aliasMap := make(map[string]string)
	var aliasesDB []AliasDB

	query := `
		SELECT id, alias, merchant
		FROM merchant_aliases;
	`
	err := db.Select(&aliasesDB, query)
	if err != nil {
		log.Fatalln("Error retrieving merchant aliases from DB: ", err)
	}

	for _, aliasDB := range aliasesDB {
		aliasMap[aliasDB.Alias] = aliasDB.Merchant
	}

	return aliasMap
}

func StoreAlias(db *sqlx.DB, alias, merchant string) {
	query := `
		INSERT INTO merchant_aliases (alias, merchant)
		VALUES (?, ?);
	`
	_, err := db.Exec(query, alias, merchant)
	if err != nil {
		log.Fatalln("Error inserting merchant alias:", err)
	}
}

func GetTransactionsByMerchant(db *sqlx.DB, merchant string) []TransactionDB {
	var aliasesDB []AliasDB

	query := `
		SELECT id, alias, merchant
		FROM merchant_aliases
		WHERE merchant = ?;
	`
	err := db.Select(&aliasesDB, query, merchant)
	if err != nil {
		log.Fatalln("Error retrieving aliases: ", err)
	}

	var aliasList []string

	for _, aliasDB := range aliasesDB {
		aliasList = append(aliasList, aliasDB.Alias)
	}

	var transactions []TransactionDB

	inQuery := `
		SELECT id, amount, card_name, last_4, merchant, date
		FROM cc_transactions
		WHERE merchant IN (?)
		ORDER BY date DESC;
	`
	query, args, err := sqlx.In(inQuery, aliasList)
	if err != nil {
		log.Fatalln("Error constructing IN query: ", err)
	}

	query = db.Rebind(query)
	err = db.Select(&transactions, query, args...)
	if err != nil {
		log.Fatalln("Error retrieving transactions: ", err)
	}

	return transactions
}

func GetTransactionsByMerchantWithinDateRange(db *sqlx.DB, merchant string, startDate, endDate time.Time) []TransactionDB {
	var aliasesDB []AliasDB

	query := `
		SELECT id, alias, merchant
		FROM merchant_aliases
		WHERE merchant = ?;
	`
	err := db.Select(&aliasesDB, query, merchant)
	if err != nil {
		log.Fatalln("Error retrieving aliases: ", err)
	}

	var aliasList []string

	for _, aliasDB := range aliasesDB {
		aliasList = append(aliasList, aliasDB.Alias)
	}

	var transactions []TransactionDB

	inQuery := `
		SELECT id, amount, card_name, last_4, merchant, date
		FROM cc_transactions
		WHERE merchant IN (?) AND date >= ? AND date <= ?
		ORDER BY date DESC;
	`
	query, args, err := sqlx.In(inQuery, aliasList, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		log.Fatalln("Error constructing IN query: ", err)
	}

	query = db.Rebind(query)
	err = db.Select(&transactions, query, args...)
	if err != nil {
		log.Fatalln("Error retrieving transactions: ", err)
	}

	return transactions
}
