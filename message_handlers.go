package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/schollz/closestmatch"
	"strings"
)

type MessageHandler struct {
	skills       SkillSet
	skillMatcher *closestmatch.ClosestMatch
}

func NewMessageHandler(file string) (*MessageHandler, error) {
	s, skillMatcher, err := LoadSkills(file)
	if err != nil {
		return nil, err
	}
	return &MessageHandler{
		skills:       s,
		skillMatcher: skillMatcher,
	}, nil
}

func (mh *MessageHandler) Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	params := strings.Split(strings.ToLower(m.Content), " ")
	switch params[0] {
	case "bothelp":
		mh.HelpHandler(s, m)
	case "skill":
		mh.SkillHandler(params[1:], s, m)
	case "links":
		mh.LinkHandler(params[1:], s, m)
	case `\o/`:
		mh.PraiseChorus(s, m)
	}
}

func (mh *MessageHandler) sendSkill(s *discordgo.Session, m *discordgo.MessageCreate, skill *Skill, embed bool) {
	if embed {
		cpcost := &discordgo.MessageEmbedField{
			Name:   "Cost",
			Value:  strings.Replace(skill.Cost+" CP", "* CP", " CP Each", 1),
			Inline: true,
		}
		prereq := &discordgo.MessageEmbedField{
			Name:   "Prerequisites",
			Value:  skill.Prereq,
			Inline: true,
		}
		description := &discordgo.MessageEmbedField{
			Name:  "Description",
			Value: skill.Desc,
		}
		embed := &discordgo.MessageEmbed{
			Title:  skill.Name,
			Color:  0x00c24d,
			Fields: []*discordgo.MessageEmbedField{cpcost, prereq, description},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	} else {
		msg := fmt.Sprintf(
			":hourglass:\n**%s**\n**Prerequisites:** %s\n**CP Cost:** %s\n```Elm\n%s```",
			skill.Name, skill.Prereq, skill.Cost, skill.Desc,
		)
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func (mh *MessageHandler) HelpHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title: "Twin Mask Bot Commands",
	}
	skills := &discordgo.MessageEmbedField{
		Name:  "Look Up Skill",
		Value: "skill [skill name to look up]",
	}
	links := &discordgo.MessageEmbedField{
		Name:  "Commonly Referenced Links",
		Value: "links list",
	}
	embed.Fields = []*discordgo.MessageEmbedField{skills, links}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func (mh *MessageHandler) PraiseChorus(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, ":innocent:\n**PRAISE CHORUS! PRAISE HIS ARDENTS!**")
}

func (mh *MessageHandler) SkillHandler(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	skill_val := strings.Join(args, " ")
	match := mh.skillMatcher.Closest(skill_val)
	if skill, ok := mh.skills[match]; ok {
		mh.sendSkill(s, m, skill, true)
	}
}

func (mh *MessageHandler) LinkHandler(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(args) == 0 {
		return
	}
	embed := &discordgo.MessageEmbed{
		Color: 0x00c24d,
	}
	switch args[0] {
	case "list":
		embed.Title = "List of Links"
		goa := &discordgo.MessageEmbedField{
			Name:  "Guild of Academics Catalogue",
			Value: "[Digitally Available Guild of Academics Documents](https://drive.google.com/open?id=1xVV7mwnWogRuMxwLWfoZUfGbR4Cl5RJS)\nShortcut: links goa",
		}
		rules := &discordgo.MessageEmbedField{
			Name:  "Twin Mask Rulebook",
			Value: "[The Official Rulebook for Twin Mask.](https://docs.wixstatic.com/ugd/c66b5f_8b6ca9b3982b43a2aedcc49182338855.pdf)\nShortcut: links rulebook",
		}
		wiki := &discordgo.MessageEmbedField{
			Name:  "Twin Mask Wiki",
			Value: "[The Official Wiki for Twin Mask.](http://twin-mask.wikia.com/)\nShortcut: links wiki",
		}
		forums := &discordgo.MessageEmbedField{
			Name:  "Old Proboard Forums",
			Value: "[The Old Twin Mask Forums](http://adelrune.proboards.com/)\nShortcut: links forum",
		}
		embed.Fields = []*discordgo.MessageEmbedField{rules, wiki, forums, goa}
	case "goa":
		embed.Title = "Guild of Academics Catalogue"
		embed.URL = "https://drive.google.com/open?id=1xVV7mwnWogRuMxwLWfoZUfGbR4Cl5RJS"
		embed.Description = "A catalogue of every Guild of Academics document that is currently available online."
	case "rules", "rulebook":
		embed.Title = "Twin Mask Rulebook"
		embed.URL = "https://docs.wixstatic.com/ugd/c66b5f_8b6ca9b3982b43a2aedcc49182338855.pdf"
		embed.Description = "The most up-to-date Twin Mask Rulebook."
	case "wiki":
		embed.Title = "Twin Mask Wiki Page"
		embed.URL = "http://twin-mask.wikia.com"
		embed.Description = "A wiki page for Twin mask where you can find, among other information, player-made character pages."
	case "forums":
		embed.Title = "Twin Maks Proboard Forum"
		embed.URL = "http://adelrune.proboards.com"
		embed.Description = "The old Twin Mask forums, which are filled with a lot of information from games long past. A must-peruse for anybody ready to dig deep into the Twin Mask universe."
	default:
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
