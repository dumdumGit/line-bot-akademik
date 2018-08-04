package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/line/line-bot-sdk-go/linebot"
)

//var untuk main.go
type TheApp struct {
	bot     *linebot.Client
	BaseUrl string
	DownDir string
}

var cm string = "@botak"
var valid bool
var exists bool

var command = map[string]string{
	"jadwal":              " jadwal",
	"kalender_akademik":   " kalender",
	"berita":              " news",
	"cari_mahasiswa_baru": " maba",
	"kelas2_baru":         " kelas2",
	"UTS":                 " uts",
	"UAS":                 " uas",
	"kelas":               " kelas",
	"stackoverflow":       " so",
	"wikipedia":           " wiki",
	"tugas":               "create task",
	"show":                "show task",
}

var param = map[string]string{
	"show":    "show",
	"buat":    "create",
	"jadwal":  "today",
	"jadwal2": "besok",
}

//var akademik
type Berita struct {
	id       string
	reporter string
	rubrik   string
	tanggal  string
	status   string
	foto     string
	judul    string
	isi      string
	gambar   string
}

type Dosen struct {
	kd       string
	nama     string
	username string
}

type Krs struct {
	id_jadwal5 int
	kelas      string
	mtkuliah   string
	dosen      string
	waktu      string
	ruang      string
	hr         string
	id_jur     int
}

type Kalender struct {
	Id       int
	Kegiatan string
	Periode  string
	Urutan   string
	Waktu    string
}

type Mhs_baru struct {
	id     int
	nopend string
	NPM    string
	Nama   string
	Kelas  string
	pil    string
}

type JadKul struct {
	Kelas       string
	Hari        string
	Mata_kuliah string
	Waktu       string
	Dosen       string
	Ruang       string
}

type table struct {
	Title   string
	Tanggal string
	status  bool
	Tb_data []JadKul
	Kal     []Kalender
	Mhs     []Mhs_baru
}

type WkParameters struct {
	Command string
	URI     []string
	Output  string
}

type task_mdl struct {
	id      int
	judul   string
	isi     string
	tanggal string
}

type dir struct {
	file   string
	folder string
	sub    string
	ex     string
}

func connection() (*sql.DB, error) {
	var db, err = sql.Open("mysql", "user:password@tcp(ip)/database")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	return db, nil
}
