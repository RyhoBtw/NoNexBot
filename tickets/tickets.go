package tickets

import (
	"NoiseDcBot"
	"NoiseDcBot/database"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

func CreateTicket(s *discordgo.Session, m *discordgo.MessageCreate) {
	userId := m.Author.ID
	messageContent := strings.TrimPrefix(m.Message.Content, "!ticket create")
	db := database.OpenDB()
	defer db.Close()
	getId := fmt.Sprintf("SELECT id FROM user WHERE user_id='%s';", userId)
	rows, err := db.Query(getId)
	if err != nil {
		log.Println(err)
	}

	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			panic(err)
		}
	}

	q := fmt.Sprintf("INSERT INTO tickets (user, channel_id, message, closed) VALUES (%v, '%s', '%s', %t);", id, m.ChannelID, messageContent, false)
	_, err = db.Query(q)
	if err != nil {
		log.Println(err)
	}
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Println(err)
	}
	q = fmt.Sprintf("SELECT id FROM tickets WHERE message='%s'", messageContent)
	rows, err = db.Query(q)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			panic(err)
		}
	}
	createSupportChannel(s, m, id, messageContent)
}

func createSupportChannel(s *discordgo.Session, m *discordgo.MessageCreate, id int, text string) {
	conf, err := NoiseDcBot.ReadBotConf("conf.yml")
	guildID := conf.GuildID
	channelName := fmt.Sprintf("supportTicket-%v", id)

	channel, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelTypeGuildText,
		Topic:    "Ticket",
		ParentID: conf.SupportCategory,
	})
	if err != nil {
		log.Println("Error creating channel: ", err)
		return
	}

	err = s.ChannelPermissionSet(channel.ID, m.Author.ID, discordgo.PermissionOverwriteTypeMember, discordgo.PermissionViewChannel, 0)
	if err != nil {
		return
	}
	err = s.ChannelPermissionSet(channel.ID, m.Author.ID, discordgo.PermissionOverwriteTypeMember, discordgo.PermissionSendMessages, 0)
	if err != nil {
		return
	}

	/*for _, role := range guild.Roles {
		if role.Name == "Admin" || role.Name == "Moderator" {
			guild.ChannelPermissionOverwrites(channel.ID).Add(discordgo.RoleID(role.ID), discordgo.Permission{
				VIEW_CHANNEL:  true,
				SEND_MESSAGES: true,
			})
		}
	}*/

	message := fmt.Sprintf("<@%s> createed a ticket about: %s", m.Author.ID, text)

	embed := &discordgo.MessageEmbed{
		Color:       0x000000, // Change to your desired color
		Title:       "Welcome to your support channel!",
		Description: message,
		//Image:       &discordgo.MessageEmbedImage{URL: m.Author.AvatarURL("original")},
		Timestamp: time.Now().UTC().Format(time.DateTime),
	}

	_, err = s.ChannelMessageSendEmbed(channel.ID, embed)
	if err != nil {
		log.Println("Error creating Message: ", err)
		return
	}

}
