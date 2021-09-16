package models

import (
	"app/enums"
	"app/utils"
	"log"
	"time"
)

type OrderEvent struct {
	OrderID     string    `json:"order_id"`
	Time        time.Time `json:"time"`
	ProductCode string    `json:"product_code"`
	Side        string    `json:"side"`
	Price       float64   `json:"price"`
	Size        float64   `json:"size"`
	Exchange    string    `json:"exchange"`
	Filled      int       `json:"filled"`
	Strategy    int       `json:"strategy"`
}

func (orderEvent *OrderEvent) BuyOrder() error {
	cmd1, err := AppDB.Prepare("INSERT INTO buy_orders (order_id,time,product_code,side,price,size,exchange,strategy) VALUES (?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Println("[ERROR]BuyOrder01:", err)
		return err
	}
	_, err = cmd1.Exec(orderEvent.OrderID, orderEvent.Time, orderEvent.ProductCode, orderEvent.Side, orderEvent.Price, orderEvent.Size, orderEvent.Exchange, orderEvent.Strategy)
	if err != nil {
		log.Println("[ERROR]BuyOrder02:", err)
		return err
	}
	return nil
}

func (orderEvent *OrderEvent) SellOrder() error {
	cmd1, err := AppDB.Prepare("INSERT INTO sell_orders (order_id, time, product_code, side, price, size, exchange,strategy) VALUES (?,?,?,?,?,?,?")
	if err != nil {
		log.Println("[ERROR]SellOrder01:", err)
		return err
	}
	_, err = cmd1.Exec(orderEvent.OrderID, orderEvent.Time, orderEvent.ProductCode, orderEvent.Side, orderEvent.Price, orderEvent.Size, orderEvent.Exchange, orderEvent.Strategy)
	if err != nil {
		log.Println("[ERROR]SellOrder02:", err)
		return err
	}
	return nil
}

func FilledCheck(productCode string) ([]string, error) {
	cmd, _ := AppDB.Prepare(`SELECT order_id FROM buy_orders WEHERE filled = 0 and order_id != '' and product_code = ? union SELECT order_id FROM sell_orders WEHER filled = 0 and order_id != '' and product_code = ?`)
	rows, err := cmd.Query(productCode, productCode)
	if err != nil {
		log.Println("Failure to Exec quey......", err)
		return nil, err
	}
	var cnt = 0
	var ids []string
	for rows.Next() {
		var orderId string

		if err := rows.Scan(&orderId); err != nil {
			log.Println("Failure to get records....", err)
			return nil, err
		}
		cnt++
		ids = append(ids, orderId)
	}
	return ids, nil
}

type BuyOrderInfo struct {
	OrderID     string  `json:"order_id"`
	Price       float64 `json:"price"`
	ProductCode string  `json:"product_code"`
	Size        float64 `json:"size"`
	Exchange    string  `json:"exchange"`
	Strategy    float64 `json:"strategy"`
}

func (buyOrderInfo *BuyOrderInfo) CalculateSellOrderPrice() float64 {
	if buyOrderInfo.Strategy == enums.Stg3BtcLtp90 {
		return utils.Round(buyOrderInfo.Price * 1.03)
	} else {
		return utils.Round(buyOrderInfo.Price * 1.015)
	}
}

func CalculateMinuteToExpire(strategy int) int {
	if strategy == enums.Stg3BtcLtp90 {
		return 1440 //1day
	} else {
		return 3600 //2.5days
	}
}

func CheckFilledBuyOrders() []BuyOrderInfo {
	rows, err := AppDB.Query(`SELECT order_id, product_code, size, exchange, strategy FROM buy_orders WHERE filled = 1 and order_id != ''`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var cnt = 0
	var buyOrderInfos []BuyOrderInfo
	for rows.Next() {
		var order_id string
		var price float64
		var product_code string
		var size float64
		var exchange string
		var strategy int

		if err := rows.Scan(&order_id, &product_code, &size, &exchange, &strategy); err != nil {
			log.Println("Failure to get records......")
		}
		cnt++
		buyOrderInfo := BuyOrderInfo{OrderID: order_id, Price: price, ProductCode: product_code, Size: size, Exchange: exchange}
		buyOrderInfos = append(buyOrderInfos, buyOrderInfo)
	}
	return buyOrderInfos
}

func UpdateFilledOrder(order_id string) error {
	cmd1, _ := AppDB.Prepare(`update buy_orders set filled = 1 where order_id = ?`)
	_, err := cmd1.Exec(order_id)
	if err != nil {
		return err
	}
	cmd2, _ := AppDB.Prepare(`update sell_orders set filled = 1 where order_id = ?`)
	_, err = cmd2.Exec(order_id)
	if err != nil {
		return err
	}
	return nil
}
