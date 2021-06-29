package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type ShadowFile struct {
	Hash string `json:"hash"`
	URL  string `json:"url"`
}

type Configs struct {
	Latest     string            `json:"latest"`
	UpdateFile map[string]string `json:"update_file"`
	Version    struct {
		Resolve string `json:"resolve"`
	} `json:"version"`
	LaunchEnvs map[string]string     `json:"launch_envs"`
	LaunchArgs []string              `json:"launch_args"`
	LocalJava  bool                  `json:"local_java"`
	Shadows    map[string]ShadowFile `json:"shadows"`
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
