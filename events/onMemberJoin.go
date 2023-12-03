package events

import (
	"NoiseDcBot"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func OnGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd, c *NoiseDcBot.Conf) {
	userId := m.User.ID

	message := fmt.Sprintf(c.JoinMessage, userId)

	_, err := s.ChannelMessageSend(c.JoinChannel, message)
	if err != nil {
		fmt.Println("error sending message,", err)
	}

	roleName := c.JoinRole
	role, err := s.State.Role(m.GuildID, roleName)
	if err != nil {
		log.Println("error getting role:", err)
		return
	}

	// Add the role to the user
	err = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, role.ID)
	if err != nil {
		log.Println("error adding role to user:", err)
	}
}
