package core

import (
	"log/slog"
	"path"

	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/qwen"
	"github.com/zhs007/goutils"
)

type Chatbot struct {
	qwenClient    *dashscopego.TongyiClient
	MgrCharacters *CharacterMgr
	MgrUsers      *UserMgr
}

// func (bot *Chatbot) sendChat(req *dashscopego.TextRequest) (string, string, error) {
// 	// req := &dashscopego.TextRequest{
// 	// 	Input: *input,
// 	// 	// Plugin: `{"pdf_extracter":{}}`,
// 	// }

// 	ctx := context.TODO()
// 	resp, err := bot.qwenClient.CreateCompletion(ctx, req)
// 	if err != nil {
// 		goutils.Error("Chatbot.sendChat:CreateCompletion",
// 			goutils.Err(err))

// 		return "", "", err
// 	}

// 	return resp.Output.Choices[0].Message.Role, resp.Output.Choices[0].Message.Content.ToString(), nil
// }

// func (bot *Chatbot) sendFileChat(req *dashscopego.FileRequest) (string, string, error) {
// 	// req := &dashscopego.TextRequest{
// 	// 	Input: *input,
// 	// 	// Plugin: `{"pdf_extracter":{}}`,
// 	// }

// 	ctx := context.TODO()
// 	resp, err := bot.qwenClient.CreateFileCompletion(ctx, req)
// 	if err != nil {
// 		goutils.Error("Chatbot.sendChat:CreateFileCompletion",
// 			goutils.Err(err))

// 		return "", "", err
// 	}

// 	return resp.Output.Choices[0].Message.Role, resp.Output.Choices[0].Message.Content.ToString(), nil
// }

// func (bot *Chatbot) SendChat(user *User) (string, string, error) {
// 	return bot.sendChat(user.character.GenRequest(user.input))
// 	// req := &dashscopego.TextRequest{
// 	// 	Input:  *user.input,
// 	// 	Plugin: `{"pdf_extracter":{}}`,
// 	// }

// 	// ctx := context.TODO()
// 	// resp, err := bot.qwenClient.CreateCompletion(ctx, req)
// 	// if err != nil {
// 	// 	goutils.Error("Chatbot.SendChat:CreateCompletion",
// 	// 		goutils.Err(err))

// 	// 	return "", "", err
// 	// }

// 	// return resp.Output.Choices[0].Message.Role, resp.Output.Choices[0].Message.Content.ToString(), nil
// }

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

	goutils.Info("NewChatbot",
		slog.Any("mgrCharacter", mgrCharacter))

	return &Chatbot{
		qwenClient:    cli,
		MgrCharacters: mgrCharacter,
		MgrUsers:      mgrUser,
	}, nil
}
