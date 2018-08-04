package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

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
	Tb_data []JadKul
}

type WkParameters struct {
	Command string
	URI     []string
	Output  string
}

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tp = template.Must(template.ParseFiles("jadwal.tmpl"))

	db, err := connection()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	request, errss := db.Query("select kelas,hari,mata_kuliah,waktu,dosen,ruang from kuliah_pta2016 where kelas='3ia07' ")
	if errss != nil {
		fmt.Println(errss.Error())
		os.Exit(0)
	}

	defer request.Close()

	data := table{}
	jad := JadKul{}
	now := time.Now()
	req := strconv.Itoa(now.Day()) + "-" + now.Month().String() + "-" + strconv.Itoa(now.Year()) + ", " + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()) + " WIB"

	data = table{
		Title:   "3ia07",
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

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}

	}

	if err = request.Err(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	if err = tp.Execute(w, data); err != nil {
		fmt.Println(err)
	}

}

func main() {

	urls := []string{
		"http://localhost:8080/pdf",
	}

	p := WkParameters{
		URI:    urls,
		Output: "test.pdf",
	}

	MakePdf(p)

	router := httprouter.New()
	router.GET("/pdf/:request", index)

	log.Fatal(http.ListenAndServe(":8000", router))

	/*
		http.HandleFunc("/pdf", index)
		http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
		fmt.Println("starting web server at http://localhost:8080/")
		http.ListenAndServe(":8080", nil)*/

}

func MakePdf(p WkParameters) error {
	var command = "../wkhtmltox/bin/./wkhtmltopdf"
	var args []string

	args = append(
		p.URI,
		p.Output,
	)

	run, errs := exec.Command(command, args...).Output()
	fmt.Printf(" -> pwd\n%s\n", string(run))

	if errs != nil {
		fmt.Println(errs)
		return errs
	}
	move_dir(p.Output)

	return errs

}

func move_dir(fname string) error {
	var folder = "pdf_file"
	var sub = "jadwal_kelas"
	var mv = folder + "/" + sub

	os.Mkdir(folder, 0777)
	os.Mkdir(mv, 0777)
	_, err := exec.Command("mv", fname, mv).Output()

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func connection() (*sql.DB, error) {
	var db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/baak")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	return db, nil
}
