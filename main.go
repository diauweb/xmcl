package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/diauweb/xmcl/game"
	"github.com/diauweb/xmcl/java"
	"github.com/diauweb/xmcl/task"
	"github.com/diauweb/xmcl/utils"
)

func main() {
	fmt.Println("java: detect java")
	java.DownloadJava()

	dat, err := ioutil.ReadFile("./Managed/versions/1.17-pre1.json")
	if err != nil {
		panic(err)
	}
	var gameVersion game.Version

	if err1 := json.Unmarshal([]byte(dat), &gameVersion); err1 != nil {
		panic(err1)
	}

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
