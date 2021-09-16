package config

import (
	"log"
	"time"

	// "gopkg.in/ini.v1"
	"gopkg.in/go-ini/ini.v1"
)

type ConfigList struct {
	MaxSell int
	MaxBuy  int

	BTCBuyAmount float64

	ApiKey      string
	ApiSecret   string
	LogFile     string
	ProductCode string

	Mysql string

	Durations map[string]time.Duration

	LineSecret string
	LineToken  string
}

var BaseURL string

var Config ConfigList

func NewConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Println(err)
	}
	durations := map[string]time.Duration{
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
	}

	Config = ConfigList{
		ApiKey:       cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret:    cfg.Section("bitflyer").Key("api_secret").String(),
		MaxSell:      cfg.Section("bitflyer").Key("max_sell_order").MustInt(),
		MaxBuy:       cfg.Section("bitflyer").Key("max_buy_order").MustInt(),
		BTCBuyAmount: cfg.Section("bitflyer").Key("btc_buy_amount").MustFloat64(),
		LogFile:      cfg.Section("tradeSetting").Key("logfile_path").String(),
		ProductCode:  cfg.Section("tradeSetting").Key("product_code").String(),
		Durations:    durations,
		Mysql:        cfg.Section("database").Key("mysql").String(),
		LineSecret:   cfg.Section("line").Key("secret").String(),
		LineToken:    cfg.Section("line").Key("token").String(),
	}
	BaseURL = cfg.Section("bitflyer").Key("base_url").String()

}
