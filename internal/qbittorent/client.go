package qbittorent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var qbittorentC qbittorent

type qbittorent struct {
	url      string
	login    string
	password string
}

type Torrent struct {
	AddedOn           int     `json:"added_on"`
	AmountLeft        int     `json:"amount_left"`
	AutoTmm           bool    `json:"auto_tmm"`
	Availability      float32 `json:"availability"`
	Category          string  `json:"category"`
	Completed         int     `json:"completed"`
	CompletionOn      int     `json:"completion_on"`
	ContentPath       string  `json:"content_path"`
	Dllimit           int     `json:"dl_limit"`
	Dlspeed           int     `json:"dlspeed"`
	Downloaded        int     `json:"downloaded"`
	DownloadedSession int     `json:"downloaded_session"`
	Eta               int     `json:"eta"`
	FLPiecePrio       bool    `json:"f_l_piece_prio"`
	ForceStart        bool    `json:"force_start"`
	Hash              string  `json:"hash"`
	LastActivity      int     `json:"last_activity"`
	MagnetUri         string  `json:"magnet_uri"`
	MaxRatio          float32 `json:"max_ratio"`
	MaxSeedingTime    int     `json:"max_seeding_time"`
	Name              string  `json:"name"`
	NumComplete       int     `json:"num_complete"`
	NumLeechs         int     `json:"num_leechs"`
	NumSeeds          int     `json:"num_seeds"`
	Priority          int     `json:"priority"`
	Progress          float32 `json:"progress"`
	Ratio             float32 `json:"ratio"`
	RatioLimit        float32 `json:"ratio_limit"`
	SavePath          string  `json:"save_path"`
	SedingTime        int     `json:"seeding_time"`
	SedingTimeLimit   int     `json:"seeding_time_limit"`
	SeenComplete      int     `json:"seen_complete"`
	SeqDl             bool    `json:"seq_dl"`
	Size              int     `json:"size"`
	State             string  `json:"state"`
	SuperSeeding      bool    `json:"super_seeding"`
	Tags              string  `json:"tags"`
	TimeActive        int     `json:"time_active"`
	TotalSize         int     `json:"total_size"`
	Tracker           string  `json:"tracker"`
	UpLimit           int     `json:"up_limit"`
	Uploaded          int     `json:"uploaded"`
	UploadedSession   int     `json:"uploaded_session"`
	UpSpeed           int     `json:"upspeed"`
}

type TorrentContent struct {
	Index        int     `json:"index"`
	Name         string  `json:"name"`
	Size         int     `json:"size"`
	Progress     float32 `json:"Progress"`
	Priority     int     `json:"priority"`
	IsSeed       bool    `json:"is_seed"`
	PieceRange   int     `json:"-"`
	Availability float32 `json:"availability"`
}

func (qbittorentC *qbittorent) prepareRequest(requestURL string) (*http.Response, error) {
	var reponse *http.Response
	var request *http.Request
	var httpC http.Client
	var err error

	requestURL = fmt.Sprintf("%s%s", qbittorentC.url, requestURL)
	request, err = http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return reponse, err
	}

	// set BasicAuth
	request.SetBasicAuth(qbittorentC.login, qbittorentC.password)

	// set request
	return httpC.Do(request)
}

// Connect to qBittorent
func Connect(url string, login string, password string) error {

	qbittorentC.url = fmt.Sprintf("https://%s/qbittorrent/api/v2/", url)
	qbittorentC.login = login
	qbittorentC.password = password

	reponse, err := qbittorentC.prepareRequest("auth/login")
	if err != nil {
		return err
	}

	if reponse.StatusCode != 200 {
		return errors.New("Login or password incorrect")
	}

	return nil
}

// List of torrent completed
func TorrentCompleted() ([]TorrentContent, error) {
	var torrents []Torrent
	var torrentContent []TorrentContent
	var torrentContentTotal []TorrentContent

	reponse, err := qbittorentC.prepareRequest("torrents/info?filter=completed")
	if err != nil {
		return torrentContentTotal, err
	}

	data, err := ioutil.ReadAll(reponse.Body)
	if err != nil {
		return torrentContentTotal, err
	}

	err = json.Unmarshal(data, &torrents)
	if err != nil {
		return torrentContentTotal, err
	}

	// Browse completed torrent to get their files
	for _, torrent := range torrents {
		torrentContent, err = torrentFiles(torrent.Hash)
		if err != nil {
			return torrentContentTotal, err
		}
		torrentContentTotal = append(torrentContentTotal, torrentContent...)
	}

	return torrentContentTotal, err
}

func torrentFiles(hash string) ([]TorrentContent, error) {
	var torrentContent []TorrentContent
	var url string

	url = fmt.Sprintf("torrents/files?hash=%s", hash)
	response, err := qbittorentC.prepareRequest(url)
	if err != nil {
		return torrentContent, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return torrentContent, err
	}

	err = json.Unmarshal(data, &torrentContent)
	if err != nil {
		return torrentContent, err
	}
	return torrentContent, nil
}
