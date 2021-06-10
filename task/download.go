package task

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/diauweb/xmcl/game"
)

func ValidateFile(art game.Artifact, baseDir string) error {
	path := fmt.Sprintf("%s/%s", baseDir, art.Path)

	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("validate: %v", err)
	}
	hash := sha1.Sum(f)

	if fmt.Sprintf("%x", hash) != art.Sha1 {
		return fmt.Errorf("validate: sha1 mismatch")
	}

	return nil
}

func ValidateFileSize(art game.Artifact, baseDir string) error {
	path := fmt.Sprintf("%s/%s", baseDir, art.Path)

	f, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	if f.Size() != int64(art.Size) {
		return fmt.Errorf("validate: size mismatch")
	}

	return nil
}

func Download(art game.Artifact, baseDir string, name string, strictValidate bool) {

	path := fmt.Sprintf("%s/%s", baseDir, art.Path)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		panic(err)
	}

	validator := ValidateFileSize
	if strictValidate {
		validator = ValidateFile
	}

	if err := validator(art, baseDir); err == nil {
		// fmt.Printf("%s [installed]\n", name)
		return
	}

	req, _ := http.NewRequest("GET", art.URL, nil)
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

	// bar := progressbar.DefaultBytes(
	// 	resp.ContentLength,
	// 	name,
	// )
	fmt.Println(name)
	if _, err := io.Copy(f, resp.Body); err != nil {
		panic(err)
	}

	if err := validator(art, baseDir); err != nil {
		panic(err)
	}
}

type Task struct {
	Name     string
	Artifact game.Artifact
}

func DownloadGroup(arts []Task, baseDir string) error {
	allDepLen := len(arts)

	const maxRoutine = 20
	guard := make(chan struct{}, maxRoutine)
	var waiter sync.WaitGroup

	fetch := func(name string, art game.Artifact) {
		guard <- struct{}{}
		waiter.Add(1)
		go func() {
			Download(art, baseDir, name, false)
			<-guard
			waiter.Done()
		}()
	}

	for i, v := range arts {
		progressName := fmt.Sprintf("[%d/%d] %s", i+1, allDepLen, v.Name)
		fetch(progressName, v.Artifact)
	}

	waiter.Wait()
	return nil
}
