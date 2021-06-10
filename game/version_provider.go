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

var resolveProviders = []GameVersionProvider{
	MojangGameVersionProvider{},
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
