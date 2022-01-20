package main

import (
	"DiscordWordle/internal/wordle/generated-code"
	"context"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

// Variables used for command line parameters
var (
	Token       string
	DatabaseUrl string
)

type response struct {
	Text  string
	Emoji string
}

var db *sql.DB

const cmdHistory = "history"
const cmdUpdate = "update"
const cmdScoreboard = "scoreboard"
const cmdPreviousWeek = "previous"
const cmdQuip = "quip"
const cmdQuipEnable = "enable"
const cmdQuipDisable = "disable"
const cmdTimeZone = "timezone"
const cmdWordle = "Wordle"
const noSolutionResult = "X"
const noSolutionGuesses = 7

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
		q := wordle.New(db)
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

		existingNickname, err := q.CountNicknameByDiscordIdAndServerId(ctx, wordle.CountNicknameByDiscordIdAndServerIdParams{
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

		var account wordle.Account
		if m.Message.GuildID != "" {
			account = getOrCreateAccount(ctx, s, m, existingAccount, existingNickname, q)
		} else {
			account = wordle.Account{
				DiscordID: m.Message.Author.ID,
				TimeZone:  "America/Chicago",
			}
		}

		routeMessageToAction(ctx, s, m, input, account, q, botMentionToken)
	}
}

func routeMessageToAction(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, input string, account wordle.Account, q *wordle.Queries, botMentionToken string) {
	var r response

	if strings.Contains(input, cmdWordle) {
		log.Println("Found a Wordle")
		gameId, guesses := extractGameGuesses(input)
		log.Println(fmt.Sprintf("%d - %d for %s", gameId, guesses, m.Author.Username))
		persistScore(ctx, m, s, account, gameId, guesses)

	} else if strings.HasPrefix(input, cmdUpdate) {
		gameId, guesses := extractGameGuesses(input)
		updateExistingScore(ctx, m, s, account, gameId, guesses)
	} else if strings.HasPrefix(input, cmdHistory) {
		getHistory(ctx, m, s, account)
	} else if strings.HasPrefix(input, cmdQuip+" "+cmdQuipEnable) {
		enableQuips(ctx, m, s)
	} else if strings.HasPrefix(input, cmdQuip+" "+cmdQuipDisable) {
		disableQuips(ctx, m, s)
	} else if strings.HasPrefix(input, cmdQuip) {
		score, quip := extractScoreQuip(input)
		persistQuip(ctx, m, s, account, score, quip)
	} else if strings.HasPrefix(input, cmdScoreboard+" "+cmdPreviousWeek) {
		getPreviousScoreboard(ctx, m, s)
	} else if strings.HasPrefix(input, cmdScoreboard) {
		getScoreboard(ctx, m, s)
	} else if strings.HasPrefix(input, cmdTimeZone) {
		updateAccountTimeZone(ctx, input, cmdTimeZone, s, m, q, account)
	} else if strings.HasPrefix(input, "help") {
		helpResponse(s, m, botMentionToken)
	} else {
		r.Text = "Wut?"
		flushEmojiAndResponseToDiscord(s, m, r)
	}
}

func extractScoreQuip(input string) (int, string) {
	var dataExp = regexp.MustCompile(`(?P<score>\d+)\s(?P<quip>.+)`)
	result := matchGroupsToStringMap(input, dataExp)
	score, _ := strconv.Atoi(result["score"])
	return score, result["quip"]
}

func extractGameGuesses(input string) (int, int) {
	var dataExp = regexp.MustCompile(fmt.Sprintf(`(?P<game_id>\d+)\s(?P<guesses>\d+|%s)`, noSolutionResult))
	result := matchGroupsToStringMap(input, dataExp)
	gameId, _ := strconv.Atoi(result["game_id"])
	var guesses int
	if strings.ToUpper(result["guesses"]) == noSolutionResult {
		guesses = noSolutionGuesses
	} else {
		guesses, _ = strconv.Atoi(result["guesses"])
	}
	return gameId, guesses
}

func matchGroupsToStringMap(input string, dataExp *regexp.Regexp) map[string]string {
	match := dataExp.FindStringSubmatch(input)
	if len(match) == 0 {
		log.Println(fmt.Sprintf("%s didn't match %s", input, dataExp))
	}
	result := make(map[string]string)
	for i, name := range dataExp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}

func flushEmojiAndResponseToDiscord(s *discordgo.Session, m *discordgo.MessageCreate, r response) {
	reactToMessage(s, m, r.Emoji)
	respondAsNewMessage(s, m, r.Text)
}

func respondAsNewMessage(s *discordgo.Session, m *discordgo.MessageCreate, response string) {
	if response != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Println("Error responding:", err)
		}
	}
}

func reactToMessage(s *discordgo.Session, m *discordgo.MessageCreate, reactionEmoji string) {
	err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, reactionEmoji)
	if err != nil {
		log.Println("Error adding and emoji:", err)
	}
}
