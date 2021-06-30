package task

import (
	"fmt"
	"sync"

	"github.com/diauweb/xmcl/game"
)

func DownloadGroup(arts []game.RemoteResource, maxRoutine int) error {
	allDepLen := len(arts)

	if maxRoutine < 1 {
		maxRoutine = 20
	}

	guard := make(chan struct{}, maxRoutine)
	var waiter sync.WaitGroup

	fetch := func(name string, art game.RemoteResource) {
		guard <- struct{}{}
		waiter.Add(1)
		go func() {
			art.Download()
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
