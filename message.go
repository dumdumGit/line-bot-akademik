package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
)

func (app *TheApp) HandleMessage(token string, message string) {
	app.bot.ReplyMessage(token, linebot.NewTextMessage(message)).Do()
}

func (app *TheApp) HandleImageMessage(dire dir, p WkParameters, ReplyToken string, message map[string]string) {

	OriginalContent := app.BaseUrl + "/static/" + dire.folder + dire.sub + p.Output
	PreviewContent := app.BaseUrl + "/static/" + "Preview/" + dire.sub + "PREVIEW_" + p.Output

	if _, err := app.bot.ReplyMessage(ReplyToken, linebot.NewTextMessage(message["image"]),
		linebot.NewImageMessage(OriginalContent, PreviewContent),
		linebot.NewTextMessage(message["pdf"])).Do(); err != nil {

		log.Print(err)
	}
}

func (app *TheApp) Help(token, source string) {
	if _, err := app.bot.PushMessage(source, linebot.NewTextMessage("Assalamu'alaikum :) "),
		linebot.NewTextMessage("Hi kak, Kenalin aku Bot Akademik siap membantu "),
		linebot.NewTextMessage("Ketik @bot help untuk menampilkan bantuan :) ")).Do(); err != nil {
		log.Println(source, err)
	}
}

/*
func (app *TheApp) HandleTemplateCarousel(ReplyToken string) {
	imageURL := app.appBaseURL + "/static/buttons/1040.jpg"
	template := linebot.NewCarouselTemplate(
		linebot.NewCarouselColumn(
			imageURL, "hoge", "fuga",
			linebot.NewURITemplateAction("Go to line.me", "https://line.me"),
			linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),
		),
		linebot.NewCarouselColumn(
			imageURL, "hoge", "fuga",
			linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは"),
			linebot.NewMessageTemplateAction("Say message", "Rice=米"),
		),
	)
	if _, err := app.bot.ReplyMessage(
		replyToken,
		linebot.NewTemplateMessage("Carousel alt text", template),
	).Do(); err != nil {
		return err
	}
}
*/
