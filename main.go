package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
		fmt.Print(c.Request)
		x, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println(x)
		c.String(400, "Hello")
	})

	router.POST("/callback", callbackHandler)

	router.Run(":" + port)
}

func callbackHandler(c *gin.Context) {
	fmt.Println(c.Request)
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(400, "123")
			log.Print(err)
		}
		return
	}
	for _, event := range events {
		fmt.Println("event= ", event)
	}
}
