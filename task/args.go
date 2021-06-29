package task

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/diauweb/xmcl/cli"
	"github.com/diauweb/xmcl/config"
	"github.com/diauweb/xmcl/game"
)

var MATCH_VARIABLE = regexp.MustCompile(`\$\{.*\}`)

func offlineUUID(name string) string {
	md5s := md5.Sum([]byte(fmt.Sprintf("OfflinePlayer:%s", name)))
	md5s[6] = md5s[6]&0x0f | 0x30
	md5s[8] = md5s[8]&0x3f | 0x80

	r := fmt.Sprintf("%x-%x-%x-%x-%x", md5s[0:4], md5s[4:6], md5s[6:8], md5s[8:10], md5s[10:])
	return r
}

func BuildArgs(game *game.Version) []string {
	gameJarPath, err := filepath.Abs(
		fmt.Sprintf("./Managed/libraries/net/minecraft/%[1]s/%[1]s.jar", game.ID))
	if err != nil {
		panic(err)
	}

	var classpath []string
	var args []string
	classpath = append(classpath, gameJarPath)

	for _, v := range game.Libraries {
		if !v.IsCompatible() {
			continue
		}
		p := fmt.Sprintf("./Managed/libraries/%s", v.Downloads.Artifact.Path)
		abs, err := filepath.Abs(p)
		if err != nil {
			panic(err)
		}
		classpath = append(classpath, abs)

		if len(v.Natives) > 0 {
			key, ok := v.Natives[runtime.GOOS]
			if ok {
				p2 := fmt.Sprintf("./Managed/libraries/%s", v.Downloads.Classifiers[key].Path)
				abs, err := filepath.Abs(p2)
				if err != nil {
					panic(err)
				}
				classpath = append(classpath, abs)
			}
		}
	}

	// preset jvm args
	// use g1gc
	args = append(args,
		// "--illegal-access=permit",
		"-XX:+UnlockExperimentalVMOptions",
		"-XX:+UseG1GC",
		"-XX:G1NewSizePercent=20",
		"-XX:G1ReservePercent=20",
		"-XX:MaxGCPauseMillis=50",
		"-XX:G1HeapRegionSize=16M",
		"-XX:-UseAdaptiveSizePolicy",
		"-XX:-OmitStackTraceInFastThrow",
		"-Xmn128m",
	)

	args = append(args,
		"-Dfml.ignoreInvalidMinecraftCertificates=true",
		"-Dfml.ignorePatchDiscrepancies=true",
	)

	gameDir, _ := filepath.Abs("./Managed/.minecraft")
	assetsDir, _ := filepath.Abs("./Managed/assets")
	natives := PrepareNatives(game)

	gameEnvs := map[string]string{
		"auth_player_name":  "Player",
		"version_name":      game.ID,
		"profile_name":      "Minecraft",
		"game_directory":    gameDir,
		"assets_root":       assetsDir,
		"assets_index_name": game.Assets,
		"auth_uuid":         offlineUUID("Player"),
		"version_type":      game.Type,
		"user_type":         "mojang",
		"launcher_name":     config.PRODUCT_NAME,
		"launcher_version":  config.GIT_BUILD,
		"natives_directory": natives,
		"classpath":         strings.Join(classpath, string(os.PathListSeparator)),
	}

	for k, v := range config.Config.LaunchEnvs {
		if v == "$required" {
			newv := cli.Ask(fmt.Sprintf("Input %s", k))
			if newv == "" {
				panic(k + " is null")
			}
			gameEnvs[k] = newv
		} else {
			gameEnvs[k] = v
		}
	}
	gameEnvs["auth_uuid"] = offlineUUID(gameEnvs["auth_player_name"])

	// legacy format convertion
	if game.MinecraftArguments != "" {
		sp := strings.Split(game.MinecraftArguments, " ")
		y := make([]interface{}, len(sp))
		for i, v := range sp {
			y[i] = v
		}

		game.Arguments.Game = y
		game.Arguments.JVM = []interface{}{
			"-Djava.library.path=${natives_directory}",
			"-Dminecraft.launcher.brand=${launcher_name}",
			"-Dminecraft.launcher.version=${launcher_version}",
			"-cp",
			"${classpath}",
		}
	}

	// todo: JVM and Game rule-specified arguments
	for _, v := range game.Arguments.JVM {
		s, ok := v.(string)
		if !ok {
			continue // todo
		}
		args = append(args, s)
	}

	args = append(args, game.Mainclass)

	for _, v := range game.Arguments.Game {
		s, ok := v.(string)
		if !ok {
			continue // todo
		}
		args = append(args, s)
	}

	args = append(args, config.Config.LaunchArgs...)
	if len(game.Tweakers) > 0 {
		if game.Mainclass != "net.minecraft.launchwrapper.Launch" {
			fmt.Println("game: tweakers: mainclass is not launchwrapper")
		}
		for _, v := range game.Tweakers {
			args = append(args, "--tweakClass", v)
		}
	}

	for i, v := range args {
		a := MATCH_VARIABLE.ReplaceAllFunc([]byte(v), func(b []byte) []byte {
			env := b[2 : len(b)-1]
			return []byte(gameEnvs[string(env)])
		})
		args[i] = string(a)
	}

	// fmt.Printf(">>> %v\n", args)
	return args
}
