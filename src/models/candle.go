package models

import (
	"app/bitflyer"
	"fmt"
	"log"
	"time"
)

type Candle struct {
	ProductCode string
	Duration    time.Duration
	Time        time.Time
	Open        float64
	Close       float64
	High        float64
	Low         float64
	Volume      float64
}

func NewCandle(productCode string, duation time.Duration, timeDate time.Time, open, close, high, low, volume float64) *Candle {
	return &Candle{
		productCode,
		duation,
		timeDate,
		open,
		close,
		high,
		low,
		volume,
	}
}

func (candle *Candle) TableName() string {
	return GetCandleTableName(candle.ProductCode, candle.Duration)
}

func (candle *Candle) Create() error {
	cmd := fmt.Sprintf("INSERT INTO %s (time, open, close, hight, low, volume) VALUES (?, ?, ?, ?, ?, ?)", candle.TableName())
	_, err := AppDB.Exec(cmd, candle.Time.Format(time.RFC3339), candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
	log.Println(candle.Time.Format(time.RFC3339))
	if err != nil {
		log.Println("DB EXEC ERROR:", err)
		return err
	}
	return err
}

func (candle *Candle) Save() error {
	cmd := fmt.Sprintf("UPDATE %s SET open = ?, close = ?, hight = ?, low = ?, volume = ? WHERE time = ?", candle.TableName())
	_, err := AppDB.Exec(cmd, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume, candle.Time.Format(time.RFC3339))
	if err != nil {
		return err
	}
	return err
}

func GetCandle(productCode string, duration time.Duration, datetime time.Time) *Candle {
	tableName := GetCandleTableName(productCode, duration)
	cmd := fmt.Sprintf("SELECT time, open, close, hight, low, volume FROM %s WHERE time = ?", tableName)
	row := AppDB.QueryRow(cmd, datetime.Format(time.RFC3339))
	var candle Candle
	err := row.Scan(&candle.Time, &candle.Open, &candle.Close, &candle.High, &candle.Low, &candle.Volume)
	if err != nil {
		return nil
	}
	return NewCandle(productCode, duration, candle.Time, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
}

func CreateCandleWithDuration(ticker bitflyer.Ticker, productCode string, duration time.Duration) bool {
	currentCandle := GetCandle(productCode, duration, ticker.TruncateDateTime(duration))
	price := ticker.GetMiddlePrice()
	if currentCandle == nil {
		candle := NewCandle(productCode, duration, ticker.TruncateDateTime(duration), price, price, price, price, ticker.Volume)
		candle.Create()
		return true
	}

	if currentCandle.High <= price {
		currentCandle.High = price
	} else if currentCandle.Low >= price {
		currentCandle.Low = price
	}
	currentCandle.Volume += ticker.Volume
	currentCandle.Close = price
	currentCandle.Save()
	return false
}
