package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func OnGuildMemberAdd(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	guild, err := session.Guild(event.GuildID)
	if err != nil {
		fmt.Println("Failed to get guild:", err)
		return
	}

	channel, err := session.UserChannelCreate(event.User.ID)
	if err != nil {
		fmt.Println("Failed to create DM channel:", err)
		return
	}

	message := fmt.Sprintf("Welcome to %s, %s!", guild.Name, event.User.Username)
	session.ChannelMessageSend(channel.ID, message)
}
