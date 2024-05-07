package core

import (
	"github.com/zhs007/dashscopego"
	"github.com/zhs007/dashscopego/qwen"
	"github.com/zhs007/goutils"
)

type Chatbot struct {
	qwenClient    *dashscopego.TongyiClient
	mgrCharacters *CharacterMgr
}

func NewChatbot(apiKey string) (*Chatbot, error) {
	model := qwen.QwenTurbo
	cli := dashscopego.NewTongyiClient(model, apiKey)

	mgrCharacter, err := LoadCharacterMgr("./cfg/characters.yaml")
	if err != nil {
		goutils.Error("NewChatbot:LoadCharacterMgr",
			goutils.Err(err))

		return nil, err
	}

	return &Chatbot{
		qwenClient:    cli,
		mgrCharacters: mgrCharacter,
	}, nil
}
