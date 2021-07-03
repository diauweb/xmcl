package remote

import "fmt"

const SHADOW_TYPE_BUNDLED = "bundled"
const SHADOW_TYPE_FILE = "file"

type ShadowFile struct {
	Path string `json:"path"`
	URL  string `json:"url"`
	Hash string `json:"hash"`
}

type SanityRule struct {
	Path string `json:"path"`
	Rule string `json:"rule"`
}

func (f ShadowFile) AsRemote() RemoteResource {
	return RemoteResource{
		ID:   f.Path,
		Type: "shadow_file",
		URL:  f.URL,
		Path: fmt.Sprintf("./.minecraft/%s", f.Path),
		Hash: f.Hash,
	}
}

type ShadowManifest struct {
	Type   string       `json:"type"`
	Bundle string       `json:"bundle"`
	Files  []ShadowFile `json:"files"`
	Sanity []SanityRule `json:"sanity"`
}
