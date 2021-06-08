package task

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/diauweb/xmcl/game"
)

func FetchLibraries(lib *[]game.Library) {
	allDepLen := len(*lib)

	const maxRoutine = 10
	guard := make(chan struct{}, maxRoutine)
	var waiter sync.WaitGroup

	fetch := func(name string, art game.Artifact) {
		guard <- struct{}{}
		waiter.Add(1)
		go func() {
			FetchLibrary(name, art)
			<-guard
			waiter.Done()
		}()
	}

	for i, v := range *lib {
		progressName := fmt.Sprintf("[%d/%d] %s", i+1, allDepLen, v.Name)

		if !v.IsCompatible() {
			fmt.Printf("[%d/%d] version: skip dependency %s\n", i+1, allDepLen, v.Name)
			continue
		}

		if native, ok := v.Natives[runtime.GOOS]; ok {
			name := fmt.Sprintf("%s-natives", v.Name)
			progressName := fmt.Sprintf("[%d/%d] %s", i+1, allDepLen, name)
			art := v.Downloads.Classifiers[native]
			fetch(progressName, art)
		}

		fetch(progressName, v.Downloads.Artifact)
	}

	waiter.Wait()
}

func getLibPath(art game.Artifact) string {
	return fmt.Sprintf("./Managed/libraries/%s", art.Path)
}

func ValidateLibrary(art game.Artifact) error {
	path := getLibPath(art)
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

func FetchLibrary(name string, art game.Artifact) {

	path := getLibPath(art)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		panic(err)
	}

	if err := ValidateLibrary(art); err == nil {
		fmt.Printf("%s [installed]\n", name)
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

	if err := ValidateLibrary(art); err != nil {
		panic(err)
	}
}
