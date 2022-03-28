package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client
var token string
var secret string

func main() {
	token = "J6bbUTM7oYuS90LFq8wl3zJJEb46aPMEHI2W7pnvl9rFxK/lkts0wwe6CsHIw41wM5oG8z+SOtZq0B0aB01BuN2oOIl8HfPD28y/l/1nZg4s7jvLHJHPtXnDHjBdnyLly6Gjvv965Q+X2mou0r/PggdB04t89/1O/w1cDnyilFU="
	secret = "701d115eea8075888df5c048afe2f0ec"
	bot, err := linebot.New(secret, token)
	fmt.Println(bot, " ", err)
	http.HandleFunc("/callback", callbackHandler)
	port := 8080
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
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
