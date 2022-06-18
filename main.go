package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func main() {
	var p Person
	fmt.Println(p)
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
		fmt.Println(eve)
		fmt.Println(string(eve))
		fmt.Println(eve[0])
		c.String(200, "test parse req pass")
	})

	router.POST("/callback", callbackHandler)

	router.Run(":" + port)
}

func callbackHandler(c *gin.Context) {
	events, err := bot.ParseRequest(c.Request)
	fmt.Println(events)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(400, err.Error())
			log.Print(err)
		}
		return
	}
	c.String(200, "success")
}

/*
[{
	"replyToken":"3e19e6a2de9e4d52912d387e85bfadaf",
	"type":"message",
	"mode":"active",
	"timestamp":1655557543382,
	"source":{"type":"user","userId":"Ua4712856c697d2d1e02d02c33343f3ea"},
	"message":{"id":"16283810978053","type":"sticker","packageId":"5788726","stickerId":"123222087","stickerResourceType":"STATIC","keywords":["Straight face"]},
	"webhookEventId":"01G5VEPNS365JCNW2XYNQ036ZJ","deliveryContext":{"isRedelivery":false}
}]
*/
