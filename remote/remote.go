package remote

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type RemoteResource struct {
	ID   string
	Type string
	Path string
	URL  string
	Hash string
}

func (r RemoteResource) realpath() string {
	return fmt.Sprintf("./Managed/%s", r.Path)
}
func (r RemoteResource) Validate() bool {
	if r.Hash == "" {
		_, err := os.Stat(r.realpath())
		return !os.IsNotExist(err)
	}

	f, err := os.Open(r.realpath())
	if err != nil {
		return false
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}

	hash := h.Sum(nil)
	if fmt.Sprintf("%x", hash) != r.Hash {
		fmt.Printf("remote_manifest: invalidate: %v found %x\n", r, hash)
	}
	return fmt.Sprintf("%x", hash) == r.Hash
}

func (r RemoteResource) ForceDownload() {
	r.ForceDownloadThreads(1)
}

func (r RemoteResource) ForceDownloadThreads(thread int) {
	path := r.realpath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		panic(err)
	}

	if thread > 1 {
		if err := r.downloadmulti(int64(thread)); err != nil {
			panic(err)
		}
	}

	req, _ := http.NewRequest("GET", r.URL, nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("request %s return status code %d", r.URL, resp.StatusCode))
	}

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

func (r RemoteResource) DownloadThreads(threads int) {

	if r.Validate() {
		return
	}

	r.ForceDownloadThreads(threads)
}

func (r RemoteResource) Download() {

	if r.Validate() {
		return
	}

	r.ForceDownload()
}

func (r RemoteResource) Unmarshal(o interface{}) {
	r.Download()

	f, err := os.ReadFile(r.realpath())
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(f, o); err != nil {
		panic(err)
	}
}
