package tickets

import (
	"NoiseDcBot"
	"NoiseDcBot/database"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func CreateTicket(text string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	userId := i.Member.User.ID
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

	q := fmt.Sprintf("INSERT INTO tickets (user, message, closed) VALUES (%v, '%s', %t);", id, text, false)
	_, err = db.Query(q)
	if err != nil {
		log.Println(err)
	}
	q = fmt.Sprintf("SELECT id FROM tickets WHERE message='%s'", text)
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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Created Channel about your Ticket",
		},
	})

	createSupportChannel(s, i, id, text)
}

func createSupportChannel(s *discordgo.Session, i *discordgo.InteractionCreate, id int, text string) {
	conf, err := NoiseDcBot.ReadBotConf("conf.yml")
	guildID := conf.GuildID
	channelName := fmt.Sprintf("supportTicket-%v", id)

	channel, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelTypeGuildText,
		Topic:    fmt.Sprintf("%v", id),
		ParentID: conf.SupportCategory,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{ID: conf.JoinRole,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Deny:  discordgo.PermissionViewChannel,
				Allow: 0},
		},
	})
	if err != nil {
		log.Println("Error creating channel: ", err)
		return
	}

	err = s.ChannelPermissionSet(channel.ID, i.Member.User.ID, discordgo.PermissionOverwriteTypeMember, discordgo.PermissionViewChannel, 0)
	if err != nil {
		return
	}
	var roleID string
	guild, err := s.Guild(guildID)
	for _, role := range guild.Roles {
		if role.ID == conf.TeamRole {
			roleID = role.ID
		}
	}
	for _, member := range guild.Members {
		if member.User.ID != i.Member.User.ID {
			for _, role := range member.Roles {
				if role == roleID {
					err = s.ChannelPermissionSet(channel.ID, member.User.ID, discordgo.PermissionOverwriteTypeMember, discordgo.PermissionViewChannel, 0)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}

	db := database.OpenDB()
	defer db.Close()
	q := fmt.Sprintf("UPDATE tickets SET channel_id = '%s' WHERE id = %v;", channel.ID, id) // doesn't work?
	test, err := db.Query(q)
	fmt.Println(test)
	if err != nil {
		log.Println(err)
	}

	message := fmt.Sprintf("<@%s> createed a ticket %s", i.Member.User.ID, text)

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

func CloseTicket(s *discordgo.Session, i *discordgo.InteractionCreate) {
	conf, err := NoiseDcBot.ReadBotConf("conf.yml")
	channelID := i.ChannelID
	db := database.OpenDB()
	defer db.Close()

	for _, role := range i.Member.Roles {
		if role == conf.TeamRole {
			var exists bool
			q := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM tickets WHERE channel_id = '%s')", channelID)
			row := db.QueryRow(q)
			err = row.Scan(&exists)
			if err != nil {
				log.Println(err)
				return
			}

			if exists {
				q = fmt.Sprintf("UPDATE tickets SET closed = 1 WHERE channel_id = '%s';", channelID)
				_, err = db.Exec(q)
				if err != nil {
					log.Println(err)
				}
				_, err = s.ChannelDelete(channelID)
				if err != nil {
					log.Println(err)
				}
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						// Note: this isn't documented, but you can use that if you want to.
						// This flag just allows you to create messages visible only for the caller of the command
						// (user who triggered the command)
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "This is not a Ticket channel!",
					},
				})
			}
		}
	}
	/*
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// Note: this isn't documented, but you can use that if you want to.
				// This flag just allows you to create messages visible only for the caller of the command
				// (user who triggered the command)
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "This is not a Ticket channel!",
			},
		})
		msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Followup message has been created, after 5 seconds it will be edited",
		})
		if err != nil {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong",
			})
			return
		}
		time.Sleep(time.Second * 5)

		content := "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted."
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: &content,
		})*/
}
