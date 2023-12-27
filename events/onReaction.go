package events

import (
	"NoiseDcBot/database"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

type data struct {
	pollID     int
	reactions  string
	maxchoices int
	role       string
	anonymus   int8
}

var lastReaction bool

func OnReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	db := database.OpenDB()
	defer db.Close()

	if r.UserID == s.State.User.ID {
		return
	}

	msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Println("Error fetching message: ", err)
		return
	}

	q := fmt.Sprintf("SELECT id, reactions, maxchoices, role, anonymus FROM polls WHERE message_id='%s'", msg.ID)
	rows, err := db.Query(q)
	if err != nil {
		log.Println("Error fetching:", err)
	}

	var data data
	for rows.Next() {
		err := rows.Scan(&data.pollID, &data.reactions, &data.maxchoices, &data.role, &data.anonymus)
		if err != nil {
			log.Println("Error fetching:", err)
		}
	}

	if data.pollID == 0 {
		return
	}

	var allowed = false

	for _, role := range r.Member.Roles {
		if role == data.role {
			allowed = true
		}
	}

	if !allowed {
		lastReaction = true
		err = s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
		if err != nil {
			log.Println("Error removing reaction: ", err)
		}
		return
	}

	q = fmt.Sprintf("SELECT id FROM user WHERE user_id='%s'", r.UserID)
	rows, err = db.Query(q)
	if err != nil {
		log.Println("Error fetching:", err)
	}

	var userID int
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			log.Println("Error fetching:", err)
		}
	}

	var exists bool
	q = fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM reactions WHERE user=%v AND poll=%v)", userID, data.pollID)
	row := db.QueryRow(q)
	err = row.Scan(&exists)
	if err != nil {
		log.Println(err)
		return
	}

	var count int

	if exists {
		q = fmt.Sprintf("SELECT reactioncount FROM reactions WHERE user = %v AND poll = %v", userID, data.pollID)
		rows, err = db.Query(q)
		if err != nil {
			log.Println(err)
		}

		for rows.Next() {
			err = rows.Scan(&count)
			if err != nil {
				log.Println("Error fetching:", err)
			}
		}

		if count == data.maxchoices {
			lastReaction = true
			err = s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
			if err != nil {
				log.Println("Error removing reaction: ", err)
			}
		} else {
			count++
			q := fmt.Sprintf("UPDATE reactions SET reactioncount=%v WHERE user=%v AND poll=%v", count, userID, data.pollID)
			_, err = db.Query(q)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		q := fmt.Sprintf("INSERT INTO reactions (user, poll, reactioncount) VALUES (%v, %v, %v);", userID, data.pollID, 1)
		_, err = db.Query(q)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Reactions: ", data.reactions, "maxchoices: ", data.maxchoices, "role", data.role, "anonymus", data.anonymus)
}

func OnReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	db := database.OpenDB()
	defer db.Close()

	if lastReaction {
		lastReaction = false
		return
	}

	msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Println("Error fetching message: ", err)
		return
	}

	q := fmt.Sprintf("SELECT id, reactions, maxchoices, role, anonymus FROM polls WHERE message_id='%s'", msg.ID)
	rows, err := db.Query(q)
	if err != nil {
		log.Println("Error fetching:", err)
	}

	var data data
	for rows.Next() {
		err := rows.Scan(&data.pollID, &data.reactions, &data.maxchoices, &data.role, &data.anonymus)
		if err != nil {
			log.Println("Error fetching:", err)
		}
	}

	if data.pollID == 0 {
		return
	}

	q = fmt.Sprintf("SELECT id FROM user WHERE user_id='%s'", r.UserID)
	rows, err = db.Query(q)
	if err != nil {
		log.Println("Error fetching:", err)
	}

	var userID int
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			log.Println("Error fetching:", err)
		}
	}

	var count int
	q = fmt.Sprintf("SELECT reactioncount FROM reactions WHERE user = %v AND poll = %v", userID, data.pollID)
	rows, err = db.Query(q)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println("Error fetching:", err)
		}
	}

	count--
	q = fmt.Sprintf("UPDATE reactions SET reactioncount=%v WHERE user=%v AND poll=%v", count, userID, data.pollID)
	_, err = db.Query(q)
	if err != nil {
		log.Println(err)
	}
}
