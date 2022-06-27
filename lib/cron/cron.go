package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func Send_Daily_message() {
	c := cron.New()

	c.AddFunc("*/1 * * * *",
		func() {
			fmt.Println("Hi")
		},
	)

	c.Start()
}
