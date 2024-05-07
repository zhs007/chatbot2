package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhs007/chatbot2/core"
	"github.com/zhs007/dashscopego"
	"github.com/zhs007/dashscopego/qwen"
	"github.com/zhs007/goutils"
)

func startServ(tgtoken string, apiKey string) error {
	chatbot, err := core.NewChatbot(apiKey)
	if err != nil {
		goutils.Error("startServ:NewChatbot",
			goutils.Err(err))

		return err
	}

	bot, err := tgbotapi.NewBotAPI(tgtoken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			input := dashscopego.TextInput{
				Messages: []dashscopego.TextMessage{
					{Role: "system", Content: &qwen.TextContent{
						Text: `你是SlotCraft的智能助手，叫SlotCraft AI，请用英文回答问题`,
					}},
				},
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}

	return nil
}

func main() {
	tgtoken := os.Getenv("TGTOKEN")
	qwtoken := os.Getenv("DASHSCOPE_API_KEY")

	startServ(tgtoken, qwtoken)
}
