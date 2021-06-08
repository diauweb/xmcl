package utils

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/diauweb/xmcl/game"
)

func download(url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(fmt.Errorf("manifest: %v", err))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return body
}

func getLocalAssetsIndex(ver *game.Version) (game.AssetsIndex, error) {
	res := ver.AssetIndex
	path := fmt.Sprintf("./Managed/assets/indexes/%s.json", res.ID)

	f, err := os.ReadFile(path)
	if err != nil {
		return game.AssetsIndex{}, err
	}

	if fmt.Sprintf("%x", sha1.Sum(f)) != res.Sha1 {
		msg := fmt.Sprintf("asset_index: %s: cache invalid", res.ID)
		fmt.Println(msg)

		return game.AssetsIndex{}, fmt.Errorf(msg)
	}

	var assets game.AssetsIndex
	if err := json.Unmarshal(f, &assets); err != nil {
		panic(err)
	}

	return assets, nil
}

func GetAssetsIndex(ver *game.Version) game.AssetsIndex {

	path := fmt.Sprintf("./Managed/assets/indexes/%s.json", ver.AssetIndex.ID)
	assets, err := getLocalAssetsIndex(ver)
	if err != nil {
		data := download(ver.AssetIndex.URL)
		if err := os.WriteFile(path, data, 0755); err != nil {
			panic(err)
		}
		assets, err = getLocalAssetsIndex(ver)
		if err != nil {
			panic(err)
		}
	}

	return assets
}
