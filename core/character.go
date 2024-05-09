package core

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/qwen"
	"github.com/zhs007/goutils"
)

type Character struct {
	Name     string   `yaml:"-" json:"-"`               // name
	Prompt   string   `yaml:"prompt" json:"prompt"`     // prompt
	Type     string   `yaml:"type" json:"type"`         // type
	Files    []string `yaml:"files" json:"files"`       // files
	Workflow []string `yaml:"workflow" json:"workflow"` // workflow
}

func (character *Character) GenRequest(input *dashscopego.TextInput) *dashscopego.TextRequest {
	plugins := character.GenPlugins()
	if plugins == nil {
		return &dashscopego.TextRequest{
			Input: *input,
		}
	}

	return &dashscopego.TextRequest{
		Input:   *input,
		Plugins: *plugins,
	}
}

func (character *Character) GenPlugins() *qwen.Plugins {
	if len(character.Files) > 0 {
		return &qwen.Plugins{qwen.PluginPDFExtracter: {}}
	}

	return nil
}

func (character *Character) IsWorkflow() bool {
	return len(character.Workflow) > 0
}

func (character *Character) ProcWorkflow(chatbot *Chatbot, msg string, onChatbot FuncOnChatbot) (string, string, error) {
	for _, v := range character.Workflow {
		c := chatbot.MgrCharacters.Get(v)
		if c != nil {
			input := c.GenInput()
			input.Messages = append(input.Messages, c.GenChatMessage(msg))

			role, ret, err := chatbot.sendChat(c.GenRequest(input))
			if err != nil {
				goutils.Error("Character.ProcWorkflow:sendChat",
					slog.String("character", v),
					goutils.Err(err))

				return "", "", err
			}

			onChatbot(role, ret)

			msg = ret
		}
	}

	return "assistant", msg, nil
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

func (character *Character) genChat(msg string) string {
	if character.Type == "simple" {
		return fmt.Sprintf("%v\n\n%v", character.Prompt, msg)
	}

	return msg
}

func (character *Character) genFile(fn string) string {
	curdir, _ := os.Getwd()
	return fmt.Sprintf("file://%v", path.Join(curdir, fn))
}

func (character *Character) GenChatMessage(msg string) dashscopego.TextMessage {
	if len(character.Files) == 0 {
		return dashscopego.TextMessage{
			Role: "user",
			Content: &qwen.TextContent{
				Text: character.genChat(msg),
			},
		}
	}

	str := ""
	for _, v := range character.Files {
		str += fmt.Sprintf(`,{"file": "%v"}`, character.genFile(v))
	}

	txtmsg := dashscopego.TextMessage{
		Role: "user",
		Content: &qwen.TextContent{
			Text: fmt.Sprintf(`[{"text": "%v"}%v]`, character.genChat(msg), str),
			// IsRaw: false,
		},
	}

	return txtmsg
}
