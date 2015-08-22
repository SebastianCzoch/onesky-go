package onesky

import (
	"crypto/md5"
	"time"
	"strconv"
	"encoding/hex"
)

const API_ADDRESS = "https://platform.api.onesky.io/"
const API_VERSION = "1"

type apiEndpoint struct {
	path string
	method string
}

type api struct {
	getFile apiEndpoint
}

var apiEndpoints = map[string]apiEndpoint{
	"getFile" : apiEndpoint{"projects/translations", "GET"},
}

type Options struct {
	Secret    string
	ApiKey    string
	ProjectID int
}

func (o *Options) DownloadFile(fileName, locale string) {

}

func (o *Options) getAuthHashAndTime() (string, string) {
	hasher := md5.New()
	time := strconv.Itoa(int(time.Now().Unix()))
    hasher.Write([]byte(time + o.Secret))

	return hex.EncodeToString(hasher.Sum(nil)), time
}

func (o *Options) getUrlForAction(action string) {
	
}
