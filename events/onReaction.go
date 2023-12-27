package events

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func OnReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.MessageID == "pollMessageID" {
		reactions, err := s.MessageReactions(r.ChannelID, r.MessageID, "", 100, "", "")
		if err != nil {
			log.Println("Error getting reactions:", err)
			return
		}

		count := len(reactions)
}
