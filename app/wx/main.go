package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/zhs007/dashscopego"
	"github.com/zhs007/dashscopego/qwen"
	"github.com/zhs007/goutils"
)

type Serv struct {
	cfg        *offConfig.Config
	wc         *wechat.Wechat
	oa         *officialaccount.OfficialAccount
	qwenClient *dashscopego.TongyiClient
}

func (serv *Serv) start(listen string) error {
	r := gin.Default()

	r.Any("/api/v1/serve", serv.OnMsg)

	return r.Run(listen)
}

func (serv *Serv) OnMsg(c *gin.Context) {
	// 传入request和responseWriter
	server := serv.oa.GetServer(c.Request, c.Writer)
	server.SkipValidate(true)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//TODO
		//回复消息：演示回复用户发送的消息
		goutils.Info("got msg",
			slog.String("text", msg.Content))

		input := dashscopego.TextInput{
			Messages: []dashscopego.TextMessage{
				{Role: "system", Content: &qwen.TextContent{
					Text: `你是SlotCraft的智能助手，叫SlotCraft AI，请用英文回答问题`,
				}},
			},
		}

		content := qwen.TextContent{
			Text:  fmt.Sprintf(`[{"text": "%v"},{"file": "https://chatbot2.oss-cn-beijing.aliyuncs.com/slotcraft.pdf"}]`, msg.Content),
			IsRaw: true,
		}

		input.Messages = append(input.Messages, dashscopego.TextMessage{
			Role:    "user",
			Content: &content,
		})

		req := &dashscopego.TextRequest{
			Input: input,
		}

		ctx := context.TODO()
		resp, err := serv.qwenClient.CreateCompletion(ctx, req)
		if err != nil {
			fmt.Print(err)

			return nil
		}

		text := message.NewText(resp.Output.Choices[0].Message.Content.ToString())

		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}

		//article1 := message.NewArticle("测试图文1", "图文描述", "", "")
		//articles := []*message.Article{article1}
		//news := message.NewNews(articles)
		//return &message.Reply{MsgType: message.MsgTypeNews, MsgData: news}

		//voice := message.NewVoice(mediaID)
		//return &message.Reply{MsgType: message.MsgTypeVoice, MsgData: voice}

		//
		//video := message.NewVideo(mediaID, "标题", "描述")
		//return &message.Reply{MsgType: message.MsgTypeVideo, MsgData: video}

		//music := message.NewMusic("标题", "描述", "音乐链接", "HQMusicUrl", "缩略图的媒体id")
		//return &message.Reply{MsgType: message.MsgTypeMusic, MsgData: music}

		//多客服消息转发
		//transferCustomer := message.NewTransferCustomer("")
		//return &message.Reply{MsgType: message.MsgTypeTransfer, MsgData: transferCustomer}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		goutils.Error("Serve Error",
			goutils.Err(err))

		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		goutils.Error("Send Error",
			goutils.Err(err))

		return
	}
}

func newServ(appid string, appSecret string, token string, aes string, apiKey string) *Serv {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()

	cfg := &offConfig.Config{
		AppID:          appid,
		AppSecret:      appSecret,
		Token:          token,
		EncodingAESKey: aes,
		Cache:          memory,
	}
	oa := wc.GetOfficialAccount(cfg)

	model := qwen.QwenTurbo
	cli := dashscopego.NewTongyiClient(model, apiKey)

	return &Serv{
		wc:         wc,
		oa:         oa,
		cfg:        cfg,
		qwenClient: cli,
	}
}

func main() {
	goutils.InitLogger2("wx", "v1.0.0", "debug", true, "./")

	appid := os.Getenv("APPID")
	appSecret := os.Getenv("APPSECRET")
	token := os.Getenv("TOKEN")
	listen := os.Getenv("LISTEN")
	aes := os.Getenv("AES")
	apiKey := os.Getenv("DASHSCOPE_API_KEY")

	serv := newServ(appid, appSecret, token, aes, apiKey)

	serv.start(listen)
}
