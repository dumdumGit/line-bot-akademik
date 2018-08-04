package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/line/line-bot-sdk-go/linebot"
	"html/template"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	router := httprouter.New()

	app, err := BotApp(
		os.Getenv("CHANNELSECRET"),
		os.Getenv("ACCESSTOKEN"),
		os.Getenv("BASEURL"),
	)
	if err != nil {
		log.Fatal(err)
	}

	router.GET("/jadwal_kelas/:request", app.kelas)
	router.GET("/Jaduts/:request", app.uts)
	router.GET("/Jaduas/:request", app.uas)
	router.GET("/jadwal_dosen/:request", app.dosen)
	router.GET("/Kalender_akademik", app.kaldamik)
	router.GET("/mahasiswa/:request", app.mahasiswa)
	router.POST("/callback", app.callbackHandler)

	router.ServeFiles("/static/*filepath", http.Dir("assets"))
	download := http.FileServer(http.Dir(app.DownDir))
	http.HandleFunc("/download/", http.StripPrefix("/download/", download).ServeHTTP)

	log.Fatal(http.ListenAndServe(":8001", router))

}

func BotApp(ChannelSecret, ChannelToken, BaseURL string) (*TheApp, error) {
	apiEndpointBase := os.Getenv("ENDPOINT_BASE")
	if apiEndpointBase == "" {
		apiEndpointBase = linebot.APIEndpointBase
	}
	bot, err := linebot.New(
		ChannelSecret,
		ChannelToken,
		linebot.WithEndpointBase(apiEndpointBase),
	)
	if err != nil {
		return nil, err
	}
	downloadDir := filepath.Join(filepath.Dir(os.Args[0]), "line-bot")
	_, err = os.Stat(downloadDir)
	if err != nil {
		if err := os.Mkdir(downloadDir, 0777); err != nil {
			return nil, err
		}
	}
	return &TheApp{
		bot:     bot,
		BaseUrl: BaseURL,
		DownDir: downloadDir,
	}, nil
}

func (app *TheApp) kelas(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tp = template.Must(template.ParseFiles("public/template/jadwal.tmpl"))
	data := app.Jadwal_kelas(ps)

	if err := tp.Execute(w, data); err != nil {
		fmt.Println(err)

	}

}

func (app *TheApp) mahasiswa(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tp = template.Must(template.ParseFiles("public/template/mhs.tmpl"))
	data := app.Maba(ps)

	if err := tp.Execute(w, data); err != nil {
		fmt.Println(err)

	}

}

func (app *TheApp) uts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tp = template.Must(template.ParseFiles("public/template/jadwal_ujian.tmpl"))
	data := app.Jadwal_uts(ps)

	if err := tp.Execute(w, data); err != nil {
		fmt.Println(err)

	}

}

func (app *TheApp) uas(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tp = template.Must(template.ParseFiles("public/template/jadwal_uas.tmpl"))
	data := app.Jadwal_uas(ps)

	if err := tp.Execute(w, data); err != nil {
		fmt.Println(err)

	}

}

func (app *TheApp) dosen(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tp = template.Must(template.ParseFiles("public/template/jadwal_dosen.tmpl"))
	data := app.Jadwal_dosen(ps)

	if err := tp.Execute(w, data); err != nil {
		fmt.Println(err)

	}

}

func (app *TheApp) kaldamik(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var tp = template.Must(template.ParseFiles("public/template/kalender.tmpl"))
	data := app.Kalender()

	if err := tp.Execute(w, data); err != nil {
		fmt.Println(err)

	}
}

func (app *TheApp) callbackHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	events, err := app.bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	log.Printf("Events : %v", events)
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				app.TextMessage(message, event.ReplyToken, event.Source)

			default:
				log.Printf("Unknown message: %v", message)
			}
		case linebot.EventTypeFollow:
			app.Help(event.ReplyToken, event.Source.UserID)
		default:
			log.Printf("Unknown event: %v", event)
		}
	}
}

