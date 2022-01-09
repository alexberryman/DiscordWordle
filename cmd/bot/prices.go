package main

import (
	"DiscordGoTurnips/internal/turnips/generated-code"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func persistTurnipPrice(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session, a turnips.Account, turnipPrice int) {
	response, turnipPriceObj, err := buildPriceObjFromInput(a, turnipPrice)

	if err != nil {
		log.Print(err)
		flushEmojiAndResponseToDiscord(s, m, response)
	}

	priceParams := turnips.CreatePriceParams{
		DiscordID: a.DiscordID,
		Price:     turnipPriceObj.Price,
		AmPm:      turnipPriceObj.AmPm,
		DayOfWeek: turnipPriceObj.DayOfWeek,
		DayOfYear: turnipPriceObj.DayOfYear,
		Year:      turnipPriceObj.Year,
		Week:      turnipPriceObj.Week,
	}

	q := turnips.New(db)
	_, err = q.CreatePrice(ctx, priceParams)

	if err != nil {
		response.Emoji = "⛔"
		response.Text = "You already created a price for this period"
	} else {
		response = turnipPriceColorfulResponse(turnipPrice)
	}
	flushEmojiAndResponseToDiscord(s, m, response)
	linkUsersCurrentPrices(s, m, AcTurnipsImageLink)
}

func updateExistingTurnipPrice(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, a turnips.Account, turnipPrice int) {
	response, turnipPriceObj, err := buildPriceObjFromInput(a, turnipPrice)
	if err != nil {
		flushEmojiAndResponseToDiscord(s, m, response)
	}

	priceParams := turnips.UpdatePriceParams{
		DiscordID: a.DiscordID,
		Price:     turnipPriceObj.Price,
		AmPm:      turnipPriceObj.AmPm,
		DayOfWeek: turnipPriceObj.DayOfWeek,
		DayOfYear: turnipPriceObj.DayOfYear,
		Year:      turnipPriceObj.Year,
	}

	q := turnips.New(db)
	_, err = q.UpdatePrice(ctx, priceParams)

	if err != nil {
		response.Emoji = "⛔"
		response.Text = "I didn't find an existing price."
	} else {
		response = turnipPriceColorfulResponse(turnipPrice)
	}

	flushEmojiAndResponseToDiscord(s, m, response)
	linkUsersCurrentPrices(s, m, AcTurnipsChartLink)
}

func buildPriceObjFromInput(a turnips.Account, turnipPrice int) (response, turnips.TurnipPrice, error) {
	accountTimeZone, err := time.LoadLocation(a.TimeZone)
	var response response

	if err != nil {
		response.Emoji = "⛔"
		response.Text = "Set a valid timezone from the `TZ database name` column https://en.wikipedia.org/wiki/List_of_tz_database_time_zones"
		return response, turnips.TurnipPrice{}, err
	}

	localTime := time.Now().In(accountTimeZone)
	var meridiem turnips.AmPm
	switch fmt.Sprint(localTime.Format("pm")) {
	case "am":
		meridiem = turnips.AmPmAm
	case "pm":
		meridiem = turnips.AmPmPm
	}

	_, week := localTime.ISOWeek()
	priceThing := turnips.TurnipPrice{
		DiscordID: a.DiscordID,
		Price:     int32(turnipPrice),
		AmPm:      meridiem,
		DayOfWeek: int32(localTime.Weekday()),
		DayOfYear: int32(localTime.YearDay()),
		Year:      int32(localTime.Year()),
		Week:      int32(week),
	}
	priceThing.AmPm = meridiem

	return response, priceThing, err
}

func turnipPriceColorfulResponse(turnipPrice int) response {
	var response response
	response.Emoji = "✅"
	if turnipPrice == 69 {
		response.Text = "nice."
	} else if turnipPrice > 0 && turnipPrice <= 100 {
		response.Text = fmt.Sprintf("Oh, your turnips are selling for %d right now? Sucks to be poor!", turnipPrice)
	} else if turnipPrice > 0 && turnipPrice <= 149 {
		response.Text = fmt.Sprintf("Oh, your turnips are selling for %d right now? Meh.", turnipPrice)
	} else if turnipPrice > 0 && turnipPrice <= 150 {
		response.Text = fmt.Sprintf("Oh, your turnips are selling for %d right now? Decent!", turnipPrice)
	} else if turnipPrice > 0 && turnipPrice < 200 {
		response.Text = fmt.Sprintf("Oh shit, your turnips are selling for %d right now? Dope!", turnipPrice)
	} else if turnipPrice >= 200 {
		response.Text = fmt.Sprintf("@everyone get in here! Someone has turnips trading for %d bells", turnipPrice)
	} else {
		response.Text = "Is that even a real number?"
		response.Emoji = "❌"
	}
	return response
}
