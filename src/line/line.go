package line

import (
	"app/config"
	"log"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func NewLine() (*linebot.Client, error) {
	bot, err := linebot.New(config.Config.LineSecret, config.Config.LineToken)
	if err != nil {
		log.Printf("action=NewLine err=%s", err)
	}
	return bot, err
}

func PostTextMessage(date, profit, count, ppt string, bot *linebot.Client) error {
	var sb strings.Builder
	sb.WriteString("【bitflyer 自動売買 収益】\n")
	sb.WriteString("date / profit / count / ppt\n")
	sb.WriteString(date + "/" + profit + "/" + count + "/" + ppt)
	message := linebot.NewTextMessage(sb.String())
	if _, err := bot.BroadcastMessage(message).Do(); err != nil {
		log.Println("action=PostTextMessage err=:", err)
		return err
	}
	return nil
}