func (app *TheApp) TextMessage(message *linebot.TextMessage, token string, source *linebot.EventSource) {
	var txtmsg string
	var regextext []string

	txtmsg = strings.ToLower(message.Text)

	switch strings.Contains(txtmsg, "@bot") {
	case strings.Contains(txtmsg, command["jadwal"]):
		switch strings.Contains(txtmsg, command["jadwal"]) {
		case strings.Contains(txtmsg, command["kelas"]):
			var SendMessage map[string]string
			SendMessage = map[string]string{}

			regex, _ := regexp.Compile(`[0-9][a-z0-9]+`)
			regextext = regex.FindAllString(txtmsg, 1)
			var key string = strings.Join(regextext, "")

			SendMessage["key"] = key

			urls := []string{
				app.BaseUrl + "/jadwal_kelas/" + key,
			}

			dire := dir{folder: "PDF/", sub: "Jadwal_Kelas/"}

			p := WkParameters{
				Command: "wkhtmltox/bin/./wkhtmltopdf",
				URI:     urls,
				Output:  key + ".pdf",
			}

			if len(regextext) != 0 {

				app.WkHTMLtoPDF(p, dire, token, SendMessage)

				p := WkParameters{Command: "wkhtmltox/bin/./wkhtmltoimage", URI: urls, Output: key + ".jpg"}
				direc := dir{folder: "Original_Image/", sub: "Jadwal_Kelas/"}
				app.WkHTMLtoImage(p, direc, token, SendMessage)

			} else {
				app.bot.ReplyMessage(token, linebot.NewTextMessage("Maaf kak,aku ga dapet jadwalnya, coba cek lagi inputan kelas yang kaka ketik")).Do()
			}

		case strings.Contains(txtmsg, command["UTS"]):
			var SendMessage map[string]string
			SendMessage = map[string]string{}

			regex, _ := regexp.Compile(`[0-9][a-z0-9]+`)
			regextext = regex.FindAllString(txtmsg, 1)
			var key string = strings.Join(regextext, "")

			SendMessage["key"] = key

			urls := []string{
				app.BaseUrl + "/Jaduts/" + key,
			}

			dire := dir{folder: "PDF/", sub: "Jadwal_UTS/"}

			p := WkParameters{
				Command: "wkhtmltox/bin/./wkhtmltopdf",
				URI:     urls,
				Output:  key + ".pdf",
			}

			if len(regextext) != 0 {
				app.WkHTMLtoPDF(p, dire, token, SendMessage)

				p := WkParameters{Command: "wkhtmltox/bin/./wkhtmltoimage", URI: urls, Output: key + ".jpg"}
				direc := dir{folder: "Original_Image/", sub: "Jadwal_UTS/"}
				app.WkHTMLtoImage(p, direc, token, SendMessage)

			} else {
				app.bot.ReplyMessage(token, linebot.NewTextMessage("Maaf kak,aku ga dapet jadwalnya, coba cek lagi inputan kelas yang kaka ketik")).Do()
			}

		case strings.Contains(txtmsg, command["UAS"]):

			var SendMessage map[string]string
			SendMessage = map[string]string{}

			regex, _ := regexp.Compile(`[0-9][a-z0-9]+`)
			regextext = regex.FindAllString(txtmsg, 1)
			var key string = strings.Join(regextext, "")

			SendMessage["key"] = key

			urls := []string{
				app.BaseUrl + "/Jaduas/" + key,
			}

			dire := dir{folder: "PDF/", sub: "Jadwal_UAS/"}

			p := WkParameters{
				Command: "wkhtmltox/bin/./wkhtmltopdf",
				URI:     urls,
				Output:  key + ".pdf",
			}

			if len(regextext) != 0 {
				app.WkHTMLtoPDF(p, dire, token, SendMessage)

				p := WkParameters{Command: "wkhtmltox/bin/./wkhtmltoimage", URI: urls, Output: key + ".jpg"}
				direc := dir{folder: "Original_Image/", sub: "Jadwal_UAS/"}
				app.WkHTMLtoImage(p, direc, token, SendMessage)

			} else {
				app.bot.ReplyMessage(token, linebot.NewTextMessage("Maaf kak,aku ga dapet jadwalnya, coba cek lagi inputan kelas yang kaka ketik")).Do()
			}

		case strings.Contains(txtmsg, "dosen"):
			var SendMessage map[string]string
			SendMessage = map[string]string{}

			regex, _ := regexp.Compile(`[A-Z]+`)
			regextext = regex.FindAllString(message.Text, -1)
			var key string = strings.Join(regextext, " ")
			var keys string = strings.Join(regextext, "_")

			SendMessage["key"] = key

			urls := []string{
				app.BaseUrl + "/jadwal_dosen/" + key,
			}

			dire := dir{folder: "PDF/", sub: "Jadwal_Dosen/"}

			p := WkParameters{
				Command: "wkhtmltox/bin/./wkhtmltopdf",
				URI:     urls,
				Output:  keys + ".pdf",
			}

			if len(regextext) != 0 {

				app.WkHTMLtoPDF(p, dire, token, SendMessage)

				p := WkParameters{Command: "wkhtmltox/bin/./wkhtmltoimage", URI: urls, Output: keys + ".jpg"}
				direc := dir{folder: "Original_Image/", sub: "Jadwal_Dosen/"}
				app.WkHTMLtoImage(p, direc, token, SendMessage)

			} else {
				app.bot.ReplyMessage(token, linebot.NewTextMessage("Maaf kak,aku ga dapet jadwalnya, coba cek lagi inputan kelas yang kaka ketik")).Do()
			}

		case strings.Contains(txtmsg, "uu"):
			app.bot.ReplyMessage(token, linebot.NewTextMessage("Maaf kak, Ujian Utama belum tersedia, karena belum ada jadwalnya")).Do()

		case strings.Contains(txtmsg, "krs"):
			app.bot.ReplyMessage(token, linebot.NewTextMessage("Maaf kak, Jadwal KRS belum tersedia, karena belum ada jadwalnya")).Do()

		default:
			app.bot.ReplyMessage(token, linebot.NewTextMessage("Maaf kak, Aku ga ngerti yang kamu minta"),
				linebot.NewTextMessage("ketik @bot /help untuk melihat daftar command")).Do()
		}

	case strings.Contains(txtmsg, command["kalender_akademik"]):
		exists = true
		var SendMessage map[string]string
		SendMessage = map[string]string{}

		SendMessage["key"] = "akademik"

		urls := []string{
			app.BaseUrl + "/Kalender_akademik/",
		}

		dire := dir{folder: "Kalender/", sub: "Akademik/"}

		pw := WkParameters{
			Command: "wkhtmltox/bin/./wkhtmltopdf",
			URI:     urls,
			Output:  "kalender.pdf",
		}

		app.WkHTMLtoPDF(pw, dire, token, SendMessage)

		p := WkParameters{Command: "wkhtmltox/bin/./wkhtmltoimage", URI: urls, Output: "kalender.jpg"}
		direc := dir{folder: "Original_Image/", sub: "Kalender/"}
		app.WkHTMLtoImage(p, direc, token, SendMessage)

	case strings.Contains(txtmsg, command["berita"]):

	case strings.Contains(txtmsg, command["cari_mahasiswa_baru"]):
		var SendMessage map[string]string
		SendMessage = map[string]string{}

		regex, _ := regexp.Compile(`[0-9][A-Z0-9]+`)
		regextext = regex.FindAllString(message.Text, 1)
		var key string = strings.Join(regextext, " ")
		var keys string = strings.Join(regextext, "_")

		SendMessage["key"] = key

		urls := []string{
			app.BaseUrl + "/mahasiswa/" + key,
		}

		dire := dir{folder: "PDF/", sub: "Maba/"}

		p := WkParameters{
			Command: "wkhtmltox/bin/./wkhtmltopdf",
			URI:     urls,
			Output:  keys + ".pdf",
		}

		if len(regextext) != 0 {

			app.WkHTMLtoPDF(p, dire, token, SendMessage)

			p := WkParameters{Command: "wkhtmltox/bin/./wkhtmltoimage", URI: urls, Output: keys + ".jpg"}
			direc := dir{folder: "Original_Image/", sub: "Maba/"}
			app.WkHTMLtoImage(p, direc, token, SendMessage)

		} else {
			app.bot.ReplyMessage(token, linebot.NewTextMessage("Maaf kak,aku ga dapet jadwalnya, coba cek lagi inputan kelas yang kaka ketik")).Do()
		}

	case strings.Contains(txtmsg, "help"):
		var msg string
		msg = "--Bantuan Bot Akademik--\n\n" + "Request Jadwal Kuliah :\n" + "@bot <jadwal> <kelas> <kelas-yang-diminta>\n Contoh: @bot aku mau minta jadwal untuk kelas 3ia07\n\n"
		msg = msg + "Request Jadwal UAS\n" + "@bot <jadwal> <UAS> <kelas-yang-diminta>\n Contoh: @bot aku mau minta jadwal UAS 3ia07\n\n"
		msg = msg + "Request Jadwal UTS" + "@bot <jadwal> <UTS> <kelas-yang-diminta>\n Contoh: @bot aku mau minta jadwal UTS 3ia07\n\n"
		msg = msg + "Request Jadwal dosen" + "@bot <jadwal> <dosen> <NAMA-KAPITAL>\n Contoh: @bot aku mau minta jadwal dosen SADAN DWI\n\n"
		msg = msg + "Request Kalender Akademik" + "@bot kalender \n Contoh: @bot aku mau minta kalender akademik\n\n"

		app.bot.ReplyMessage(token, linebot.NewTextMessage(msg)).Do()

	case txtmsg == "":
		app.Help(token, source.UserID)

	default:
		app.bot.ReplyMessage(token, linebot.NewTextMessage("Maaf kak, Aku ga ngerti yang kamu minta"),
			linebot.NewTextMessage("ketik @bot /help untuk melihat daftar command")).Do()

	}
}

