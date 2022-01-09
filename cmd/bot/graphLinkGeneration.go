package main

import (
	"DiscordGoTurnips/internal/turnips/generated-code"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type dailyPrice struct {
	DayOfWeek      int
	MorningPrice   int32
	AfternoonPrice int32
}

type response struct {
	Text  string
	Emoji string
}

//AcTurnipsChartLink Formatting string for web links to ac-turnip
const AcTurnipsChartLink = "%s: <https://ac-turnip.com/share?f=%s>\n"

//AcTurnipsImageLink Formatting string for image links to ac-turnip
const AcTurnipsImageLink = "%s: https://ac-turnip.com/p-%s.png\n"

func linkServersCurrentPrices(s *discordgo.Session, m *discordgo.MessageCreate, linkFormat string) {
	q := turnips.New(db)
	ctx := context.Background()
	prices, err := q.GetWeeksPriceHistoryByServer(ctx, m.GuildID)
	if err != nil {
		log.Println("error fetching prices: ", err)
	}

	response := buildPriceGraphs(prices, linkFormat)
	flushEmojiAndResponseToDiscord(s, m, response)
}

func linkUsersCurrentPrices(s *discordgo.Session, m *discordgo.MessageCreate, linkFormat string) {
	var response response

	q := turnips.New(db)
	ctx := context.Background()
	prices, err := q.GetWeeksPriceHistoryByAccount(ctx, turnips.GetWeeksPriceHistoryByAccountParams{
		DiscordID: m.Author.ID,
		ServerID:  m.GuildID,
	})

	if err != nil {
		log.Println("error fetching prices: ", err)
	}

	data := make([]turnips.GetWeeksPriceHistoryByServerRow, 0)
	for _, v := range prices {
		p := turnips.GetWeeksPriceHistoryByServerRow(v)
		data = append(data, p)
	}

	response = buildPriceGraphs(data, linkFormat)

	flushEmojiAndResponseToDiscord(s, m, response)
}

func linkAccountsPreviousPrices(m *discordgo.MessageCreate, s *discordgo.Session, offset int, linkFormat string) {
	q := turnips.New(db)
	ctx := context.Background()

	week, err := getCurrentWeek(ctx, m, q)
	if err != nil {
		log.Println("error fetching current week: ", err)
	}

	prices, err := q.GetHistoricalWeekPriceHistoryByAccount(ctx, turnips.GetHistoricalWeekPriceHistoryByAccountParams{
		DiscordID: m.Author.ID,
		ServerID:  m.GuildID,
		Week:      int32(week - offset),
	})
	if err != nil {
		log.Println("error fetching prices: ", err)
	}

	data := make([]turnips.GetWeeksPriceHistoryByServerRow, 0)
	for _, v := range prices {
		p := turnips.GetWeeksPriceHistoryByServerRow(v)
		data = append(data, p)
	}

	response := buildPriceGraphs(data, linkFormat)
	flushEmojiAndResponseToDiscord(s, m, response)
}

func linkServersPreviousPrices(m *discordgo.MessageCreate, s *discordgo.Session, offset int, linkFormat string) {
	q := turnips.New(db)
	ctx := context.Background()

	week, err := getCurrentWeek(ctx, m, q)
	if err != nil {
		log.Println("error fetching current week: ", err)
	}

	prices, err := q.GetHistoricalWeekPriceHistoryByServer(ctx, turnips.GetHistoricalWeekPriceHistoryByServerParams{
		ServerID: m.GuildID,
		Week:     int32(week - offset),
	})
	if err != nil {
		log.Println("error fetching prices: ", err)
	}

	data := make([]turnips.GetWeeksPriceHistoryByServerRow, 0)
	for _, v := range prices {
		p := turnips.GetWeeksPriceHistoryByServerRow(v)
		data = append(data, p)
	}

	response := buildPriceGraphs(data, linkFormat)
	flushEmojiAndResponseToDiscord(s, m, response)
}

func getCurrentWeek(ctx context.Context, m *discordgo.MessageCreate, q *turnips.Queries) (int, error) {
	account, _ := q.GetAccount(ctx, m.Author.ID)
	accountTimeZone, err := time.LoadLocation(account.TimeZone)
	localTime := time.Now().In(accountTimeZone)
	_, week := localTime.ISOWeek()
	return week, err
}

func buildPriceGraphs(prices []turnips.GetWeeksPriceHistoryByServerRow, format string) response {

	return buildPriceUri(prices, format)
}

func buildPriceUri(prices []turnips.GetWeeksPriceHistoryByServerRow, format string) response {
	priceMap := make(map[string]map[string]dailyPrice)

	for _, value := range prices {
		wp := getEmptyWeeklyPrices()
		if _, ok := priceMap[value.Nickname]; ok {
			updateMorningOrAfterNoonPrice(value, priceMap)
		} else {
			priceMap[value.Nickname] = wp
			updateMorningOrAfterNoonPrice(value, priceMap)
		}
	}

	var response response
	turnipLink := make(map[string]string)
	for nickname, prices := range priceMap {
		for _, d := range dayRange(monday, saturday) {
			if _, ok := turnipLink[nickname]; !ok {
				turnipLink[nickname] = ""
			}

			if prices[fmt.Sprint(d)].MorningPrice != 0 {
				turnipLink[nickname] += fmt.Sprintf("-%d", prices[fmt.Sprint(d)].MorningPrice)
			} else {
				turnipLink[nickname] += "-"
			}
			if prices[fmt.Sprint(d)].AfternoonPrice != 0 {
				turnipLink[nickname] += fmt.Sprintf("-%d", prices[fmt.Sprint(d)].AfternoonPrice)
			} else {
				turnipLink[nickname] += "-"
			}
		}
		response.Text += fmt.Sprintf(format, nickname, turnipLink[nickname])
	}

	response.Emoji = "🔗"
	return response
}

func dayRange(min, max Weekday) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = int(min) + i
	}
	return a
}

func getEmptyWeeklyPrices() map[string]dailyPrice {
	w := newWeeklyPrices()

	for _, d := range dayRange(monday, saturday) {
		dp := dailyPrice{
			DayOfWeek:      d,
			MorningPrice:   0,
			AfternoonPrice: 0,
		}
		w[fmt.Sprintf("%d", d)] = dp
	}
	return w
}

func newWeeklyPrices() map[string]dailyPrice {
	w := make(map[string]dailyPrice)
	return w
}

func updateMorningOrAfterNoonPrice(value turnips.GetWeeksPriceHistoryByServerRow, priceMap map[string]map[string]dailyPrice) {
	if value.AmPm == turnips.AmPmAm {
		tempPrice := priceMap[value.Nickname][fmt.Sprint(value.DayOfWeek)]
		tempPrice.MorningPrice = value.Price
		priceMap[value.Nickname][fmt.Sprint(value.DayOfWeek)] = tempPrice
	} else {
		tempPrice := priceMap[value.Nickname][fmt.Sprint(value.DayOfWeek)]
		tempPrice.AfternoonPrice = value.Price
		priceMap[value.Nickname][fmt.Sprint(value.DayOfWeek)] = tempPrice
	}
}
