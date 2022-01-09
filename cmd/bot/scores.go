package main

import (
	turnips "DiscordWordle/internal/turnips/generated-code"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func persistScore(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session, a turnips.Account, gameId int, guesses int) {
	response, scoreObj := buildScoreObjFromInput(a, gameId, guesses)

	scoreParams := turnips.CreateScoreParams{
		DiscordID: a.DiscordID,
		GameID:    scoreObj.GameID,
		Guesses:   scoreObj.Guesses,
	}

	q := turnips.New(db)
	_, err := q.CreateScore(ctx, scoreParams)

	if err != nil {
		response.Emoji = "⛔"
		response.Text = "You already created a price for this game, try updating it if it's wrong"
	} else {
		response = scoreColorfulResponse(guesses)
	}
	flushEmojiAndResponseToDiscord(s, m, response)
}

func updateExistingScore(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, a turnips.Account, gameId int, guesses int) {
	response, turnipScoreObj := buildScoreObjFromInput(a, gameId, guesses)

	priceParams := turnips.UpdateScoreParams{
		DiscordID: a.DiscordID,
		GameID:    turnipScoreObj.GameID,
		Guesses:   turnipScoreObj.Guesses,
	}

	q := turnips.New(db)
	_, err := q.UpdateScore(ctx, priceParams)

	if err != nil {
		response.Emoji = "⛔"
		response.Text = "I didn't find an existing price."
	} else {
		response = scoreColorfulResponse(guesses)
	}

	flushEmojiAndResponseToDiscord(s, m, response)
}

func buildScoreObjFromInput(a turnips.Account, gameId int, guesses int) (response, turnips.WordleScore) {
	var response response

	scoreThing := turnips.WordleScore{
		DiscordID: a.DiscordID,
		GameID:    int32(gameId),
		Guesses:   int32(guesses),
	}

	return response, scoreThing
}

func scoreColorfulResponse(guesses int) response {
	var response response
	response.Emoji = "✅"
	if guesses == 69 {
		response.Text = "nice."
	} else if guesses == 1 {
		response.Text = fmt.Sprintf("Oh, you got it in %d guess? Cheater. No one belives you.", guesses)
	} else if guesses == 2 {
		response.Text = fmt.Sprintf("%d? Dope.", guesses)
	} else if guesses == 3 {
		response.Text = fmt.Sprintf("Three is par for the course")
	} else if guesses == 4 {
		response.Text = fmt.Sprintf("%d isn't that bad but you can do better", guesses)
	} else if guesses == 5 {
		response.Text = fmt.Sprintf("Hey, at least you're not the guy that got 6 today.")
	} else if guesses == 6 {
		response.Text = fmt.Sprintf("You're really taking down to the wire using %d gusses", guesses)
	} else {
		response.Text = "Is that even a real number? Did you fail to guess it?"
		response.Emoji = "❌"
	}
	return response
}
