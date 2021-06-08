package game

type Asset struct {
	Hash string `json:"hash"`
	Size int    `json:"size"`
}

type AssetsIndex struct {
	Objects map[string]Asset `json:"objects"`
}
