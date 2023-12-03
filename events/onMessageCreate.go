package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func OnMssageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the message is a poll
	if strings.HasPrefix(m.Content, "!poll") {
		// Split the message content into the poll question and the reactions
		parts := strings.Split(m.Content, "\"")
		pollQuestion := parts[1]
		reactions := strings.Split(parts[3], " ")

		// Send the poll question
		_, _ = s.ChannelMessageSend(m.ChannelID, pollQuestion)

		// Add the reactions for the poll options
		for _, reaction := range reactions {
			_ = s.MessageReactionAdd(m.ChannelID, m.ID, reaction)
		}
	}

	// Check if the message is a poll result request
	if m.Content == "!results" {
		// Get the message with the poll
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
