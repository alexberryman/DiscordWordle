package main

import (
	"DiscordWordle/internal/wordle/generated-code"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
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
const hardModeIndicator = "*"
const noSolutionGuesses = 7

func init() {
	Token = os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		log.Fatal().Msg("DISCORD_TOKEN must be set")
	}

	DatabaseUrl = os.Getenv("DATABASE_URL")
	if DatabaseUrl == "" {
		log.Fatal().Msg("DATABASE_URL must be set")
	}

	dbConnection, err := sql.Open("postgres", DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot connect to database: %s", DatabaseUrl)
	}

	db = dbConnection
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating Discord session")
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("error opening connection to discord over websocket")
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Info().Msg("Bot is now running.  Press CTRL-C to exit.")
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
		log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Msg("Failed to replace mentions")
		return
	}

	botMentionToken := fmt.Sprintf("@%s", botName)
	wordleScoreDetected, err := mentionlessWordleScoreDetection(tokenizedContent)
	if strings.HasPrefix(tokenizedContent, botMentionToken) || wordleScoreDetected {
		input := strings.TrimSpace(strings.Replace(tokenizedContent, botMentionToken, "", 1))
		q := wordle.New(db)
		ctx := context.Background()

		r.Emoji = "❌"
		existingAccount, err := q.CountAccountsByDiscordId(ctx, m.Author.ID)
		if err != nil {
			log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID)
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
			log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID)
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
		gameId, guesses, err := extractGameGuesses(input)
		if err != nil {
			log.Error().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdWordle).Err(err).Msg("Error parsing guess count")
		}
		log.Info().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdWordle).Int("guesses", guesses).Int("game_id", gameId).Msg("Found a Wordle")
		persistScore(ctx, m, s, account, gameId, guesses)

	} else if strings.HasPrefix(input, cmdUpdate) {
		gameId, guesses, err := extractGameGuesses(input)
		if err != nil {
			log.Error().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdUpdate).Err(err).Msg("Error parsing guess count")
		}
		log.Info().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdUpdate).Int("guesses", guesses).Int("game_id", gameId).Msg("Updated a Wordle")
		updateExistingScore(ctx, m, s, account, gameId, guesses)
	} else if strings.HasPrefix(input, cmdHistory) {
		getHistory(ctx, m, s, account)
	} else if strings.HasPrefix(input, cmdQuip+" "+cmdQuipEnable) {
		enableQuips(ctx, m, s)
	} else if strings.HasPrefix(input, cmdQuip+" "+cmdQuipDisable) {
		disableQuips(ctx, m, s)
	} else if strings.HasPrefix(input, cmdQuip) {
		score, quip, err := extractScoreQuip(input)
		if err != nil {
			log.Error().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", cmdQuip).Err(err).Msg("Error parsing quip")
		}
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
		log.Info().Str("server_id", m.GuildID).Str("input", input).Str("author", m.Author.ID).Str("command", "").Msg("Failed to match command")
		r.Text = "Wut?"
		r.Emoji = "🤷"
		flushEmojiAndResponseToDiscord(s, m, r)
	}
}

func extractScoreQuip(input string) (int, string, error) {
	var dataExp = regexp.MustCompile(`(?P<score>\d+)\s(?P<quip>.+)`)

	result, err := matchGroupsToStringMap(input, dataExp)
	if err != nil {
		return 0, "", err
	}

	score, _ := strconv.Atoi(result["score"])
	return score, result["quip"], nil
}

func extractGameGuesses(input string) (int, int, error) {
	var dataExp = regexp.MustCompile(fmt.Sprintf(`(?P<game_id>\d+)\s(?P<guesses>\d+|%s)`, noSolutionResult))
	result, err := matchGroupsToStringMap(input, dataExp)
	if err != nil {
		return 0, 0, err
	}
	gameId, _ := strconv.Atoi(result["game_id"])
	var guesses int
	if strings.ToUpper(result["guesses"]) == noSolutionResult {
		guesses = noSolutionGuesses
	} else {
		guesses, _ = strconv.Atoi(result["guesses"])
	}
	return gameId, guesses, nil
}

func mentionlessWordleScoreDetection(input string) (bool, error) {
	var dataExp = regexp.MustCompile(fmt.Sprintf(`Wordle (?P<game_id>\d+)\s(?P<guesses>\d+|%s)/6[\%s]?\n`, noSolutionResult, hardModeIndicator))
	result, err := matchGroupsToStringMap(input, dataExp)
	if err != nil {
		return false, err
	}

	return len(result) > 0, nil
}

func matchGroupsToStringMap(input string, dataExp *regexp.Regexp) (map[string]string, error) {
	match := dataExp.FindStringSubmatch(input)
	result := make(map[string]string)
	if len(match) == 0 {
		errorMessage := fmt.Sprintf("%s didn't match %s", input, dataExp)
		return result, errors.New(errorMessage)

	}
	for i, name := range dataExp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result, nil
}

func flushEmojiAndResponseToDiscord(s *discordgo.Session, m *discordgo.MessageCreate, r response) {
	reactToMessage(s, m, r.Emoji)
	respondAsNewMessage(s, m, r.Text)
}

func respondAsNewMessage(s *discordgo.Session, m *discordgo.MessageCreate, response string) {
	if response != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Str("response", response).Msg("Error responding")
		}
	}
}

func reactToMessage(s *discordgo.Session, m *discordgo.MessageCreate, reactionEmoji string) {
	if reactionEmoji != "" {
		err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, reactionEmoji)
		if err != nil {
			log.Error().Err(err).Str("server_id", m.GuildID).Str("content", m.Content).Str("author", m.Author.ID).Str("reaction", reactionEmoji).Msg("Error reacting")
		}
	}
}
