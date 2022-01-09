package main

import (
	"DiscordGoTurnips/internal/turnips/generated-code"
	"context"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// Variables used for command line parameters
var (
	Token       string
	DatabaseUrl string
)

var db *sql.DB

//Weekday Named integer for weekdays
type Weekday int

const (
	sunday Weekday = iota
	monday
	tuesday
	wednesday
	thursday
	friday
	saturday
)

const cmdGraph = "graph"
const cmdTimeZone = "timezone"
const cmdUpdate = "update"

func init() {
	Token = os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		log.Println("DISCORD_TOKEN must be set")
	}

	DatabaseUrl = os.Getenv("DATABASE_URL")
	if DatabaseUrl == "" {
		log.Println("DATABASE_URL must be set")
	}

	dbConnection, err := sql.Open("postgres", DatabaseUrl)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	db = dbConnection
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	_ = dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var r response

	// Ignore all messages created by the bot itself
	botName := s.State.User.Username
	if m.Author.ID == s.State.User.ID {
		return
	}

	tokenizedContent, err := m.ContentWithMoreMentionsReplaced(s)
	if err != nil {
		log.Println("Failed to replace mentions:", err)
		return
	}

	botMentionToken := fmt.Sprintf("@%s", botName)
	if strings.HasPrefix(tokenizedContent, botMentionToken) {
		input := strings.TrimSpace(strings.Replace(tokenizedContent, botMentionToken, "", 1))
		q := turnips.New(db)
		ctx := context.Background()

		r.Emoji = "❌"
		existingAccount, err := q.CountAccountsByDiscordId(ctx, m.Author.ID)
		if err != nil {
			log.Println(err)
			r.Text = "Nice work! You broke the one thing that made people happy."
			r.Emoji = "🔥"
			flushEmojiAndResponseToDiscord(s, m, r)
			return
		}

		existingNickname, err := q.CountNicknameByDiscordId(ctx, turnips.CountNicknameByDiscordIdParams{
			DiscordID: m.Author.ID,
			ServerID:  m.GuildID,
		})
		if err != nil {
			log.Println(err)
			r.Text = "Nice work! You broke the one thing that made people happy."
			r.Emoji = "🔥"
			flushEmojiAndResponseToDiscord(s, m, r)
			return
		}

		account := getOrCreateAccount(ctx, s, m, existingAccount, existingNickname, q)

		routeMessageToAction(ctx, s, m, input, account, q, botMentionToken)
	}
}

func routeMessageToAction(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, input string, account turnips.Account, q *turnips.Queries, botMentionToken string) {
	var r response

	if turnipPrice, err := strconv.Atoi(input); err == nil {
		persistTurnipPrice(ctx, m, s, account, turnipPrice)
	} else if strings.Contains(input, cmdGraph) {
		historyInput := strings.TrimSpace(strings.Replace(input, cmdGraph, "", 1))
		if historyInput == "" {
			linkUsersCurrentPrices(s, m, AcTurnipsChartLink)
		} else if historyInput == "all" {
			linkServersCurrentPrices(s, m, AcTurnipsChartLink)
		} else if offset, err := strconv.Atoi(historyInput); err == nil {
			linkAccountsPreviousPrices(m, s, offset*(-1), AcTurnipsChartLink)
		} else if strings.HasPrefix(historyInput, "all") {
			historicalServerInput := strings.TrimSpace(strings.Replace(historyInput, "all", "", 1))
			if offset, err := strconv.Atoi(historicalServerInput); err == nil {
				linkServersPreviousPrices(m, s, offset*(-1), AcTurnipsChartLink)
			} else {
				r.Text = "That isn't a valid week offset. Use -1, -2, -3 etc..."
				r.Emoji = "⏰"
				flushEmojiAndResponseToDiscord(s, m, r)
			}
		} else {
			r.Emoji = "⛔"
			r.Text = "That is not a valid history request"
			flushEmojiAndResponseToDiscord(s, m, r)
		}

	} else if strings.Contains(input, cmdUpdate) {
		updateInput := strings.TrimSpace(strings.Replace(input, cmdUpdate, "", 1))
		if updateTurnipPrice, err := strconv.Atoi(updateInput); err == nil {
			updateExistingTurnipPrice(ctx, s, m, account, updateTurnipPrice)
		} else {
			r.Emoji = "⛔"
			r.Text = "That is not a valid price"
			flushEmojiAndResponseToDiscord(s, m, r)
		}

	} else if strings.HasPrefix(input, cmdTimeZone) {
		updateAccountTimeZone(ctx, input, cmdTimeZone, s, m, q, account)
	} else if strings.HasPrefix(input, "help") {
		helpResponse(s, m, botMentionToken, cmdGraph, cmdTimeZone)
	} else {
		r.Text = "Wut?"
		flushEmojiAndResponseToDiscord(s, m, r)
	}
}

func flushEmojiAndResponseToDiscord(s *discordgo.Session, m *discordgo.MessageCreate, r response) {
	reactToMessage(s, m, r.Emoji)
	respondAsNewMessage(s, m, r.Text)
}

func respondAsNewMessage(s *discordgo.Session, m *discordgo.MessageCreate, response string) {
	_, err := s.ChannelMessageSend(m.ChannelID, response)
	if err != nil {
		log.Println("Error responding:", err)
	}
}

func reactToMessage(s *discordgo.Session, m *discordgo.MessageCreate, reactionEmoji string) {
	err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, reactionEmoji)
	if err != nil {
		log.Println("Error adding and emoji:", err)
	}
}
