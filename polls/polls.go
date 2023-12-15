package polls

import (
	"NoiseDcBot/database"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

func CreatePoll(s *discordgo.Session, m *discordgo.MessageCreate) {
	date := time.Now().Format(time.DateTime)
	db := database.OpenDB()
	defer db.Close()

	parts := strings.Split(m.Content, "\"")

	question := parts[1]
	fmt.Println(parts[0])
	parts[0] = parts[0][13:]
	splitStr := strings.Split(parts[0], ":")

	var emojis []string
	var texts []string

	for i := 0; i < len(splitStr); i++ {
		if i%2 == 0 {
			emojis = append(emojis, splitStr[i])
		} else {
			texts = append(texts, splitStr[i])
		}
	}

	getId := fmt.Sprintf("SELECT id FROM user WHERE user_id='%s';", m.Author.ID)
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
	q := fmt.Sprintf("INSERT INTO polls (user, message, channel_id, reactions, date) VALUES ('%v', '%s', '%s', '%s', '%s');", id, m.Message.Content, m.ChannelID, emojis, date)
	_, err = db.Exec(q)
	if err != nil {
		fmt.Println(err)
	}

	question = fmt.Sprintf("**Question** \n %s \n \n **Coices** \n", question)

	for i, _ := range emojis {
		question = fmt.Sprintf("%s \n %s %s \n", question, emojis[i], texts[i])
	}

	question = fmt.Sprintf("\n %s \n **Settings** \n ", question)

	embed := &discordgo.MessageEmbed{
		Color:       0x000000, // Change to your desired color
		Description: question,
		Timestamp:   date,
	}

	pollMessage, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		fmt.Println("Error creating Message: ", err)
		return
	}
	for _, emoji := range emojis {
		err := s.MessageReactionAdd(m.ChannelID, pollMessage.ID, emoji)
		if err != nil {
			fmt.Println("Error creating Message: ", err)
		}
	}

	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		fmt.Printf("Error deleting message: %v", err)
	}

}
