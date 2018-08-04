package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/line/line-bot-sdk-go/linebot"
)

func rowExists(query string, args ...interface{}) bool {
	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	query = fmt.Sprintf("SELECT exists (%s)", query)
	errs := db.QueryRow(query, args...).Scan(&exists)
	if errs != nil && errs != sql.ErrNoRows {
		log.Fatalf("error checking if row exists '%s' %v", args, errs)
	}

	return exists
}

/*
func (app *TheApp) News() []Berita {
	data := Berita{}

	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())

	}
	defer db.Close()

	if rowExists("select tanggal,foto,judul,isi from news") {

		request, errss := db.Query("select tanggal,foto,judul,isi from news")

		if errss != nil && errss != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", errss)
			confirm()
		}

		defer request.Close()

		jad := JadKul{}
		now := time.Now()
		req := strconv.Itoa(now.Day()) + "-" + now.Month().String() + "-" + strconv.Itoa(now.Year()) + ", " + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()) + " WIB"

		data = table{
			Title:   ps.ByName("request"),
			Tanggal: req,
		}

		for request.Next() {

			var err = request.Scan(&jad.Kelas, &jad.Hari, &jad.Mata_kuliah, &jad.Waktu, &jad.Dosen, &jad.Ruang)

			data.Tb_data = append(data.Tb_data, JadKul{
				Kelas:       jad.Kelas,
				Hari:        jad.Hari,
				Mata_kuliah: jad.Mata_kuliah,
				Waktu:       jad.Waktu,
				Dosen:       jad.Dosen,
				Ruang:       jad.Ruang,
			})

			if err != nil && err != sql.ErrNoRows {
				log.Fatalf("error checking if row exists %v", err)
				data.status = false
				confirm()
			}

		}

		if err != nil && err != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", err)
			data.status = false
			confirm()
		}

		return data
	}
	return data

}
*/

func (app *TheApp) Maba(ps httprouter.Params) table {
	data := table{}
	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	if rowExists("select npm,nama,kelas from tabel where nama like '%" + ps.ByName("request") + "%' or npm like '%" + ps.ByName("request") + "%' or kelas like '%" + ps.ByName("request") + "%'") {

		request, err := db.Query("select npm,nama,kelas from mhs_baru2017 where nama like '%" + ps.ByName("request") + "%' or npm like '%" + ps.ByName("request") + "%' or kelas like '%" + ps.ByName("request") + "%'")

		if err != nil && err != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", err)
			confirm()
		}

		defer request.Close()

		jad := Mhs_baru{}
		now := time.Now()
		req := strconv.Itoa(now.Day()) + "-" + now.Month().String() + "-" + strconv.Itoa(now.Year()) + ", " + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()) + " WIB"

		data = table{
			Title:   ps.ByName("request"),
			Tanggal: req,
		}

		for request.Next() {

			var err = request.Scan(&jad.NPM, &jad.Nama, &jad.Kelas)

			data.Mhs = append(data.Mhs, Mhs_baru{
				NPM:   jad.NPM,
				Nama:  jad.Nama,
				Kelas: jad.Kelas,
			})

			if err != nil && err != sql.ErrNoRows {
				log.Fatalf("error checking if row exists %v", err)
				data.status = false
				confirm()
			}

		}

		if err != nil && err != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", err)
			data.status = false
			confirm()
		}

		return data
	}
	return data
}

func (app *TheApp) Kalender() table {
	data := table{}
	kal := Kalender{}
	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())

	}
	defer db.Close()

	request, errss := db.Query("select Kegiatan,Urutan,Waktu from tabel where IDPeriode=(SELECT max(IDPeriode) from kaldamik_b) order by Urutan")

	if errss != nil && errss != sql.ErrNoRows {
		log.Fatalf("error checking if row exists %v", errss)

	}

	defer request.Close()

	for request.Next() {

		var err = request.Scan(&kal.Kegiatan, &kal.Urutan, &kal.Waktu)

		data.Kal = append(data.Kal, Kalender{
			Kegiatan: kal.Kegiatan,
			Urutan:   kal.Urutan,
			Waktu:    kal.Waktu,
		})

		if err != nil && err != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", err)
		}

	}

	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("error checking if row exists %v", err)

	}

	return data

}

