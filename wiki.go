package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

var data map[string]interface{}

type wikiparse struct {
	id      string
	title   string
	content string
	link    string
}

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}

func getResponse(url string) []byte {

	res, err := http.Get(url)
	check_error(err)
	rbody, reader_error := ioutil.ReadAll(res.Body)
	check_error(reader_error)

	return rbody

}

func main() {
	var url = "http://baak.gunadarma.ac.id/index.php?stateid=jadkul&substep=search&cari=3ia07&bywhat=kelas"
	result := getResponse(url)

	json.Unmarshal(result, &data)

	fmt.Printf("result : %v \n", data)

}
