package main

import (
	"context"
	"fmt"
	"os"

	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/qwen"
)

func main() {
	model := qwen.QwenTurbo
	token := os.Getenv("DASHSCOPE_API_KEY")

	if token == "" {
		panic("token is empty")
	}

	cli := dashscopego.NewTongyiClient(model, token)

	input := dashscopego.TextInput{
		Messages: []dashscopego.TextMessage{
			{Role: "system", Content: &qwen.TextContent{Text: "你是一个翻译助手，将我说的话原封不动的翻译成英文，然后再把你翻译的英文翻译回中文，然后分别用 英文是 ... 中文是 ... 的句式来回答我。"}},
		},
	}

	for {
		fmt.Print("Please input text with ENTER\r\n")

		str := ""
		fmt.Scanln(&str)

		input.Messages = append(input.Messages, dashscopego.TextMessage{
			Role:    "user",
			Content: &qwen.TextContent{Text: str},
		})

		// (可选 SSE开启) 需要流式输出时 通过该 Callback Function 获取实时显示的结果
		// 开启 SSE 时的 request_id/finish_reason/token usage 等信息在调用完成统一返回(resp)
		streamCallbackFn := func(ctx context.Context, chunk []byte) error {
			// fmt.Print(string(chunk))
			return nil
		}
		req := &dashscopego.TextRequest{
			Input:       input,
			StreamingFn: streamCallbackFn,
		}

		ctx := context.TODO()
		resp, err := cli.CreateCompletion(ctx, req)
		if err != nil {
			fmt.Print(err)

			break
		}

		fmt.Println(resp.Output.Choices[0].Message.Role)
		fmt.Println(resp.Output.Choices[0].Message.Content.ToString())

		input.Messages = append(input.Messages, dashscopego.TextMessage{
			Role:    resp.Output.Choices[0].Message.Role,
			Content: resp.Output.Choices[0].Message.Content,
		})

		fmt.Println(resp.RequestID)
		fmt.Println(resp.Output.Choices[0].FinishReason)
		fmt.Println(resp.Usage.TotalTokens)
	}
}
