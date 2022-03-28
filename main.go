package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func main() {
	secret := os.Getenv("SECRET")
	token := os.Getenv("TOKEN")
	fmt.Println(secret)
	fmt.Println(token)
	bot, err := linebot.New(secret, token)
	log.Println("Bot:", bot, " err:", err)
	router := gin.Default()
	router.POST("/callback", callbackHandler)
	port := os.Getenv("PORT")
	fmt.Println("port= ", port)
	//addr := fmt.Sprintf(":%s", port)
	router.Run(":" + port)
}

func callbackHandler(c *gin.Context) {
	fmt.Println("HELLO")
	fmt.Println(bot.ParseRequest(c.Request))
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
}
