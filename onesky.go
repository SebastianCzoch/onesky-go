package onesky

import (
	"errors"
	"crypto/md5"
	"time"
	"strconv"
	"encoding/hex"
	"net/url"
)

const API_ADDRESS = "https://platform.api.onesky.io"
const API_VERSION = "1"

type Options struct {
	Secret    string
	ApiKey    string
	ProjectID int
}

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

func (o *Options) DownloadFile(fileName, locale string) {

}

func (o *Options) getAuthHashAndTime() (string, string) {
	hasher := md5.New()
	time := strconv.Itoa(int(time.Now().Unix()))
    hasher.Write([]byte(time + o.Secret))

	return hex.EncodeToString(hasher.Sum(nil)), time
}

func (o *Options) getUrlForEndpoint(endpointName string) (string, error) {
	if _, ok := apiEndpoints[endpointName]; !ok {
		return "", errors.New("Endpoint not found!")
	}

	address, err := url.Parse(API_ADDRESS + "/" + API_VERSION + "/" + apiEndpoints[endpointName].path)
	if err != nil {
		return "", errors.New("Can not parse url address!")
	}
	
	return address.String(), nil	
}
