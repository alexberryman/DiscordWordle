package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func helpResponse(s *discordgo.Session, m *discordgo.MessageCreate, botMentionToken string, cmdHistory string, CmdTimeZone string, cmdQuip string) {
	var response response
	response.Text = fmt.Sprintf("`%s` - register a score for the current Wordle game. Only one is score is allowed for each game.\n"+
		"`%s` - update existing Wordle score for a game\n"+
		"`%s` - get your past Wordle scores\n"+
		"`%s` - Add your own sass for the bot to use as a reply for a specific number\n"+
		"`%s` - set your local timezone <https://en.wikipedia.org/wiki/List_of_tz_database_time_zones>\n"+
		"Report issue or help improve this bot at <https://github.com/alexberryman/DiscordWordle>\n",
		fmt.Sprintf("%s Wordle 204 5/6 <emoji blocks>", botMentionToken),
		fmt.Sprintf("%s update 204 2/6", botMentionToken),
		fmt.Sprintf("%s %s", botMentionToken, cmdHistory),
		fmt.Sprintf("%s %s 3 the three blind mice could do better", botMentionToken, cmdQuip),
		fmt.Sprintf("%s %s America/New_York", botMentionToken, CmdTimeZone),
	)

	response.Emoji = "💁"

	flushEmojiAndResponseToDiscord(s, m, response)
}
