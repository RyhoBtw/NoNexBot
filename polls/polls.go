package polls

import (
	"NoiseDcBot"
	"NoiseDcBot/database"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

func CreatePoll(question string, maxChoices int, allowedRole *discordgo.Role, text string, anonymous bool, answers []string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	conf, err := NoiseDcBot.ReadBotConf("conf.yml")
	date := time.Now().Format(time.DateTime)
	db := database.OpenDB()
	defer db.Close()

	var emojis []string

	for _, answer := range answers {
		parts := strings.SplitN(answer, " ", 2)
		emojis = append(emojis, parts[0])
	}

	log.Println("Emojis", emojis)

	getId := fmt.Sprintf("SELECT id FROM user WHERE user_id='%s';", i.Member.User.ID)
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

	if maxChoices == 0 {
		maxChoices = 1
	}

	var anonym int
	if anonymous {
		anonym = 1
	} else {
		anonym = 0
	}

	var roleID string
	if allowedRole == nil {
		roleID = conf.JoinRole
	} else {
		roleID = allowedRole.ID
	}

	q := fmt.Sprintf("INSERT INTO polls (user, message, channel_id, reactions, date, maxchoices, role, anonymus) VALUES ('%v', '%s', '%s', '%s', '%s', '%v', '%s', '%v');", id, i.Data.Type().String(), i.ChannelID, emojis, date, maxChoices, roleID, anonym)
	_, err = db.Exec(q)
	if err != nil {
		fmt.Println(err)
	}

	question = fmt.Sprintf("**Question** \n %s \n \n **Coices** \n", question)

	for i, _ := range answers {
		question = fmt.Sprintf("%s \n %s \n", question, answers[i])
	}

	question = fmt.Sprintf("\n %s \n **Settings** \n", question)

	question = fmt.Sprintf("%s Max choises: %v \n", question, maxChoices)
	question = fmt.Sprintf("%s Anonym: %t \n", question, anonymous)
	question = fmt.Sprintf("%s allowed Role: <@&%s> \n", question, roleID)

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
