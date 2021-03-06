package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/nece099/base/encrypt"
)

func DownloadFile(url string) ([]byte, string, error) {

	ext := path.Ext(url)
	filename := fmt.Sprintf("%v%v", encrypt.MD5(url), ext)

	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	pix, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return pix, filename, nil
}
