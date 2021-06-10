package main

import (
	"fmt"

	"github.com/diauweb/xmcl/config"
	"github.com/diauweb/xmcl/game"
	"github.com/diauweb/xmcl/java"
	"github.com/diauweb/xmcl/task"
	"github.com/diauweb/xmcl/update"
	"github.com/diauweb/xmcl/utils"
)

func main() {

	fmt.Printf("%s %s\n", config.PRODUCT_NAME, config.GIT_BUILD)

	config.InitConfig()
	update.Update()

	fmt.Println("java: detect java")
	java.DownloadJava()

	resolveVersion := config.Config.Version.Resolve
	gameVersion := game.ResolveVersion(resolveVersion)

	fmt.Println("version: download dependencies")
	task.FetchLibraries(&gameVersion.Libraries)
	// fmt.Printf("%+v", gameVersion)
	fmt.Println("version: download assets")
	assets := utils.GetAssetsIndex(&gameVersion)
	task.FetchAssets(&assets)

	// game itself
	gameJarPath := fmt.Sprintf("net/minecraft/%[1]s/%[1]s.jar", gameVersion.ID)
	gameJar := gameVersion.Downloads.Client
	gameJar.Path = gameJarPath
	task.Download(gameJar, "./Managed/libraries", fmt.Sprintf("[1/1] Minecraft %s", gameVersion.ID), true)

	args := task.BuildArgs(&gameVersion)
	java.RunJava(args, "./Managed/.minecraft")
}
