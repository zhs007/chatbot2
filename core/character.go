package core

import (
	"github.com/zhs007/dashscopego"
	"github.com/zhs007/dashscopego/qwen"
)

type Character struct {
	Name   string   `yaml:"name" json:"name"`     // name
	Prompt string   `yaml:"prompt" json:"prompt"` // prompt
	Files  []string `yaml:"files" json:"files"`   // files
}

func (character *Character) GenInput() *dashscopego.TextInput {
	return &dashscopego.TextInput{
		Messages: []dashscopego.TextMessage{
			{Role: "system", Content: &qwen.TextContent{
				Text: character.Prompt,
			}},
		},
	}
}