func (app *TheApp) Jadwal_dosen(ps httprouter.Params) table {
	data := table{}
	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())

	}
	defer db.Close()

	if rowExists("select kelas,hari,mata_kuliah,waktu,dosen,ruang from tabel where dosen like '%" + ps.ByName("request") + "%' ") {

		request, errss := db.Query("select kelas,hari,mata_kuliah,waktu,dosen,ruang from kuliah_pta2017  where dosen  like '%" + ps.ByName("request") + "%' ")

		if errss != nil && errss != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", errss)
			confirm()
		}

		defer request.Close()

		jad := JadKul{}
		now := time.Now()
		req := strconv.Itoa(now.Day()) + "-" + now.Month().String() + "-" + strconv.Itoa(now.Year()) + ", " + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()) + " WIB"

		data = table{
			Title:   ps.ByName("request"),
			Tanggal: req,
		}

		for request.Next() {

			var err = request.Scan(&jad.Kelas, &jad.Hari, &jad.Mata_kuliah, &jad.Waktu, &jad.Dosen, &jad.Ruang)

			data.Tb_data = append(data.Tb_data, JadKul{
				Kelas:       jad.Kelas,
				Hari:        jad.Hari,
				Mata_kuliah: jad.Mata_kuliah,
				Waktu:       jad.Waktu,
				Dosen:       jad.Dosen,
				Ruang:       jad.Ruang,
			})

			if err != nil && err != sql.ErrNoRows {
				log.Fatalf("error checking if row exists %v", err)
				data.status = false
				confirm()
			}

		}

		if err != nil && err != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", err)
			data.status = false
			confirm()
		}

		return data
	}
	return data
}

func (app *TheApp) Jadwal_kelas(ps httprouter.Params) table {
	data := table{}
	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())

	}
	defer db.Close()

	if rowExists("select kelas,hari,mata_kuliah,waktu,dosen,ruang from tabel where kelas=? ", ps.ByName("request")) {

		request, errss := db.Query("select kelas,hari,mata_kuliah,waktu,dosen,ruang from kuliah_pta2017 where kelas=? ", ps.ByName("request"))

		if errss != nil && errss != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", errss)
			confirm()
		}

		defer request.Close()

		jad := JadKul{}
		now := time.Now()
		req := strconv.Itoa(now.Day()) + "-" + now.Month().String() + "-" + strconv.Itoa(now.Year()) + ", " + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()) + " WIB"

		data = table{
			Title:   ps.ByName("request"),
			Tanggal: req,
		}

		for request.Next() {

			var err = request.Scan(&jad.Kelas, &jad.Hari, &jad.Mata_kuliah, &jad.Waktu, &jad.Dosen, &jad.Ruang)

			data.Tb_data = append(data.Tb_data, JadKul{
				Kelas:       jad.Kelas,
				Hari:        jad.Hari,
				Mata_kuliah: jad.Mata_kuliah,
				Waktu:       jad.Waktu,
				Dosen:       jad.Dosen,
				Ruang:       jad.Ruang,
			})

			if err != nil && err != sql.ErrNoRows {
				log.Fatalf("error checking if row exists %v", err)
				data.status = false
				confirm()
			}

		}

		if err != nil && err != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", err)
			data.status = false
			confirm()
		}

		return data
	}
	return data
}

