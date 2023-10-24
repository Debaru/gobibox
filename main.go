package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"time"

	"github.com/Debaru/gobibox/internal/gobibox"
	"github.com/Debaru/gobibox/internal/qbittorent"
)

var (
	url            string
	login          string
	password       string
	port           int
	dataRepository string
)

func init() {
	dataRepository = "data/"
}

func main() {
	var err error
	flag.StringVar(&url, "url", "", "URL")
	flag.StringVar(&login, "login", "", "Login")
	flag.StringVar(&password, "password", "", "Password")
	flag.IntVar(&port, "port", 21, "FTP port")
	flag.Parse()

	if _, err := os.Stat(dataRepository); os.IsNotExist(err) {
		err = os.Mkdir(dataRepository, 0740)
		if err != nil {
			panic(err)
		}
	}

	for true {
		err = qbittorent.Connect(url, login, password)
		if err != nil {
			panic(err)
		}

		torrents, err := qbittorent.TorrentCompleted()
		if err != nil {
			panic(err)
		}

		files := gobibox.GetFilesToDownload(torrents)
		data, _ := json.MarshalIndent(&files, "  ", "  ")
		ioutil.WriteFile("data.json", data, 0740)

		gobibox.Url = url
		gobibox.Login = login
		gobibox.Password = password
		gobibox.Port = port

		gobibox.Download(files)

		time.Sleep(60 * time.Second)
	}
}
