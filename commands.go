package NoiseDcBot

import "github.com/bwmarrin/discordgo"

func CreateTicket() (ticket *discordgo.ApplicationCommand) {
	ticket = &discordgo.ApplicationCommand{
		Name:        "ticketcreate",
		Description: "Creates a Support ticket",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "What is your ticket about?",
				Required:    true,
			},
		},
	}
	return ticket
}

func CloseTicket() (ticket *discordgo.ApplicationCommand) {
	ticket = &discordgo.ApplicationCommand{
		Name:        "ticketclose",
		Description: "Closes a Support ticket",
	}
	return ticket
}

func Poll() (poll *discordgo.ApplicationCommand) {
	poll = &discordgo.ApplicationCommand{
		Name:        "poll",
		Description: "Creates a normal Poll",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "question",
				Description: "What is the question?",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "maxchoices",
				Description: "How many choices are allowed per user? (Default 1)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "allowedrole",
				Description: "Allowed role which can vote",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "Text to appear above the poll (introduction, ping, role,...)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer1",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer2",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer3",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer4",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer5",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer6",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer7",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer8",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer9",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer10",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer12",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer13",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer14",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer15",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "answer16",
				Description: "Put emoji first than put awnser option",
				Required:    false,
			},
		},
	}
	return poll
}