func confirm() {

}

func (app *TheApp) WkHTMLtoPDF(p WkParameters, dire dir, token string, SendMessage map[string]string) {
	var args []string
	SendMessage["pdf"] = "Aku juga punya dalam bentuk PDF kalau diperlukan \n"
	SendMessage["pdf"] = SendMessage["pdf"] + app.BaseUrl + "/static/" + dire.folder + dire.sub + p.Output

	args = append(
		p.URI,
		p.Output,
	)

	_, errs := exec.Command(p.Command, args...).Output()

	if errs != nil {
		fmt.Println(errs)
	}
	app.move_dir(dire, p, token)

}

func (app *TheApp) WkHTMLtoImage(p WkParameters, dire dir, token string, SendMessage map[string]string) {
	var args []string

	args = append(
		p.URI,
		p.Output,
	)

	_, errs := exec.Command(p.Command, args...).Output()

	if errs != nil {
		fmt.Println(errs)
	}
	app.move_dir(dire, p, token)
	app.create_prev(p, dire, token, SendMessage)

}

func (app *TheApp) move_dir(dire dir, p WkParameters, token string) {
	var root = "assets/"

	os.Mkdir(root, 0777)
	os.Mkdir(root+dire.folder, 0777)
	os.Mkdir(root+dire.folder+dire.sub, 0777)
	_, err := exec.Command("mv", p.Output, root+dire.folder+dire.sub).Output()

	if err != nil {
		fmt.Println(err)
	}

}

