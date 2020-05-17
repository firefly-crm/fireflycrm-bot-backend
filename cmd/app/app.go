package main

import "github.com/firefly-crm/fireflycrm-bot-backend/service"

func main() {
	bot := service.MustNewBot("token")
	_ = service.Service{
		Bot:       bot,
		OrderBook: nil,
		Users:     nil,
		Storage:   nil,
	}
}
