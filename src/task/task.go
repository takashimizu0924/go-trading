package task

import (
	"app/bitflyer"
	"app/config"
	"app/enums"

	// "app/line"
	"app/models"
	"app/utils"
	"log"
	"runtime"
	"time"

	"github.com/carlescere/scheduler"
)

func StartBfService() {
	log.Println(" [StartBfService] start")
	config.NewConfig()

	apiClient := bitflyer.NewBitflyer(
		config.Config.ApiKey,
		config.Config.ApiSecret,
		config.Config.MaxBuy,
		config.Config.MaxSell,
	)
	ticker := func() {
		getTicker("BTC_JPY", apiClient)
	}
	buyingJob1 := func() {
		placeBuyOrder(enums.Stg1BtcLtp997, "BTC_JPY", config.Config.BTCBuyAmount, apiClient)
	}
	sellingJob1 := func() {
		sellPlaceOrder("BTC_JPY", apiClient)
	}
	scheduler.Every(1).Minutes().Run(ticker)
	scheduler.Every(1).Minutes().Run(buyingJob1)
	scheduler.Every(1).Minutes().Run(sellingJob1)
	runtime.Goexit()
}

func getTicker(product_code string, apiClient *bitflyer.APIClient) {
	get_ticker, err := apiClient.GetTicker(product_code)
	if err != nil {
		log.Printf("action=main/GetTicker err=%s", err.Error())
		return
	}
	log.Printf("BestBid:%10.2f BestAsk:%10.2f", get_ticker.BestBid, get_ticker.BestAsk)
}

func placeBuyOrder(strategy int, producCode string, size float64, apiClient *bitflyer.APIClient) {
	log.Printf("strategy:%v", strategy)
	log.Println(" [buying Job] start of job")

	shouldSkip := false

	buyPrice := 0.0

	if !shouldSkip {
		ticker, _ := apiClient.GetTicker(producCode)
		if strategy < 10 {
			buyPrice = utils.CalculateBuyPrice(ticker.Ltp, ticker.BestBid, strategy)
		} else {
			return
		}

		minuteToExpire := models.CalculateMinuteToExpire(strategy)

		order := &bitflyer.Order{
			ProductCode:     producCode,
			ChildOrderType:  "LIMIT",
			Side:            "BUY",
			Price:           buyPrice,
			Size:            size,
			MinuteToExpires: minuteToExpire,
			TimeInForce:     "GTC",
		}
		utc, _ := time.LoadLocation("UTC")
		utc_current_date := time.Now().In(utc)
		event := models.OrderEvent{
			OrderID:     "12345678910",
			Time:        utc_current_date,
			ProductCode: producCode,
			Side:        order.Side,
			Price:       buyPrice,
			Size:        size,
			Exchange:    "bitflyer",
			Strategy:    strategy,
		}
		// bot, err := line.NewLine()
		// if err != nil {
		// 	log.Println(err)
		// }
		// err = line.PostTextMessage("20210915", "150", "0.001", "2", bot)
		// if err != nil {
		// 	log.Println(err)
		// }

		err := event.BuyOrder()
		if err != nil {
			log.Println("BuyOrder failed.....")
		} else {
			log.Println("BuyOrder Success!!!!!")
		}

		log.Printf("order_price:%10.2f order_size:%g", order.Price, order.Size)
		balance, err := apiClient.GetBalance()
		if err != nil {
			log.Println(err)
		}
		log.Println("My Balance :", balance)

		// active_orders, err := apiClient.GetActiveBuyOrders(producCode,"COMPLETED")
		// if err != nil {
		// 	log.Println(err)
		// }
		// log.Println("active_order_Length:",len(*active_orders))
		// for i,val := range *active_orders{
		// 	log.Printf("No:%d completedPrice:%10.2f ",i,val.AveragePrice)
		// }
		// もし注文がACTIVE
		// 注文は見送り
		// もし注文がCOMPETED
		// PlaceOrderを実行させる
		// if len(*active_orders) == 0 {
		// 	// ここで指値注文
		// 	res, err := apiClient.PlaceOrder(order)
		// 	if err != nil || res == nil {
		// 		log.Println("購入失敗、、、")
		// 		shouldSkip = true
		// 	}
		// } else {
		// 	log.Println("still watting.....")
		// }

	}
}

func sellPlaceOrder(product_code string, apiClient *bitflyer.APIClient) {
	log.Println(" [SellOrder start] ")
	ticker, _ := apiClient.GetTicker(product_code)
	active_orders, err := apiClient.GetActiveBuyOrders(product_code, "COMPLETED")
	if err != nil {
		log.Println(err)
	}
	// log.Println("active_orders",*active_orders)
	for _, val := range *active_orders {
		res := utils.CalculateSellPrice(val.Size, ticker.BestBid, val.AveragePrice)
		log.Println("CalculeteSellResponse", res)
	}

}

func StreamIngectionData() {
	var tickerChanel = make(chan bitflyer.Ticker)
	apiClient := bitflyer.NewBitflyer(config.Config.ApiKey, config.Config.ApiSecret, config.Config.MaxBuy, config.Config.MaxSell)
	go apiClient.RealTimeGetTicker("BTC_JPY", tickerChanel)
	for ticker := range tickerChanel {
		log.Printf("action=StreamIngectionData, %v", ticker)
		for _, duration := range config.Config.Durations {
			isCreated := models.CreateCandleWithDuration(ticker, "BTC_JPY", duration)
			if isCreated == true {
				//TODO
			}
		}
	}
}