func (app *TheApp) create_prev(p WkParameters, dire dir, token string, SendMessage map[string]string) {
	var root = "assets/"
	var prev = "Preview/"

	os.Mkdir(root+prev, 0777)
	os.Mkdir(root+prev+dire.sub, 0777)
	src, err := imaging.Open(root + dire.folder + dire.sub + p.Output)
	if err != nil {
		log.Fatalf("Open failed: %v", err)
	}

	err = imaging.Save(src, root+prev+dire.sub+"PREVIEW_"+p.Output)

	if err != nil {
		log.Fatalf("Save failed: %v", err)
	}

	app.preview(dire, p, token, SendMessage)
}

func (app *TheApp) preview(dire dir, p WkParameters, token string, SendMessage map[string]string) {
	var root = "assets/"
	var prev = "Preview/"

	src, err := imaging.Open(root + prev + dire.sub + "PREVIEW_" + p.Output)
	if err != nil {
		log.Fatalf("Open failed: %v", err)
	}

	resize := imaging.Resize(src, 240, 240, imaging.Lanczos)

	fmt.Println("SUCCESS, cImg dimension for preview image :", resize.Bounds())

	if _, err := os.Stat(root + prev + dire.sub + "PREVIEW_" + p.Output); os.IsNotExist(err) {
		fmt.Println("image is Exist")
	} else {
		os.Remove(root + prev + dire.sub + "PREVIEW_" + p.Output)
	}

	err = imaging.Save(resize, root+prev+dire.sub+"PREVIEW_"+p.Output)

	if err != nil {
		log.Fatalf("Save failed: %v", err)
	}

	if exists {

		SendMessage["image"] = dire.sub + " untuk " + SendMessage["key"] + "\n Aku dapet nih kak datanya"
		app.HandleImageMessage(dire, p, token, SendMessage)
	} else {
		app.bot.ReplyMessage(token, linebot.NewTextMessage("Maaf kak,aku ga dapet datanya"),
			linebot.NewTextMessage("coba cek lagi inputan anda masukkan ya!"),
			linebot.NewTextMessage("Data yang anda cari adalah "+SendMessage["key"])).Do()

	}
}
