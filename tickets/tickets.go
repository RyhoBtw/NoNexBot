package tickets

import (
	"NoiseDcBot"
	"NoiseDcBot/database"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
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
	createSupportChannel(s, m, id)
}

func createSupportChannel(s *discordgo.Session, m *discordgo.MessageCreate, id int) {
	conf, err := NoiseDcBot.ReadBotConf("DBConf.yml")
	guildID := conf.GuildID
	channelName := fmt.Sprintf("supportTicket -", id)

	_, err = s.GuildChannelCreate(guildID, channelName, discordgo.ChannelTypeGuildText)
	if err != nil {
		fmt.Println("error creating channel,", err)
		return
	}
	//_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
	//	ParentID: conf.SupportCategory,
	//})
}
