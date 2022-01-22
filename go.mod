module DiscordWordle

go 1.14

// Comment below is needed for heroku-buildpack-go https://github.com/heroku/heroku-buildpack-go/issues/301

// +heroku goVersion go1.14

require (
	github.com/bwmarrin/discordgo v0.23.2
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/lib/pq v1.10.4
	github.com/rs/zerolog v1.26.1
	github.com/rubenv/sql-migrate v1.0.0
	github.com/sirupsen/logrus v1.5.0 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
)
