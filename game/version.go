package game

import (
	"runtime"
	"time"
)

type Version struct {
	Arguments              Arguments      `json:"arguments"`
	AssetIndex             AssetIndexLink `json:"assetIndex"`
	Assets                 string         `json:"assets"`
	ComplianceLevel        int            `json:"complianceLevel"`
	Downloads              Downloads      `json:"downloads"`
	ID                     string         `json:"id"`
	JavaVersion            JavaVersion    `json:"javaVersion"`
	Libraries              []Library      `json:"libraries"`
	Logging                Logging        `json:"logging"`
	Mainclass              string         `json:"mainClass"`
	Minimumlauncherversion int            `json:"minimumLauncherVersion"`
	Releasetime            time.Time      `json:"releaseTime"`
	Time                   time.Time      `json:"time"`
	Type                   string         `json:"type"`
}

type Arguments struct {
	Game []interface{} `json:"game"`
	JVM  []interface{} `json:"jvm"`
}

type AssetIndexLink struct {
	ID        string `json:"id"`
	Sha1      string `json:"sha1"`
	Size      int    `json:"size"`
	Totalsize int    `json:"totalSize"`
	URL       string `json:"url"`
}

type Downloads struct {
	Client         Artifact `json:"client"`
	ClientMappings Artifact `json:"client_mappings"`
	Server         Artifact `json:"server"`
	ServerMappings Artifact `json:"server_mappings"`
}

type JavaVersion struct {
	Component    string `json:"component"`
	Majorversion int    `json:"majorVersion"`
}
type Artifact struct {
	Path string `json:"path,omitempty"`
	Sha1 string `json:"sha1"`
	Size int    `json:"size"`
	URL  string `json:"url"`
}

type Logging struct {
	Client struct {
		Argument string `json:"argument"`
		File     struct {
			ID   string `json:"id"`
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"file"`
		Type string `json:"type"`
	} `json:"client"`
}

type LibraryRule struct {
	Action string `json:"action"`
	OS     struct {
		Name string `json:"name"`
	} `json:"os,omitempty"`
}

type Library struct {
	Name      string `json:"name"`
	Downloads struct {
		Artifact    Artifact            `json:"artifact"`
		Classifiers map[string]Artifact `json:"classifiers,omitempty"`
	} `json:"downloads"`
	Natives map[string]string `json:"natives,omitempty"`
	Rules   []LibraryRule     `json:"rules,omitempty"`
}

func (lib *Library) IsCompatible() bool {
	compatible := true
	for _, r := range lib.Rules {
		switch r.Action {
		case "allow":
			if r.OS.Name != "" && r.OS.Name != runtime.GOOS {
				compatible = false
			}
		case "disallow":
			if r.OS.Name == runtime.GOOS {
				compatible = false
			}
		}
	}

	return compatible
}
