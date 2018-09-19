package main

import (
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/schollz/closestmatch"
)

type MessageHandler struct {
	skills SkillSet
	matcher *closestmatch.ClosestMatch
}

func NewMessageHandler(file string) (*MessageHandler, error) {
	s, matcher, err := LoadSkills(file)
	if err != nil {
		return nil, err
	}
	return &MessageHandler{
		skills: s,
		matcher: matcher,
	}, nil
}

func (mh *MessageHandler) Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	params := strings.Split(strings.ToLower(m.Content), " ")
	fmt.Println(params)
	switch params[0]{
	case "skill":
		mh.SkillHandler(params[1:], s, m)
	}
}

func (mh *MessageHandler) SkillHandler(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	skill_val := strings.Join(args, " ")
	match := mh.matcher.Closest(skill_val)
	if skill, ok := mh.skills[match]; ok {
		s.ChannelMessageSend(m.ChannelID, "```" + skill.Desc + "```")
	}
}
