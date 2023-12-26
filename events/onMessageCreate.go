package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func OnMssageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!results" {
		message, _ := s.ChannelMessage(m.ChannelID, m.ID)

		// Check if the message has reactions
		if len(message.Reactions) > 0 {
			// Loop through the reactions and print the results
			for _, reaction := range message.Reactions {
				fmt.Printf("%s: %d\n", reaction.Emoji.Name, reaction.Count)
			}
		}
	}
}
