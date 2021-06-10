package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Configs struct {
	Latest     string            `json:"latest"`
	UpdateFile map[string]string `json:"update_file"`
	Version    struct {
		Resolve string `json:"resolve"`
	} `json:"version"`
}

var Config Configs

func InitConfig() {
	if MODE == "DEBUG" {
		f, _ := os.ReadFile("./config.json")
		if err := json.Unmarshal(f, &Config); err != nil {
			panic(err)
		}
		return
	}

	if CONFIG_ENDPOINT == "" {
		panic("configuration endpoint is not set")
	}

	req, _ := http.NewRequest("GET", CONFIG_ENDPOINT, nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(fmt.Errorf("config: %v", err))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &Config); err != nil {
		panic(err)
	}

}
