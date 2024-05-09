package core

import (
	"context"

	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/qwen"
	"github.com/zhs007/goutils"
)

type FileRequest struct {
	BasicRequest
}

func (req *FileRequest) Start(chatbot *Chatbot, msg string) (*Message, error) {
	input := &dashscopego.FileInput{
		Messages: []dashscopego.FileMessage{},
	}

	for _, v := range req.History {
		fcl := &qwen.FileContentList{}

		fcl.SetText(v.Message)

		input.Messages = append(input.Messages, dashscopego.FileMessage{
			Role:    v.Role,
			Content: fcl,
		})
	}

	fcl := &qwen.FileContentList{}

	fcl.SetText(req.Character.genChat(msg))
	for _, v := range req.Character.Files {
		fcl.SetFile(req.Character.genFile(v))
	}

	input.Messages = append(input.Messages, dashscopego.FileMessage{
		Role:    "user",
		Content: fcl,
	})

	treq := &dashscopego.FileRequest{
		Input: *input,
	}

	resp, err := chatbot.qwenClient.CreateFileCompletion(context.TODO(), treq)
	if err != nil {
		goutils.Error("FileRequest.Start:CreateFileCompletion",
			goutils.Err(err))

		return nil, err
	}

	return &Message{
		Role:    resp.Output.Choices[0].Message.Role,
		Message: resp.Output.Choices[0].Message.Content.ToString(),
	}, nil
}