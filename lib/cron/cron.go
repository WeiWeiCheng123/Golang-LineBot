package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func Send_Daily_message() {
	c := cron.New()

	c.AddFunc("* 15 * * *",
		func() {
			fmt.Println("Hi")
			fmt.Println(time.Now())
		},
	)

	c.Start()
}
