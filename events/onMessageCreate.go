package events

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"strings"
)

const KuteGoAPIURL = "https://kutego-api-xxxxxx-ew.a.run.app"

type Gopher struct {
	Name string `json: "name"`
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func OnMssageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!gopher" {

		//Call the KuteGo API and retrieve our cute Dr Who Gopher
		response, err := http.Get(KuteGoAPIURL + "/gopher/" + "dr-who")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			_, err = s.ChannelFileSend(m.ChannelID, "dr-who.png", response.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get dr-who Gopher! :-(")
		}
	}

	if m.Content == "!random" {

		//Call the KuteGo API and retrieve a random Gopher
		response, err := http.Get(KuteGoAPIURL + "/gopher/random/")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			_, err = s.ChannelFileSend(m.ChannelID, "random-gopher.png", response.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get random Gopher! :-(")
		}
	}

	if m.Content == "!gophers" {

		//Call the KuteGo API and display the list of available Gophers
		response, err := http.Get(KuteGoAPIURL + "/gophers/")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			// Transform our response to a []byte
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}

			// Put only needed informations of the JSON document in our array of Gopher
			var data []Gopher
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println(err)
			}

			// Create a string with all of the Gopher's name and a blank line as separator
			var gophers strings.Builder
			for _, gopher := range data {
				gophers.WriteString(gopher.Name + "\n")
			}

			// Send a text message with the list of Gophers
			_, err = s.ChannelMessageSend(m.ChannelID, gophers.String())
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get list of Gophers! :-(")
		}
	}
}
