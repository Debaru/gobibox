package gobibox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	ftpbox "github.com/Debaru/gobibox/internal/ftp"
	"github.com/Debaru/gobibox/internal/qbittorent"
)

var (
	Url      string
	Login    string
	Password string
	Port     int
)

type FileGoBiBox struct {
	Name           string `json:"name"`
	Downloaded     bool   `json:"downloaded"`
	Size           int    `json:"size"`
	SizeDownloaded int    `json:"size_downloaded"`
}

func Download(files []FileGoBiBox) error {
	var err error

	// Connexion to FTP
	_ftp, err := ftpbox.Connect(Url, Login, Password, Port)
	if err != nil {
		return err
	}

	for i, file := range files {
		if !file.Downloaded {

			// Log
			fmt.Println("Download - ", file.Name)

			//local := filepath.Join(".", file.Name)
			err := _ftp.Download(file.Name)
			if err != nil {
				return err
			}

			files[i].Downloaded = true

			data, _ := json.MarshalIndent(&files, "  ", "  ")
			ioutil.WriteFile("data.json", data, 0740)
		}
	}

	return nil
}

func search(files []FileGoBiBox, torrentName string) bool {
	for _, f := range files {
		if f.Name == torrentName {
			return true
		}
	}
	return false
}

func GetFilesToDownload(torrents []qbittorent.TorrentContent) []FileGoBiBox {
	var files []FileGoBiBox
	var isFound bool

	data, err := ioutil.ReadFile("data.json")
	if err == nil {
		json.Unmarshal(data, &files)
	}

	// Browse files
	for _, torrent := range torrents {
		isFound = search(files, torrent.Name)

		if !isFound {
			var f = FileGoBiBox{
				Name:       torrent.Name,
				Size:       torrent.Size,
				Downloaded: false,
			}

			files = append(files, f)
		}
	}

	return files
}
