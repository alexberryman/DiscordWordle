package main

import (
	"DiscordWordle/internal/wordle/generated-code"
	"context"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

func updateAccountTimeZone(ctx context.Context, input string, CmdTimeZone string, s *discordgo.Session, m *discordgo.MessageCreate, q *wordle.Queries, a wordle.Account) {
	var response response

	timezoneInput := strings.TrimSpace(strings.Replace(input, CmdTimeZone, "", 1))
	_, err := time.LoadLocation(timezoneInput)

	if err != nil {
		response.Emoji = "⛔"
		response.Text = "Set a valid timezone from the `TZ database name` column https://en.wikipedia.org/wiki/List_of_tz_database_time_zones"
		flushEmojiAndResponseToDiscord(s, m, response)
	} else {
		response.Emoji = "✅"
	}

	_, _ = q.UpdateTimeZone(ctx, wordle.UpdateTimeZoneParams{
		DiscordID: a.DiscordID,
		TimeZone:  timezoneInput,
	})

	flushEmojiAndResponseToDiscord(s, m, response)
}

func getOrCreateAccount(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, existingAccount int64, existingNickname int64, q *wordle.Queries) wordle.Account {
	var account wordle.Account
	var nickname wordle.Nickname
	if existingAccount > 0 {
		account, _ = q.GetAccount(ctx, m.Author.ID)
		reactToMessage(s, m, "👤")
	} else {
		account, _ = q.CreateAccount(ctx, m.Author.ID)
		reactToMessage(s, m, "🆕")
	}

	var name string
	if m.Member.Nick != "" {
		name = m.Member.Nick
	} else {
		name = m.Author.Username
	}

	if existingNickname > 0 {
		nickname, _ = q.GetNickname(ctx, wordle.GetNicknameParams{
			DiscordID: m.Author.ID,
			ServerID:  m.GuildID,
		})
		if nickname.Nickname != name {
			var err error
			nickname, err = q.UpdateNickname(ctx, wordle.UpdateNicknameParams{
				DiscordID: m.Author.ID,
				Nickname:  name,
				ServerID:  m.GuildID,
			})
			if err != nil {
				log.Println("Failed to update nickname")
			} else {
				reactToMessage(s, m, "🔁")
			}
		}

	} else {
		nickname, _ = q.CreateNickname(ctx, wordle.CreateNicknameParams{
			DiscordID: m.Author.ID,
			ServerID:  m.GuildID,
			Nickname:  name,
		})

		reactToMessage(s, m, "🆕")
	}
	return account
}
