package config

import "time"

type (
	Config struct {
		Db              string
		Rabbit          string
		TgToken         string
		WatcherInterval time.Duration
	}
)