func (app *TheApp) Jadwal_uts(ps httprouter.Params) table {
	data := table{}
	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())

	}
	defer db.Close()

	if rowExists("select kelas,hari,mata_kuliah,waktu,ruang from tabel where kelas=? ", ps.ByName("request")) {

		request, errss := db.Query("select kelas,hari,mata_kuliah,waktu,ruang from uts_pta2016 where kelas=? ", ps.ByName("request"))

		if errss != nil && errss != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", errss)
			confirm()
		}

		defer request.Close()

		jad := JadKul{}
		now := time.Now()
		req := strconv.Itoa(now.Day()) + "-" + now.Month().String() + "-" + strconv.Itoa(now.Year()) + ", " + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()) + " WIB"

		data = table{
			Title:   ps.ByName("request"),
			Tanggal: req,
		}

		for request.Next() {

			var err = request.Scan(&jad.Kelas, &jad.Hari, &jad.Mata_kuliah, &jad.Waktu, &jad.Ruang)

			data.Tb_data = append(data.Tb_data, JadKul{
				Kelas:       jad.Kelas,
				Hari:        jad.Hari,
				Mata_kuliah: jad.Mata_kuliah,
				Waktu:       jad.Waktu,
				Ruang:       jad.Ruang,
			})

			if err != nil && err != sql.ErrNoRows {
				log.Fatalf("error checking if row exists %v", err)
				data.status = false
				confirm()
			}

		}

		if err != nil && err != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", err)
			data.status = false
			confirm()
		}

		return data
	}
	return data
}

func (app *TheApp) Jadwal_uas(ps httprouter.Params) table {
	data := table{}
	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())

	}
	defer db.Close()

	if rowExists("select kelas,hari,mata_kuliah,waktu,ruang from tabel where kelas=? ", ps.ByName("request")) {

		request, errss := db.Query("select kelas,hari,mata_kuliah,waktu,ruang from tabel where kelas=? ", ps.ByName("request"))

		if errss != nil && errss != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", errss)
			confirm()
		}

		defer request.Close()

		jad := JadKul{}
		now := time.Now()
		req := strconv.Itoa(now.Day()) + "-" + now.Month().String() + "-" + strconv.Itoa(now.Year()) + ", " + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()) + " WIB"

		data = table{
			Title:   ps.ByName("request"),
			Tanggal: req,
		}

		for request.Next() {

			var err = request.Scan(&jad.Kelas, &jad.Hari, &jad.Mata_kuliah, &jad.Waktu, &jad.Ruang)

			data.Tb_data = append(data.Tb_data, JadKul{
				Kelas:       jad.Kelas,
				Hari:        jad.Hari,
				Mata_kuliah: jad.Mata_kuliah,
				Waktu:       jad.Waktu,
				Ruang:       jad.Ruang,
			})

			if err != nil && err != sql.ErrNoRows {
				log.Fatalf("error checking if row exists %v", err)
				data.status = false
				confirm()
			}

		}

		if err != nil && err != sql.ErrNoRows {
			log.Fatalf("error checking if row exists %v", err)
			data.status = false
			confirm()
		}

		return data
	}
	return data
}

func (app *TheApp) task(msg string, token string) {
	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	now := time.Now()
	date := now.Year()

	_, err = db.Exec("insert into tabel values (?, ?, ?, ?)", msg, "", date)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("insert success!")

	message := "Task " + msg + " telah dibuat !"

	app.HandleMessage(token, message)
}

func (app *TheApp) show_task(replyToken string) {
	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	request, err := db.Query("select id, judul, tanggal from tabel ")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)

	}

	defer request.Close()
	var result []task_mdl

	for request.Next() {
		var data = task_mdl{}
		var err = request.Scan(&data.id, &data.judul, &data.tanggal)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}

		result = append(result, data)
		if err = request.Err(); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}

	for _, data := range result {
		template := linebot.NewConfirmTemplate(
			data.judul,
			linebot.NewMessageTemplateAction("Yes", "Yes!"),
			linebot.NewMessageTemplateAction("No", "No!"),
		)
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTemplateMessage("Your Task", template),
		).Do(); err != nil {
			return
		}
	}

}
