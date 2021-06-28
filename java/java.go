package java

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/diauweb/xmcl/config"
	"github.com/gookit/color"
	"github.com/schollz/progressbar/v3"
)

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		re := regexp.MustCompile(`jdk-[^/]+`)
		fpath := re.ReplaceAllString(f.Name, "")
		fpath = filepath.Join(dest, fpath)
		// fmt.Println(f.Name, fpath)

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

const JDK_PATH = "https://aka.ms/download-jdk/microsoft-jdk-16.0.1.9.1-windows-x64.zip"

func DownloadJava() {

	if runtime.GOOS != "windows" {
		return
	}

	if _, err := os.Stat("./Managed/java/bin/javaw.exe"); !os.IsNotExist(err) {
		return
	}

	f, err := os.CreateTemp("", "jvav")
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest("GET", JDK_PATH, nil)
	resp, _ := http.DefaultClient.Do(req)

	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"OpenJDK",
	)
	// fmt.Println(name)
	if _, err := io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
		panic(err)
	}

	if err := unzip(f.Name(), "./Managed/java/"); err != nil {
		panic(err)
	}
}

func GetJava() string {
	switch runtime.GOOS {
	case "windows":
		if config.Config.LocalJava {
			color.LightYellow.Println("config: Assume local java exists")
			return "javaw.exe"
		}

		if f, err := filepath.Abs("./Managed/java/bin/javaw.exe"); err == nil {
			return f
		} else {
			panic(err)
		}
	case "linux":
		return "/usr/bin/java"
	}
	panic("not implemented")
}

func RunJava(args []string, cwd string) {
	cmd := exec.Command(GetJava(), args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cwd

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}
