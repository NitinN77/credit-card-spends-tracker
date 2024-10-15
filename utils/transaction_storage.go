package utils

import (
	"database/sql/driver"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type CustomDate struct {
	time.Time
}

func (ct *CustomDate) Scan(value interface{}) error {

	if v, ok := value.(string); ok {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		ct.Time = t
		return nil
	}
	return fmt.Errorf("failed to scan CustomTime: expected string, got %T", value)
}

func (ct CustomDate) Value() (driver.Value, error) {
	return ct.Time.Format("2006-01-02"), nil
}

type FetchedRangeDB struct {
	StartDate CustomDate `db:"start_date"`
	EndDate   CustomDate `db:"end_date"`
}

type FetchedRange struct {
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
}

func GetOverlappingRanges(startDate, endDate time.Time, db *sqlx.DB) []FetchedRange {
	var rangesDB []FetchedRangeDB
	var ranges []FetchedRange
	query := `
        SELECT start_date, end_date
        FROM fetched_date_ranges
        WHERE (start_date <= ? AND end_date >= ?) 
        OR (start_date >= ? AND start_date <= ?) 
        OR (end_date >= ? AND end_date <= ?)
		ORDER BY start_date;
	`
	err := db.Select(&rangesDB, query, endDate, startDate, startDate, endDate, startDate, endDate)
	if err != nil {
		log.Fatalf("Unable to fetch date ranges from DB: %v", err)
	}
	for _, fetchedRange := range rangesDB {
		ranges = append(ranges, FetchedRange{fetchedRange.StartDate.Time, fetchedRange.EndDate.Time})
	}
	return ranges
}

func CalculateMissingRanges(startDate, endDate time.Time, fetchedRanges []FetchedRange) []FetchedRange {
	missingRanges := []FetchedRange{}
	currentStart := startDate

	for _, fetched := range fetchedRanges {
		if currentStart.Before(fetched.StartDate) {
			missingRanges = append(missingRanges, FetchedRange{StartDate: currentStart, EndDate: fetched.StartDate.AddDate(0, 0, -1)})
		}
		if currentStart.Before(fetched.EndDate) {
			currentStart = fetched.EndDate.AddDate(0, 0, 1)
		}
	}

	if currentStart.Before(endDate) {
		missingRanges = append(missingRanges, FetchedRange{StartDate: currentStart, EndDate: endDate})
	}

	return missingRanges
}
