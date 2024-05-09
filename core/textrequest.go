package core

import (
	"context"

	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/qwen"
	"github.com/zhs007/goutils"
)

type TextRequest struct {
	BasicRequest
}

func (req *TextRequest) Start(chatbot *Chatbot, msg string) (*Message, error) {
	input := &dashscopego.TextInput{
		Messages: []dashscopego.TextMessage{},
	}

	for _, v := range req.History {
		input.Messages = append(input.Messages, dashscopego.TextMessage{
			Role: v.Role,
			Content: &qwen.TextContent{
				Text: v.Message,
			},
		})
	}

	input.Messages = append(input.Messages, dashscopego.TextMessage{
		Role: "user",
		Content: &qwen.TextContent{
			Text: req.Character.genChat(msg),
		},
	})

	treq := &dashscopego.TextRequest{
		Input: *input,
	}

	resp, err := chatbot.qwenClient.CreateCompletion(context.TODO(), treq)
	if err != nil {
		goutils.Error("TextRequest.Start:CreateCompletion",
			goutils.Err(err))

		return nil, err
	}

	return &Message{
		Role:    resp.Output.Choices[0].Message.Role,
		Message: resp.Output.Choices[0].Message.Content.ToString(),
	}, nil
}
