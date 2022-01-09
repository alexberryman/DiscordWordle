module DiscordGoTurnips

go 1.14

// Comment below is needed for heroku-buildpack-go https://github.com/heroku/heroku-buildpack-go/issues/301

// +heroku goVersion go1.14

require (
	github.com/bwmarrin/discordgo v0.20.3
	github.com/gobuffalo/packr/v2 v2.8.0 // indirect
	github.com/karrick/godirwalk v1.15.6 // indirect
	github.com/lib/pq v1.3.0
	github.com/rubenv/sql-migrate v0.0.0-20200402132117-435005d389bc
	github.com/sirupsen/logrus v1.5.0 // indirect
	golang.org/x/crypto v0.0.0-20200414173820-0848c9571904 // indirect
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a // indirect
	golang.org/x/sys v0.0.0-20200413165638-669c56c373c4 // indirect
)
