package telegram

import (
	"fmt"
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhs007/chatbot2/core"
	"github.com/zhs007/goutils"
)

type Serv struct {
	bot     *tgbotapi.BotAPI
	chatbot *core.Chatbot
}

func (serv *Serv) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := serv.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			goutils.Info(fmt.Sprintf("[%s] %s", update.Message.From.UserName, update.Message.Text))

			uid := strings.TrimSpace(update.Message.From.UserName)
			user := serv.chatbot.MgrUsers.GetUser(uid, serv.chatbot.MgrCharacters)
			if user == nil {
				goutils.Error("Serv.Start:GetUser",
					slog.String("uid", uid),
					goutils.Err(core.ErrGetUser))

				return core.ErrGetUser
			}

			goutils.Info("user character",
				slog.String("character", user.CharacterName))

			txt := strings.TrimSpace(update.Message.Text)
			character := serv.chatbot.MgrCharacters.Get(txt)
			if character != nil {
				user.SetCharacter(character)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "OK, I got it.")
				msg.ReplyToMessageID = update.Message.MessageID

				serv.bot.Send(msg)
			} else {
				// user.AddChat(txt)

				role, msg, err := user.ProcChat(serv.chatbot, txt, func(role string, ret string) {
					msgRet := tgbotapi.NewMessage(update.Message.Chat.ID, ret)
					msgRet.ReplyToMessageID = update.Message.MessageID

					serv.bot.Send(msgRet)
				})
				if err != nil {
					goutils.Error("Serv.Start:ProcChat",
						goutils.Err(err))

					return err
				}

				user.AddReply(role, msg)

				// msgRet := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
				// msgRet.ReplyToMessageID = update.Message.MessageID

				// serv.bot.Send(msgRet)
			}

			// input := dashscopego.TextInput{
			// 	Messages: []dashscopego.TextMessage{
			// 		{Role: "system", Content: &qwen.TextContent{
			// 			Text: `你是SlotCraft的智能助手，叫SlotCraft AI，请用英文回答问题`,
			// 		}},
			// 	},
			// }

			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID

			// serv.bot.Send(msg)
		}
	}

	return nil
}

func NewServ(tgToken string, apiKey string, cfgPath string) (*Serv, error) {
	chatbot, err := core.NewChatbot(apiKey, cfgPath)
	if err != nil {
		goutils.Error("StartServ:NewChatbot",
			goutils.Err(err))

		return nil, err
	}

	// mgrCharacters, err := core.LoadCharacterMgr(path.Join(cfgPath, "characters.yaml"))
	// if err != nil {
	// 	goutils.Error("StartServ:LoadCharacterMgr",
	// 		goutils.Err(err))

	// 	return nil, err
	// }

	// mgrUsers, err := core.LoadUserMgr(path.Join(cfgPath, "users.yaml"))
	// if err != nil {
	// 	goutils.Error("StartServ:LoadUserMgr",
	// 		goutils.Err(err))

	// 	return nil, err
	// }

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		goutils.Error("StartServ:NewBotAPI",
			goutils.Err(err))

		return nil, err
	}

	bot.Debug = true

	goutils.Info(fmt.Sprintf("Authorized on account %v", bot.Self.UserName))

	return &Serv{
		bot:     bot,
		chatbot: chatbot,
	}, nil
}
