package service

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
)

func MustNewBot(token string) *tg.BotAPI {
	var c *http.Client
	transport, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		transport.DisableKeepAlives = true
		var rt http.RoundTripper = transport
		c = &http.Client{Transport: rt}
	} else {
		c = http.DefaultClient
	}

	bot, err := tg.NewBotAPIWithClient(token, c)
	if err != nil {
		panic(err)
	}

	return bot
}
