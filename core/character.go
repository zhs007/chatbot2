package core

import (
	"fmt"

	"github.com/zhs007/dashscopego"
	"github.com/zhs007/dashscopego/qwen"
)

type Character struct {
	Name   string   `yaml:"-" json:"-"`           // name
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

func (character *Character) GenChatMessage(msg string) dashscopego.TextMessage {
	if len(character.Files) == 0 {
		return dashscopego.TextMessage{
			Role: "user",
			Content: &qwen.TextContent{
				Text: msg,
			},
		}
	}

	str := ""
	for _, v := range character.Files {
		str += fmt.Sprintf(`,{"file": "%v"}`, v)
	}

	txtmsg := dashscopego.TextMessage{
		Role: "user",
		Content: &qwen.TextContent{
			Text:  fmt.Sprintf(`[{"text": "%v"}%v]`, msg, str),
			IsRaw: false,
		},
	}

	return txtmsg
}
