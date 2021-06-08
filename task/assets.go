package task

import (
	"fmt"

	"github.com/diauweb/xmcl/game"
)

const OBJECT = "https://resources.download.minecraft.net/%s"

func FetchAssets(assets *game.AssetsIndex) {
	tasks := make([]Task, 0, len(assets.Objects))
	for k, v := range assets.Objects {
		path := fmt.Sprintf("%s/%s", v.Hash[:2], v.Hash)
		// fmt.Printf("%s %v\n", k, path)
		tasks = append(tasks, Task{Name: k, Artifact: game.Artifact{
			Sha1: v.Hash,
			Size: v.Size,
			Path: path,
			URL:  fmt.Sprintf(OBJECT, path),
		}})
	}

	DownloadGroup(tasks, "./Managed/assets/objects")
}
