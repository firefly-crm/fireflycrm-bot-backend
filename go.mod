module github.com/firefly-crm/fireflycrm-bot-backend

go 1.14

replace github.com/firefly-crm/common => ../common

require (
	github.com/DarthRamone/modulbank-go v0.0.5
	github.com/badoux/checkmail v0.0.0-20181210160741-9661bd69e9ad
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/envoyproxy/protoc-gen-validate v0.3.0 // indirect
	github.com/firefly-crm/common v0.0.2
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/gojuno/goose v1.0.0
	github.com/gojuno/minimock/v3 v3.0.6
	github.com/golangci/golangci-lint v1.27.0
	github.com/hexdigest/gowrap v1.1.7
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334 // indirect
	github.com/jmoiron/sqlx v1.2.1-0.20190826204134-d7d95172beb5
	github.com/lyft/protoc-gen-star v0.4.15 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/twitchtv/twirp v5.10.2+incompatible // indirect
)
