package ftpbox

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTP struct {
	con      *ftp.ServerConn
	homeName string
}

func Connect(url, login, password string, port int) (FTP, error) {
	var _ftp FTP
	var err error

	addr := fmt.Sprintf("%s:%d", url, port)
	_ftp.con, err = ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return _ftp, err
	}

	err = _ftp.con.Login(login, password)
	if err != nil {
		return _ftp, err
	}

	_ftp.homeName = login
	return _ftp, nil
}

func (client *FTP) setRootPath() string {
	return fmt.Sprintf("/home/%s/torrents/qbittorrent/", client.homeName)
}

func (client *FTP) Download(fileToDownload string) error {
	var err error

	rootPath := client.setRootPath()
	pathFileToDownload := path.Join(rootPath, fileToDownload)
	//localName := strings.Replace(fileToDownload, rootPath, "", -1)
	fileToDownload = path.Join("data", fileToDownload)
	localDir := path.Dir(fileToDownload)
	//_, fileName := path.Split(fileToDownload)

	if len(localDir) > 1 {
		err = os.MkdirAll(localDir, 0740)
		if err != nil {
			return err
		}
	}

	r, err := client.con.Retr(pathFileToDownload)
	if err != nil {
		return err
	}

	out, err := os.Create(fileToDownload)
	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}
	r.Close()

	if err != nil {
		return err
	}

	return nil
}

/*func (client *FTP) Download(fileToDownload string, offset int) (*ftp.Response, error) {
	rootPath := client.setRootPath()
	pathFileToDownload := path.Join(rootPath, fileToDownload)

	return client.con.RetrFrom(pathFileToDownload, uint64(offset))
}*/
