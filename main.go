package main

import (
	"fmt"

	"github.com/diauweb/xmcl/config"
	"github.com/diauweb/xmcl/game"
	"github.com/diauweb/xmcl/java"
	"github.com/diauweb/xmcl/task"
	"github.com/diauweb/xmcl/update"
	"github.com/gookit/color"
)

func main() {

	color.Style{color.Bold}.Printf("%s %s\n", config.PRODUCT_NAME, config.GIT_BUILD)

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
	gameJarPath := fmt.Sprintf("net/minecraft/%[1]s/%[1]s.jar", gameVersion.ID)
	gameJar := gameVersion.Downloads.Client
	gameJar.Path = gameJarPath
	task.Download(gameJar, "./Managed/libraries", fmt.Sprintf("[1/1] Minecraft %s", gameVersion.ID), true)

	fmt.Println("config: install shadows")
	task.ApplyShadow()

	args := task.BuildArgs(&gameVersion)
	defer task.CleanNatives(&gameVersion)

	java.RunJava(args, "./Managed/.minecraft")
}
