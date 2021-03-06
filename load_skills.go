package main

import (
	"encoding/csv"
	"fmt"
	"github.com/schollz/closestmatch"
	"io"
	"os"
	"strings"
)

type Skill struct {
	Name   string
	Cost   string
	Prereq string
	Desc   string
}

type SkillSet map[string]*Skill

func recordToSkill(record []string) (*Skill, error) {
	if len(record) != 4 {
		fmt.Printf("%+v\n", record)
		return nil, fmt.Errorf("Record not a skill")
	}
	return &Skill{
		Name:   record[0],
		Cost:   record[1],
		Prereq: record[2],
		Desc:   record[3],
	}, nil
}

func LoadSkills(file string) (SkillSet, *closestmatch.ClosestMatch, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	toret := map[string]*Skill{}
	csvr := csv.NewReader(f)
	csvr.Comment = '#'
	skillNames := []string{}
	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		skill, err := recordToSkill(record)
		if err != nil {
			return nil, nil, err
		}
		key := strings.ToLower(skill.Name)
		skillNames = append(skillNames, key)
		toret[key] = skill
		if strings.HasSuffix(key, ": apprentice") {
			key_p := strings.TrimSuffix(key, ": apprentice")
			toret[key_p] = skill
			skillNames = append(skillNames, key_p)
		}
	}
	bagSizes := []int{8}
	return toret, closestmatch.New(skillNames, bagSizes), nil
}
