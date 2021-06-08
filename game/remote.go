package game

import (
	"fmt"
	"os"
	"reflect"
)

type RemoteManifest struct {
	ID       string
	Type     string
	Path     string
	URL      string
	Hash     string
	WrapType reflect.Type
}

func (r RemoteManifest) Validate() error {
	f, err := os.ReadFile(fmt.Sprintf("./Managed/%s", r.Path))
	_, _ = f, err

	return nil
}
