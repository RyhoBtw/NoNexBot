package events

import (
	"NoiseDcBot"
	"NoiseDcBot/database"
	"NoiseDcBot/tickets"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func OnMssageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	conf, err := NoiseDcBot.ReadBotConf("conf.yml")
	if err != nil {
		fmt.Println(err)
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!ticket create") {
		tickets.CreateTicket(s, m)
	}

	guild, err := s.Guild(conf.GuildID)

	for _, channel := range guild.Channels {
		if channel.ParentID == conf.SupportCategory {
			supportID := channel.Topic
			db := database.OpenDB()
			defer db.Close()
			q := fmt.Sprintf()
		}
	}

	if strings.HasPrefix(m.Content, "!poll") {
		parts := strings.Split(m.Content, "\"")
		pollQuestion := parts[1]
		reactions := strings.Split(parts[3], " ")

		_, _ = s.ChannelMessageSend(m.ChannelID, pollQuestion)

		for _, reaction := range reactions {
			_ = s.MessageReactionAdd(m.ChannelID, m.ID, reaction)
		}
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
