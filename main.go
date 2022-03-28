package main

import (
	"fmt"
	"log"
	"os"
	//"strconv"

	"github.com/WeiWeiCheng123/Golang-LineBot/lib/config"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func main() {
	bot, err := linebot.New(config.GetStr("TOKEN"), config.GetStr("SECRET"))
	log.Println("Bot:", bot, " err:", err)
	router := gin.Default()
	router.POST("/callback", callbackHandler)
	port := os.Getenv("PORT")
	fmt.Println("port= ", port)
	//addr := fmt.Sprintf(":%s", port)
	router.Run(":" + port)
}

func callbackHandler(c *gin.Context) {
	fmt.Println("Hello ", c.Request)
	events, err := bot.ParseRequest(c.Request)
	fmt.Println("env= ", events)
	fmt.Println("err= ", err)
	/*
	events, err := bot.ParseRequest(c.Request)
	fmt.Println("env= ", events)
	fmt.Println("err= ", err)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(400, "123")
		} else {
			c.String(500, "456")
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				// GetMessageQuota: Get how many remain free tier push message quota you still have this month. (maximum 500)
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				// message.ID: Msg unique ID
				// message.Text: Msg text
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("msg ID:"+message.ID+":"+"Get:"+message.Text+" , \n OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
					log.Print(err)
				}

			// Handle only on Sticker message
			case *linebot.StickerMessage:
				var kw string
				for _, k := range message.Keywords {
					kw = kw + "," + k
				}

				outStickerResult := fmt.Sprintf("收到貼圖訊息: %s, pkg: %s kw: %s  text: %s", message.StickerID, message.PackageID, kw, message.Text)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outStickerResult)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
	*/
}