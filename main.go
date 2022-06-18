package main

import (
	"encoding/json"
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
	port := os.Getenv("PORT")
	fmt.Println("secret", secret)
	fmt.Println("token", token)
	fmt.Println("port ", port)

	bot, err := linebot.New(secret, token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bot:", bot, " err:", err)
	router := gin.Default()

	router.POST("/", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			log.Println(err)
		}
		// fmt.Println(events, err)
		eve, _ := json.Marshal(events)
		fmt.Println(string(eve))
		c.String(200, "test parse req pass")
	})

	router.POST("/callback", callbackHandler)

	router.Run(":" + port)
}

func callbackHandler(c *gin.Context) {
	fmt.Println(c.Request)
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(400, err.Error())
			log.Print(err)
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
	c.String(200, "success")
}
