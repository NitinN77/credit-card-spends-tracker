package utils

import (
	"fmt"
	"strings"
	"time"
)

func GetDateNDaysAgo(timeNow time.Time, n int) string {
	dateNDaysAgo := timeNow.AddDate(0, 0, -n)
	return dateNDaysAgo.Format("2006/01/02")
}

func GenerateFromEmailsQuery(sourceEmails []string) string {
	var sourceEmailsQuery []string
	for _, source := range sourceEmails {
		sourceEmailsQuery = append(sourceEmailsQuery, fmt.Sprintf("from:%s", source))
	}
	sourceEmailsStr := strings.Join(sourceEmailsQuery, " ")
	return sourceEmailsStr
}
