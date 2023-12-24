package events

import (
	"NoiseDcBot/polls"
	"github.com/bwmarrin/discordgo"
)

func OnInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the interaction is a command
	if i.Type == discordgo.InteractionApplicationCommand {
		// Check if the command name is "poll"
		if i.ApplicationCommandData().Name == "poll" {
			var question string
			var maxChoices int
			var allowedRole *discordgo.Role
			var text string
			var anonymous bool
			var answers []string

			// Iterate over all options
			for _, opt := range i.ApplicationCommandData().Options {
				switch opt.Name {
				case "question":
					question = opt.StringValue()
				case "maxchoices":
					maxChoices = int(opt.IntValue())
				case "allowedrole":
					allowedRole = opt.RoleValue(s, i.GuildID)
				case "text":
					text = opt.StringValue()
				case "anonymous":
					anonymous = opt.BoolValue()
				default:
					// Assume any other option is an answer
					if opt.StringValue() != "" {
						answers = append(answers, opt.StringValue())
					}
				}
			}

			polls.CreatePoll(question, maxChoices, allowedRole, text, anonymous, answers, s, i)
		}
	}
}
