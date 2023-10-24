package ftpbox

import (
	"encoding/json"
	"io/ioutil"
)

type FileGoBiBox struct {
	Path       string `json:"path"`
	Size       uint64 `json:"size"`
	Downloaded bool   `json:"downloaded"`
}

type Config struct {
	Url      string        `json:"url"`
	Port     int           `json:"port"`
	Login    string        `json:"login"`
	Password string        `json:"password"`
	List     []FileGoBiBox `json:"list"`
}

func (c *Config) Search(path string) bool {
	for _, file := range c.List {
		if file.Path == path {
			return true
		}
	}
	return false
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(&c, "", "")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("gobibox.json", data, 0740)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Load() error {
	data, err := ioutil.ReadFile("gobibox.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	return nil
}
