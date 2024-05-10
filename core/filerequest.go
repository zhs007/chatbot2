package core

import (
	"context"

	"github.com/zhs007/dashscopego"
	"github.com/zhs007/dashscopego/qwen"
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

	fcl := &qwen.FileContentList{
		{
			Text: req.Character.genChat(msg),
		},
		{
			// File: "https://qianwen-res.oss-cn-beijing.aliyuncs.com/QWEN_TECHNICAL_REPORT.pdf",
			File: req.Character.genFile(req.Character.Files[0]),
		},
	}

	// fcl.SetText(req.Character.genChat(msg))
	// for _, v := range req.Character.Files {
	// 	fcl.SetFile(req.Character.genFile(v))
	// }

	input.Messages = append(input.Messages, dashscopego.FileMessage{
		Role:    "user",
		Content: fcl,
	})

	treq := &dashscopego.FileRequest{
		Input:   *input,
		Plugins: *req.Character.GenPlugins(),
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
