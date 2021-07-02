package task

import (
	"fmt"
	"sync"

	"github.com/diauweb/xmcl/remote"
)

func DownloadGroup(arts []remote.RemoteResource, maxRoutine int) error {
	allDepLen := len(arts)

	if maxRoutine < 1 {
		maxRoutine = 20
	}

	guard := make(chan struct{}, maxRoutine)
	var waiter sync.WaitGroup

	fetch := func(progressName string, art remote.RemoteResource) {
		guard <- struct{}{}
		waiter.Add(1)
		go func() {
			if !art.Validate() {
				fmt.Println(progressName)
				art.ForceDownload()
			}
			<-guard
			waiter.Done()
		}()
	}

	for i, v := range arts {
		progressName := fmt.Sprintf("[%d/%d] %s", i+1, allDepLen, v.ID)
		fetch(progressName, v)
	}

	waiter.Wait()
	return nil
}
