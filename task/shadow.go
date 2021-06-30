package task

import (
	"fmt"

	"github.com/diauweb/xmcl/config"
	"github.com/diauweb/xmcl/game"
)

func ApplyShadow() {
	tasks := make([]game.RemoteResource, 0, len(config.Config.Shadows))
	for k, v := range config.Config.Shadows {
		remote := game.RemoteResource{
			ID:   k,
			Type: "game_shadow",
			Path: fmt.Sprintf("./.minecraft/%s", k),
			URL:  v.URL,
			Hash: v.Hash,
		}

		tasks = append(tasks, remote)
	}

	err := DownloadGroup(tasks, 3)
	if err != nil {
		panic(err)
	}
}
