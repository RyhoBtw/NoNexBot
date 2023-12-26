package main

import (
	"NoiseDcBot"
	"NoiseDcBot/database"
	"NoiseDcBot/events"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func main() {

	_, err := NoiseDcBot.ReadDBConf("DBConf.yml")
	if err != nil {
		log.Println(err)
	}

	_ = database.ConnectToDatabase()

	c, err := NoiseDcBot.ReadBotConf("conf.yml")
	if err != nil {
		log.Println(err)
	}
	if c.Token == "" {
		fmt.Println("DISCORD_BOT_TOKEN environment variable not set.")
		return
	}

	session, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		fmt.Println("Failed to create Discord session:", err)
		return
	}

	_, err = session.ApplicationCommandCreate("1106901093729964153", c.GuildID, NoiseDcBot.Poll())
	if err != nil {
		log.Println("Failed:", err)
	}
	_, err = session.ApplicationCommandCreate("1106901093729964153", c.GuildID, NoiseDcBot.CreateTicket())
	if err != nil {
		log.Println("Failed:", err)
	}
	_, err = session.ApplicationCommandCreate("1106901093729964153", c.GuildID, NoiseDcBot.CloseTicket())
	if err != nil {
		log.Println("Failed:", err)
	}

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	session.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		events.OnGuildMemberAdd(s, m, c)
	})
	session.AddHandler(onReady)
	session.AddHandler(events.OnInteraction)

	err = session.Open()
	if err != nil {
		fmt.Println("Failed to open Discord session:", err)
		return
	}

	defer session.Close()

	select {}
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	/*_, err := s.ChannelMessageSend("823971053545062474", "Bot is online")
	if err != nil {
		log.Println(err)
	}*/

	go deletChannel(s)
}

func deletChannel(s *discordgo.Session) {
	c, err := NoiseDcBot.ReadBotConf("conf.yml")
	if err != nil {
		log.Println(err)
	}
	channelID := c.StreamChannel
	channel, err := s.Channel(channelID)
	if err != nil {
		log.Println("Channel with ID:", channelID, "not Found")
	}

	for {
		day := time.Now().Day()
		if day == 1 {
			messages, err := s.ChannelMessages(channelID, channel.MessageCount, "", "", "")
			if err != nil {
				log.Println(err)
			}

			messageIDs := make([]string, len(messages))

			for i, message := range messages {
				messageIDs[i] = message.ID
			}

			for i, _ := range messages {
				err := s.ChannelMessageDelete(channelID, messageIDs[i])
				if err != nil {
					log.Println(err)
				}
			}
		}

		time.Sleep(24 * time.Hour)
	}
}
