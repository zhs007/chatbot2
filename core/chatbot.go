package core

import (
	"context"
	"path"

	"github.com/zhs007/dashscopego"
	"github.com/zhs007/dashscopego/qwen"
	"github.com/zhs007/goutils"
)

type Chatbot struct {
	qwenClient    *dashscopego.TongyiClient
	MgrCharacters *CharacterMgr
	MgrUsers      *UserMgr
}

func (bot *Chatbot) SendChat(user *User) (string, string, error) {
	req := &dashscopego.TextRequest{
		Input:  *user.input,
		Plugin: `{"pdf_extracter":{}}`,
	}

	ctx := context.TODO()
	resp, err := bot.qwenClient.CreateCompletion(ctx, req)
	if err != nil {
		goutils.Error("Chatbot.SendChat:CreateCompletion",
			goutils.Err(err))

		return "", "", err
	}

	return resp.Output.Choices[0].Message.Role, resp.Output.Choices[0].Message.Content.ToString(), nil
}

func NewChatbot(apiKey string, cfgPath string) (*Chatbot, error) {
	model := qwen.QwenTurbo
	cli := dashscopego.NewTongyiClient(model, apiKey)

	mgrCharacter, err := LoadCharacterMgr(path.Join(cfgPath, "characters.yaml"))
	if err != nil {
		goutils.Error("NewChatbot:LoadCharacterMgr",
			goutils.Err(err))

		return nil, err
	}

	mgrUser, err := LoadUserMgr(path.Join(cfgPath, "users.yaml"))
	if err != nil {
		goutils.Error("NewChatbot:LoadUserMgr",
			goutils.Err(err))

		return nil, err
	}

	mgrUser.Rebuild(mgrCharacter)

	return &Chatbot{
		qwenClient:    cli,
		MgrCharacters: mgrCharacter,
		MgrUsers:      mgrUser,
	}, nil
}
