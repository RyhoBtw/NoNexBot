package main

import (
	"NoiseDcBot/events"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type conf struct {
	Token string `yaml:"token"`
}

func main() {
	c, err := readConf("conf.yml")
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

	session.AddHandler(onReady)
	session.AddHandler(events.OnGuildMemberAdd)
	session.AddHandler(events.OnMssageCreate)

	err = session.Open()
	if err != nil {
		fmt.Println("Failed to open Discord session:", err)
		return
	}

	defer session.Close()

	select {}
}

func onReady(session *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

}

func readConf(filename string) (*conf, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &conf{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}
