package core

import (
	"github.com/zhs007/dashscopego"
	"github.com/zhs007/dashscopego/qwen"
)

type Message struct {
	Role    string `yaml:"role" json:"role"`       // role
	Message string `yaml:"message" json:"message"` // message
}

type User struct {
	UserID        string                 `yaml:"userID" json:"userID"`               // userID
	History       []*Message             `yaml:"history" json:"history"`             // history
	CharacterName string                 `yaml:"characterName" json:"characterName"` // characterName
	input         *dashscopego.TextInput `yaml:"-" json:"-"`                         // -
	character     *Character             `yaml:"-" json:"-"`                         // -
}

func (user *User) SetCharacter(character *Character) {
	if user.CharacterName != character.Name {
		user.History = nil

		user.CharacterName = character.Name
	}

	user.character = character
	user.rebuild(character)
}

func (user *User) rebuild(character *Character) {
	if character != nil {
		user.input = character.GenInput()
	} else {
		user.input = &dashscopego.TextInput{
			Messages: []dashscopego.TextMessage{
				{Role: "system", Content: &qwen.TextContent{
					Text: "你是Zerro的AI机器人。",
				}},
			},
		}
	}

	for _, v := range user.History {
		user.input.Messages = append(user.input.Messages, dashscopego.TextMessage{
			Role:    v.Role,
			Content: &qwen.TextContent{Text: v.Message},
		})
	}
}

func (user *User) AddChat(msg string) {
	user.input.Messages = append(user.input.Messages, user.character.GenChatMessage(msg))

	user.History = append(user.History, &Message{
		Role:    "user",
		Message: msg,
	})
}

func (user *User) AddReply(role string, msg string) {
	user.input.Messages = append(user.input.Messages, dashscopego.TextMessage{
		Role:    role,
		Content: &qwen.TextContent{Text: msg},
	})

	user.History = append(user.History, &Message{
		Role:    role,
		Message: msg,
	})
}

func (user *User) Rebuild(mgrCharacters *CharacterMgr) {
	character := mgrCharacters.Get(user.CharacterName)
	user.rebuild(character)
}

func NewUser(uid string) *User {
	return &User{
		UserID: uid,
	}
}
