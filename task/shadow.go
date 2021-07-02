package task

import (
	"crypto/sha1"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/diauweb/xmcl/config"
	"github.com/diauweb/xmcl/remote"
)

type shadowSanityRuleProcessor func(path string, file fs.DirEntry) error

func newProvisionedSanityRuleProcessor(manifest *remote.ShadowManifest, rule remote.SanityRule) shadowSanityRuleProcessor {

	checklist := map[string]*remote.ShadowFile{}

	for _, v := range manifest.Files {
		if strings.HasPrefix(v.Path, rule.Path) {
			abspath, _ := filepath.Abs(fmt.Sprintf("./Managed/.minecraft/%s", v.Path))
			checklist[abspath] = &v
		}
	}

	scopeAbsPath, _ := filepath.Abs(rule.Path)
	return func(path string, file fs.DirEntry) error {
		abspath, _ := filepath.Abs(path)
		if file.IsDir() {
			return nil
		}

		if strings.HasPrefix(abspath, scopeAbsPath) {
			ck, ok := checklist[abspath]
			if !ok {
				fmt.Printf("unknown file: %s", abspath)
				// os.Remove(abspath)
				return nil
			}

			f, err := os.ReadFile(abspath)
			if err != nil {
				panic(err)
			}

			hash := sha1.Sum(f)
			if fmt.Sprintf("%x", hash) != ck.Hash {
				return fmt.Errorf("sanity check failed %s sha1=%x want=%s", path, hash, ck.Hash)
			}
		}

		return nil
	}
}

func ApplyShadow() {

	// temporary disable
	return

	tasks := []remote.RemoteResource{}

	rManifest := remote.RemoteResource{
		URL:  config.Config.ShadowManifest,
		Path: "cache/shadow_manifest.json",
	}
	rManifest.ForceDownload()

	var shadowManifest remote.ShadowManifest
	rManifest.Unmarshal(&shadowManifest)

	_ = filepath.WalkDir("./Managed/.minecraft", func(path string, d fs.DirEntry, _ error) error {
		i, _ := d.Info()
		fmt.Printf(":: %s %v\n", path, i)

		return nil
	})
	panic("")
	err := DownloadGroup(tasks, 3)
	if err != nil {
		panic(err)
	}
}
