package main

import (
	"fmt"
	"os"

	"github.com/diauweb/xmcl/config"
	"github.com/diauweb/xmcl/game"
	"github.com/diauweb/xmcl/java"
	"github.com/diauweb/xmcl/remote"
	"github.com/diauweb/xmcl/task"
	"github.com/diauweb/xmcl/update"
	"github.com/gookit/color"
)

func main() {

	color.Style{color.Bold}.Printf("%s %s\n", config.PRODUCT_NAME, config.GIT_BUILD)

	if err := os.MkdirAll("./Managed/cache", 0755); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("./Managed/.minecraft", 0755); err != nil {
		panic(err)
	}

	config.InitConfig()
	update.Update()

	fmt.Println("java: detect java")
	java.DownloadJava()

	resolveVersion := config.Config.Version.Resolve
	gameVersion := game.ResolveVersion(resolveVersion)

	fmt.Println("version: download dependencies")
	task.FetchLibraries(&gameVersion.Libraries)

	fmt.Println("version: download assets")
	var assets game.AssetsIndex
	gameVersion.AssetIndex.AsRemote().Unmarshal(&assets)

	task.FetchAssets(&assets)

	// game itself
	gameJar := remote.RemoteResource{
		ID:   gameVersion.ID,
		Type: "game",
		URL:  gameVersion.Downloads.Client.URL,
		Path: fmt.Sprintf("libraries/net/minecraft/%[1]s/%[1]s.jar", gameVersion.ID),
		Hash: gameVersion.Downloads.Client.Sha1,
	}
	fmt.Printf("[1/1] Minecraft %s\n", gameVersion.ID)
	gameJar.Download()

	fmt.Println("config: install shadows")
	task.ApplyShadow()

	args := task.BuildArgs(&gameVersion)
	defer task.CleanNatives(&gameVersion)
	defer os.RemoveAll("./Managed/cache")

	java.RunJava(args, "./Managed/.minecraft")
}
