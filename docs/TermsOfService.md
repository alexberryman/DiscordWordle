# Terms Of Service
By inviting this bot and using its features you are agreeing to the below-mentioned Terms and Privacy Policy.

You acknowledge that you are free to invite the bot to any server where you have permission to invite guests to that server, and that privilege may be revoked for you if you break the [Terms of Service](./TermsOfService.md).

Inviting the bot allows it to collect specific data as described in its [Privacy Policy](./TermsOfService.md). The intended usage of this data is the core functionalities of the bot such as parsing messages for Wordle scores, displaying Wordle scores, handling commands to change guild-specific settings. The Wordle Scores from multiple servers may be processed to create public stats (Average # of guesses, # of players, etc.) about Wordle games.

## Intended Age
The bot may not be used by individuals under the minimal age described in Discord's Terms of Service.

## Affiliation¶
The Bot is not affiliated with, supported or made by Discord Inc.
Any direct connection to Discord or any of its Trademark objects is purely coincidental. We do not claim to have the copyright ownership of any of Discord's assets, trademarks or other intellectual property.

# Liability
The owner of the bot may not be made liable for individuals breaking these Terms at any given time.
He has faith in the end users being truthful about their information and not misusing this bot or The Services of Discord Inc in a malicious way.

We reserve the right to update these terms at our own discretion, giving you a 1-Week (7 days) period to opt out of these terms if you're not agreeing with the new changes. You may opt out by Removing the bot from any Server you have the rights for.

# Privacy Policy
## Usage of Data
The bot may use stored data, as defined below, for different features including but not limited to: Parsing messages for Wordle scores, displaying Wordle scores, handling commands to change guild-specific settings. The Wordle Scores from multiple servers may be processed to create public stats (Average # of guesses, # of players, etc.) about Wordle games.

No usage of data outside of the aforementioned cases will happen and the data is not shared with any 3rd-party site or service.

## Stored Information
The bot may store the information automatically when being invited to a new Discord Server. See [SQL schema migrations](https://github.com/alexberryman/DiscordWordle/tree/main/internal/wordle/schema) for implementation of stored data. The following data is kept in the database:
- `server_id` to relate data to a specific Discord Server
- `discord_id` to relate data to a specific Discord User
- `nicknames` to display the Discord User's preferred name in the context of the server on the Wordle `scoreboard`
- `time_zone` to allow the bot to determine the local date/time for a Discord User
- `quip` user-generated content that the bot will use to respond to specific Wordle scores
- `quips.uses` count of how many times a `quip` has been used in order to increase the variety of responses 
- `wordle_scores.game_id` the game number shared when a Discord User posts a Wordle Score block
- `wordle_scores.guesses` the number of guesses taken to solve the Wordle puzzle shared when a Discord User posts a Wordle Score block

## Updating Data
The data may be updated when using specific commands.
Updating data will require the input of an end user, and data that can be seen as sensitive, such as content of a message, may need to be stored when using certain commands.

No other actions may update the stored information at any given time.

## Removal of Data
### Automatic removal
Not implemented at this time

### Manual removal
Manual removal of the data can be requested through Discord by messaging `BearsInTheSky#9588`. The Discord ID you use to message this user will be verified that you own the data you are asking to be removed.