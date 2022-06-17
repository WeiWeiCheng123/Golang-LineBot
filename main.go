package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	log.Println("Bot:", bot, " err:", err)
	router := gin.Default()
	router.POST("/callback", callbackHandler)

	router.Run(":" + port)
}

func callbackHandler(c *gin.Context) {
	events, err := ParseRequest1(c.Request)
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
		fmt.Println("event= ", event)
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
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

func ParseRequest1(r *http.Request) ([]*linebot.Event, error) {
	return ParseRequest(os.Getenv("SECRET"), r)
}

func ParseRequest(channelSecret string, r *http.Request) ([]*linebot.Event, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("1")
		return nil, err
	}
	if !validateSignature(channelSecret, r.Header.Get("X-Line-Signature"), body) {
		fmt.Println("2")
		return nil, errors.New("ErrInvalidSignature")
	}

	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err = json.Unmarshal(body, request); err != nil {
		fmt.Println("3")
		return nil, err
	}
	fmt.Println("done")
	return request.Events, nil
}

func validateSignature(channelSecret, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))

	_, err = hash.Write(body)
	if err != nil {
		return false
	}

	return hmac.Equal(decoded, hash.Sum(nil))
}
