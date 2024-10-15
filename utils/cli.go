package utils

import (
	"os"

	"github.com/NitinN77/credit-card-spends-tracker/global"
	"github.com/pelletier/go-toml/v2"
)

type AppConfig struct {
	UserEmail       string               `toml:"user_email"`
	SourceEmails    []string             `toml:"source_emails"`
	AxisCardDetails []global.CardDetails `toml:"axis_card_details"`
	HDFCCardDetails []global.CardDetails `toml:"hdfc_card_details"`
}

func GetAppConfig() AppConfig {
	var appConfig AppConfig

	doc, err := os.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal([]byte(doc), &appConfig)
	if err != nil {
		panic(err)
	}
	return appConfig
}
