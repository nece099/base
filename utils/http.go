package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/zen099/onetube/server/base/crypto"
)

func DownloadFile(url string) ([]byte, string, error) {

	ext := path.Ext(url)
	filename := fmt.Sprintf("%v%v", crypto.MD5(url), ext)

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
