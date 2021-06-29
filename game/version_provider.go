package game

import (
	"fmt"
	"strings"
)

type GameVersionProvider interface {
	GetSchema() string
	Resolve(version string) (Version, error)
}

type MojangGameVersionProvider struct{}

func (g MojangGameVersionProvider) GetSchema() string {
	return "mojang"
}

var remoteGameManifest = RemoteManifest{
	ID:   "version_manifest",
	Type: "version_manifest",
	URL:  VERSION_MANIFEST_URL,
	Path: "versions/version_manifest.json",
}

func (g MojangGameVersionProvider) Resolve(version string) (Version, error) {
	var gameList VersionManifest
	remoteGameManifest.ForceDownload()
	remoteGameManifest.Unmarshal(&gameList)

	for _, v := range gameList.Versions {
		if v.ID == version {
			var version Version
			v.AsRemote().Unmarshal(&version)

			return version, nil
		}
	}

	return Version{}, fmt.Errorf("No such game version %s", version)
}

//////////

type ForgeGameVersionProvider struct{}

func (g ForgeGameVersionProvider) GetSchema() string {
	return "forge"
}

const FORGE_VERSION_URL = "https://meta.multimc.org/v1/net.minecraftforge/%s.json"

type multimcGameVersion struct {
	Tweakers           []string  `json:"+tweakers,omitempty"`
	Libraries          []Library `json:"libraries"`
	Mainclass          string    `json:"mainClass"`
	MinecraftArguments string    `json:"minecraftArguments"`
	Requires           []struct {
		Equals string `json:"equals"`
		UID    string `json:"uid"`
	} `json:"requires"`
	Version string `json:"version"`
	// ignore others...
}

func (g ForgeGameVersionProvider) Resolve(version string) (Version, error) {
	var remoteForgeManifest = RemoteManifest{
		ID:   version,
		Type: "version",
		Path: fmt.Sprintf("versions/forge-%s.mmc.json", version),
		URL:  fmt.Sprintf(FORGE_VERSION_URL, version),
	}
	var forgeVersion multimcGameVersion

	remoteForgeManifest.ForceDownload()
	remoteForgeManifest.Unmarshal(&forgeVersion)

	var gameVersion string
	for _, v := range forgeVersion.Requires {
		if v.UID == "net.minecraft" {
			gameVersion = v.Equals
		}
	}
	if gameVersion == "" {
		panic("found no related game version to forge")
	}

	gVersion, err := MojangGameVersionProvider{}.Resolve(gameVersion)
	if err != nil {
		return Version{}, err
	}

	if forgeVersion.MinecraftArguments != "" {
		gVersion.MinecraftArguments = forgeVersion.MinecraftArguments
	}

	if len(forgeVersion.Tweakers) > 0 {
		gVersion.Tweakers = forgeVersion.Tweakers
	}

	gVersion.Mainclass = forgeVersion.Mainclass
	gVersion.Libraries = append(gVersion.Libraries, forgeVersion.Libraries...)

	return gVersion, nil
}

var resolveProviders = []GameVersionProvider{
	MojangGameVersionProvider{},
	ForgeGameVersionProvider{},
}

func ResolveVersion(raw string) Version {
	for _, v := range resolveProviders {
		if strings.HasPrefix(raw, v.GetSchema()+":") {
			k := strings.SplitN(raw, ":", 2)

			v, err := v.Resolve(k[1])
			if err != nil {
				panic(err)
			}
			return v
		}
	}

	panic("version: can't resolve schema " + raw)
}
