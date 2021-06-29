package task

import (
	"fmt"

	"github.com/diauweb/xmcl/config"
	"github.com/diauweb/xmcl/game"
)

func ApplyShadow() {
	for k, v := range config.Config.Shadows {
		remote := game.RemoteManifest{
			ID:   k,
			Type: "game_shadow",
			Path: fmt.Sprintf("./.minecraft/%s", k),
			URL:  v.URL,
			Hash: v.Hash,
		}

		if !remote.Validate() {
			fmt.Println(k)
			remote.ForceDownload()
		}
	}
}
