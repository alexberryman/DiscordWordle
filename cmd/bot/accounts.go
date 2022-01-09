package main

import (
	"DiscordGoTurnips/internal/turnips/generated-code"
	"context"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

func updateAccountTimeZone(ctx context.Context, input string, CmdTimeZone string, s *discordgo.Session, m *discordgo.MessageCreate, q *turnips.Queries, a turnips.Account) {
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

	_, _ = q.UpdateTimeZone(ctx, turnips.UpdateTimeZoneParams{
		DiscordID: a.DiscordID,
		TimeZone:  timezoneInput,
	})

	flushEmojiAndResponseToDiscord(s, m, response)
}

func getOrCreateAccount(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, existingAccount int64, existingNickname int64, q *turnips.Queries) turnips.Account {
	var account turnips.Account
	var nickname turnips.Nickname
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
		nickname, _ = q.GetNickname(ctx, turnips.GetNicknameParams{
			DiscordID: m.Author.ID,
			ServerID:  m.GuildID,
		})
		if nickname.Nickname != name {
			var err error
			nickname, err = q.UpdateNickname(ctx, turnips.UpdateNicknameParams{
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
		nickname, _ = q.CreateNickname(ctx, turnips.CreateNicknameParams{
			DiscordID: m.Author.ID,
			ServerID:  m.GuildID,
			Nickname:  name,
		})

		reactToMessage(s, m, "🆕")
	}
	return account
}
