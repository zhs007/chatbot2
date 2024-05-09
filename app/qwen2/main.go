package main

import (
	"os"

	"github.com/zhs007/chatbot2/core"
	"github.com/zhs007/goutils"
)

func main() {
	goutils.InitLogger2("qwen", "v1.0.0", "debug", true, "./")

	token := os.Getenv("DASHSCOPE_API_KEY")

	chatbot, err := core.NewChatbot(token, "./cfg")
	if err != nil {
		goutils.Error("NewChatbot",
			goutils.Err(err))

		return
	}

	user := chatbot.MgrUsers.GetUser("zhs007", chatbot.MgrCharacters)
	user.ProcChat(chatbot, "组件有哪些", func(msg *core.Message) {
		goutils.Info(msg.Message)
	})
}
