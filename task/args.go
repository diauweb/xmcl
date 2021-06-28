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

	args = append(args, "-cp", strings.Join(classpath, string(os.PathListSeparator)))
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

	// brand infomations
	args = append(args,
		"-Dminecraft.launcher.version="+config.GIT_BUILD,
		"-Dminecraft.launcher.brand="+config.PRODUCT_NAME,
	)

	args = append(args, game.Mainclass)

	gameDir, _ := filepath.Abs("./Managed/.minecraft")
	assetsDir, _ := filepath.Abs("./Managed/assets")

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

	for _, v := range game.Arguments.Game {
		s, ok := v.(string)
		if !ok {
			continue
		}
		a := MATCH_VARIABLE.ReplaceAllFunc([]byte(s), func(b []byte) []byte {
			env := b[2 : len(b)-1]
			return []byte(gameEnvs[string(env)])
		})
		args = append(args, string(a))
	}

	args = append(args, config.Config.LaunchArgs...)

	return args
}
