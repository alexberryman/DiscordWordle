package main

import (
	wordle "DiscordWordle/internal/wordle/generated-code"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func persistScore(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session, a wordle.Account, gameId int, guesses int) {
	response, scoreObj := buildScoreObjFromInput(a, gameId, guesses)

	scoreParams := wordle.CreateScoreParams{
		DiscordID: a.DiscordID,
		GameID:    scoreObj.GameID,
		Guesses:   scoreObj.Guesses,
	}

	q := wordle.New(db)
	_, err := q.CreateScore(ctx, scoreParams)

	if err != nil {
		response.Emoji = "⛔"
		response.Text = "You already created a price for this game, try updating it if it's wrong"
	} else {
		response = scoreColorfulResponse(guesses, m)
	}
	flushEmojiAndResponseToDiscord(s, m, response)
}

func getScores(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session, a wordle.Account) {

	historyByAccountParams := wordle.GetScoreHistoryByAccountParams{
		DiscordID: a.DiscordID,
		ServerID:  m.GuildID,
	}

	q := wordle.New(db)
	scores, err := q.GetScoreHistoryByAccount(ctx, historyByAccountParams)

	var response response

	if err != nil {
		response.Emoji = "⛔"
		response.Text = "Not finding any previous scores"
	} else {
		response.Emoji = "👍"
		response.Text = fmt.Sprintf("Found dem %d scores, boss!", len(scores))
		for _, v := range scores {
			response.Text += fmt.Sprintf("\n game: %d - %d/6", v.GameID, v.Guesses)
		}
	}
	flushEmojiAndResponseToDiscord(s, m, response)
}

func updateExistingScore(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session, a wordle.Account, gameId int, guesses int) {
	response, wordlecoreObj := buildScoreObjFromInput(a, gameId, guesses)

	priceParams := wordle.UpdateScoreParams{
		DiscordID: a.DiscordID,
		GameID:    wordlecoreObj.GameID,
		Guesses:   wordlecoreObj.Guesses,
	}

	q := wordle.New(db)
	_, err := q.UpdateScore(ctx, priceParams)

	if err != nil {
		response.Emoji = "⛔"
		response.Text = "I didn't find an existing price."
	} else {
		response = scoreColorfulResponse(guesses, m)
	}

	flushEmojiAndResponseToDiscord(s, m, response)
}

func buildScoreObjFromInput(a wordle.Account, gameId int, guesses int) (response, wordle.WordleScore) {
	var response response

	scoreThing := wordle.WordleScore{
		DiscordID: a.DiscordID,
		GameID:    int32(gameId),
		Guesses:   int32(guesses),
	}

	return response, scoreThing
}

func scoreColorfulResponse(guesses int, m *discordgo.MessageCreate) response {
	var response response
	response = selectResponseText(guesses, m, response)
	response = selectResponseEmoji(guesses, response)
	return response
}

func selectResponseText(guesses int, m *discordgo.MessageCreate, response response) response {
	if guesses >= 0 && guesses <= 6 {
		responseParams := wordle.GetResponseByScoreParams{
			ScoreValue:         guesses,
			InsideJokeServerID: m.GuildID,
		}

	} else if guesses == 69 {
		response.Text = "nice."
	} else {
		response.Text = "Is that even a real number? Did you fail to guess it?"
	}

	return response
}

func selectResponseEmoji(guesses int, response response) response {
	if guesses == 69 {
		response.Emoji = "♋️"
	} else if guesses == 0 {
		response.Emoji = "0️⃣"
	} else if guesses == 1 {
		response.Emoji = "1️⃣"
	} else if guesses == 2 {
		response.Emoji = "2️⃣"
	} else if guesses == 3 {
		response.Emoji = "3️⃣"
	} else if guesses == 4 {
		response.Emoji = "4️⃣"
	} else if guesses == 5 {
		response.Emoji = "5️⃣"
	} else if guesses == 6 {
		response.Emoji = "6️⃣"
	} else {
		response.Emoji = "❌"
	}

	return response
}
