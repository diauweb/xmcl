package game

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type RemoteManifest struct {
	ID   string
	Type string
	Path string
	URL  string
	Hash string
}

func (r RemoteManifest) realpath() string {
	return fmt.Sprintf("./Managed/%s", r.Path)
}
func (r RemoteManifest) Validate() bool {
	if r.Hash == "" {
		_, err := os.Stat(r.realpath())
		return !os.IsNotExist(err)
	}

	f, err := os.ReadFile(r.realpath())
	if err != nil {
		return false
	}

	hash := sha1.Sum(f)
	// debug
	if fmt.Sprintf("%x", hash) != r.Hash {
		fmt.Printf("remote_manifest: invalidate: %v found %x\n", r, hash)
	}
	return fmt.Sprintf("%x", hash) == r.Hash
}

func (r RemoteManifest) ForceDownload() {
	path := r.realpath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		panic(err)
	}

	req, _ := http.NewRequest("GET", r.URL, nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	f, err1 := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	if err1 != nil {
		panic(err1)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		panic(err)
	}

	if !r.Validate() {
		panic("r: validation failure")
	}
}

func (r RemoteManifest) Download() {

	if r.Validate() {
		return
	}

	r.ForceDownload()
}

func (r RemoteManifest) Unmarshal(o interface{}) {
	r.Download()

	f, err := os.ReadFile(r.realpath())
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(f, o); err != nil {
		panic(err)
	}
}
