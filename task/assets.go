package task

import (
	"fmt"

	"github.com/diauweb/xmcl/game"
	"github.com/diauweb/xmcl/remote"
)

const OBJECT = "https://resources.download.minecraft.net/%s"

func FetchAssets(assets *game.AssetsIndex) {
	tasks := make([]remote.RemoteResource, 0, len(assets.Objects))
	for k, v := range assets.Objects {
		path := fmt.Sprintf("%s/%s", v.Hash[:2], v.Hash)
		rm := remote.RemoteResource{
			ID:   k,
			Type: "assets_object",
			Path: fmt.Sprintf("assets/objects/%s/%s", v.Hash[:2], v.Hash),
			URL:  fmt.Sprintf(OBJECT, path),
			Hash: v.Hash,
		}
		tasks = append(tasks, rm)
	}

	if err := DownloadGroup(tasks, 20); err != nil {
		panic(err)
	}
}
