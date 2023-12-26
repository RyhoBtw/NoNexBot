package polls

import (
	"NoiseDcBot/database"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

func CreatePoll(question string, maxChoices int, allowedRole *discordgo.Role, text string, anonymous bool, answers []string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	date := time.Now().Format(time.DateTime)
	db := database.OpenDB()
	defer db.Close()

	var emojis []string

	for _, answer := range answers {
		parts := strings.SplitN(answer, " ", 2)
		emojis = append(emojis, parts[0])
	}

	log.Println("Emojis", emojis)

	getId := fmt.Sprintf("SELECT id FROM user WHERE user_id='%s';", i.ID)
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
	q := fmt.Sprintf("INSERT INTO polls (user, message, channel_id, reactions, date) VALUES ('%v', '%s', '%s', '%s', '%s');", id, i.Data.Type().String(), i.ChannelID, emojis, date)
	_, err = db.Exec(q)
	if err != nil {
		fmt.Println(err)
	}

	question = fmt.Sprintf("**Question** \n %s \n \n **Coices** \n", question)

	for i, _ := range answers {
		question = fmt.Sprintf("%s \n %s \n", question, answers[i])
	}

	question = fmt.Sprintf("\n %s \n **Settings** \n ", question)

	embed := &discordgo.MessageEmbed{
		Title:       text,
		Color:       0x000000, // Change to your desired color
		Description: question,
		Timestamp:   date,
	}

	pollMessage, err := s.ChannelMessageSendEmbed(i.ChannelID, embed)
	if err != nil {
		fmt.Println("Error creating Message: ", err)
		return
	}
	for _, emoji := range emojis {
		err := s.MessageReactionAdd(i.ChannelID, pollMessage.ID, emoji)
		if err != nil {
			fmt.Println("Error creating Message: ", err)
		}
	}

}
