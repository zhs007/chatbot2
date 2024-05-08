package main

import (
	"os"

	"github.com/zhs007/chatbot2/telegram"
	"github.com/zhs007/goutils"
)

func main() {
	goutils.InitLogger2("telegrambot", "v1.0.0", "debug", true, "./")

	tgtoken := os.Getenv("TGTOKEN")
	qwtoken := os.Getenv("DASHSCOPE_API_KEY")

	serv, err := telegram.NewServ(tgtoken, qwtoken, "./cfg")
	if err != nil {
		goutils.Error("NewServ",
			goutils.Err(err))

		return
	}

	err = serv.Start()
	if err != nil {
		goutils.Error("Start",
			goutils.Err(err))

		return
	}
}
