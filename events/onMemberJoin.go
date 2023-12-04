package events

import (
	"NoiseDcBot"
	"NoiseDcBot/database"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func OnGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd, c *NoiseDcBot.BotConf) {
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

	db := database.OpenDB()
	defer db.Close()

	q := fmt.Sprintf("INSERT INTO user (user_id, username, joinDate) VALUES ('%s', '%s', '%v');", userId, m.User.Username, time.Now().Format("2006-01-02"))
	_, err = db.Query(q)
	if err != nil {
		log.Println("ERROT: Insert %v", err)
	}
}
