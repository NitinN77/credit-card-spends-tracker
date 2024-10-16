package utils

import (
	"log"
	"time"
)

func DateStringToTime(dateString string) time.Time {
	parsedStartDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		log.Fatalln("Error parsing date string: ", err)
	}
	return parsedStartDate
}
