package task

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/diauweb/xmcl/config"
	"github.com/diauweb/xmcl/remote"
)

// provisionedSanityCheck deletes all files not in file list from specified path
// file integrity is handled in download process
func provisionedSanityCheck(manifest *remote.ShadowManifest, rule remote.SanityRule) {

	getpath := func(path string) string {
		p, _ := filepath.Abs(fmt.Sprintf("./Managed/.minecraft/%s", path))
		return p
	}

	dirpath := getpath(rule.Path)
	checklist := map[string]remote.ShadowFile{}

	for _, v := range manifest.Files {
		if strings.HasPrefix(v.Path, rule.Path) {
			abspath := getpath(v.Path)
			checklist[abspath] = v
		}
	}

	err := filepath.WalkDir("./Managed/.minecraft", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		abspath, _ := filepath.Abs(path)

		if !strings.HasPrefix(abspath, dirpath) {
			return err
		}

		_, ok := checklist[abspath]
		if !ok {
			fmt.Printf("integrity: %s is unknown\n", filepath.Base(abspath))
			os.Remove(abspath)
		}
		return err
	})
	if err != nil {
		panic(err)
	}
}

func ApplyShadow() {

	if config.Config.ShadowManifest == "" {
		return
	}

	tasks := []remote.RemoteResource{}

	rManifest := remote.RemoteResource{
		URL:  config.Config.ShadowManifest,
		Path: "cache/shadow_manifest.json",
	}
	rManifest.ForceDownload()

	var shadowManifest remote.ShadowManifest
	rManifest.Unmarshal(&shadowManifest)

	if shadowManifest.Type != "files" {
		panic(fmt.Errorf("unrecognized sanity type %s", shadowManifest.Type))
	}

	for _, v := range shadowManifest.Files {
		tasks = append(tasks, v.AsRemote())
	}

	err := DownloadGroup(tasks, 10)
	if err != nil {
		panic(err)
	}

	for _, v := range shadowManifest.Sanity {
		switch v.Rule {
		case "provisioned":
			provisionedSanityCheck(&shadowManifest, v)
		default:
			panic(fmt.Errorf("unrecognized sanity type %s", v.Rule))
		}
	}

}
