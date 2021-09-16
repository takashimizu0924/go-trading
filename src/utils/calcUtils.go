package utils

import (
	"app/enums"
	"log"
	"math"
)

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func CalculateBuyPrice(ltp, low float64, strategy int) float64 {
	log.Printf("LTP:%10.2f  BestBid:%10.2f",ltp,low)
	if strategy == enums.Stg0BtcLtplow7{
		return Round(ltp*0.3 + low*0.7)
	} else if strategy == enums.Stg1BtcLtp997{
		return Round(ltp * 0.997)
	} else if strategy == enums.Stg2BtcLtp98{
		return Round(ltp * 0.98)
	} else {
		return Round(ltp * 0.95)
	}
}

func CalculateSellPrice(size,nowPrice,completedPrice float64) string{
	//持っているBTC(手数料込)より現在価格の方が高かったら売り
	var shouldSell = "売り"
	var waitSell = "待て"
	var sellSize = size
	var sellPrice = completedPrice * sellSize
	var price = nowPrice *sellSize
	log.Printf("Sell Price:%10.2f  NowPrice:%10.2f",sellPrice,price)
	if sellPrice <= price {
		log.Println("価格上昇しています")
		return shouldSell
	} else if sellPrice >= price {
		log.Println("価格下降しています")
		return waitSell
	} else {
		log.Println("価格変動なしです")
		return waitSell
	}
}

// 売る数量　＊　0.00015 = 手数料込みの売り数量　＊　現在価格