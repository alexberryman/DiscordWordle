package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func helpResponse(s *discordgo.Session, m *discordgo.MessageCreate, botMentionToken string) {
	var response response
	response.Text = fmt.Sprintf("`%s` - register a score for the current Wordle game. Only one is score is allowed for each game.\n"+
		"`%s` - update existing Wordle score for a game\n"+
		"`%s` - get your past Wordle scores\n"+
		"`%s` - Turn off quip responses, the bot will continue to add emojis to the request if it understood the input\n"+
		"`%s` - Turn quip responses back on if previously disabled\n"+
		"`%s` - Add your own sass for the bot to use as a reply for a specific number\n"+
		"`%s` - view the scoreboard of you and your friends\n"+
		"`%s` - view last week's (game number/7) scoreboard\n"+
		"Join the DiscordWordle server for help with the bot: https://discord.gg/PcZ74rBsTE\n"+
		"Report issue or help improve this bot at <https://github.com/alexberryman/DiscordWordle>\n",
		fmt.Sprintf("%s %s 204 5/6 <emoji blocks>", botMentionToken, cmdWordle),
		fmt.Sprintf("%s %s 204 2/6", botMentionToken, cmdUpdate),
		fmt.Sprintf("%s %s", botMentionToken, cmdHistory),
		fmt.Sprintf("%s %s %s", botMentionToken, cmdQuip, cmdQuipDisable),
		fmt.Sprintf("%s %s %s", botMentionToken, cmdQuip, cmdQuipEnable),
		fmt.Sprintf("%s %s 3 Wow, you seem really smart!", botMentionToken, cmdQuip),
		fmt.Sprintf("%s %s", botMentionToken, cmdScoreboard),
		fmt.Sprintf("%s %s %s", botMentionToken, cmdScoreboard, cmdPreviousWeek),
	)

	response.Emoji = "💁"

	flushEmojiAndResponseToDiscord(s, m, response)
}
