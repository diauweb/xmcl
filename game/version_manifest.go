package game

import (
	"fmt"
	"strings"
	"time"
)

const VERSION_MANIFEST_URL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

type VersionMeta struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	URL         string    `json:"url"`
	Time        time.Time `json:"time"`
	Releasetime time.Time `json:"releaseTime"`
}

type VersionManifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []VersionMeta `json:"versions"`
}

func (v VersionMeta) AsRemote() RemoteManifest {
	s := strings.Split(v.URL, "/")
	hash := s[len(s)-2]

	return RemoteManifest{
		ID:   v.ID,
		Type: "version",
		URL:  v.URL,
		Hash: hash,
		Path: fmt.Sprintf("versions/%s.json", v.ID),
	}
}
