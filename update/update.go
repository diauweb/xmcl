package update

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/diauweb/xmcl/config"
	"github.com/staktrace/go-update"
)

func Update() {
	if config.MODE == "DEBUG" {
		fmt.Println("Debug mode, skip update")
		return
	}

	if f, err := os.Stat("./.xmcl.exe.old"); err == nil {
		os.Remove(f.Name())
	}

	if config.Config.Latest == config.GIT_BUILD {
		return
	}

	fmt.Printf("Latest version is %s, running %s\n", config.Config.Latest, config.GIT_BUILD)
	resp, err := http.Get(config.Config.UpdateFile[runtime.GOOS])
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Update completed. Please restart launcher.")
	os.Exit(0)
}
