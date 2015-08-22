package onesky

import (
	"errors"
	"crypto/md5"
	"time"
	"strconv"
	"encoding/hex"
	"net/url"
	"net/http"
	"io/ioutil"
	"fmt"
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
	"getFile" : apiEndpoint{"projects/%d/translations", "GET"},
}

func (o *Options) DownloadFile(fileName, locale string) (string, error) {
	_, err := o.getUrlForEndpoint("getFile")
	if err != nil {
		return "", err
	}
	
	address, err := o.getUrlForEndpoint("getFile")
	if err != nil {
		return "", err
	}
	
	response, err := http.Get(address)
	if err != nil {
		return "", err
	}

	res, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return "", err
	}
	
	return string(res), nil
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

	urlWithProjectID := fmt.Sprintf(apiEndpoints[endpointName].path, o.ProjectID)
	address, err := url.Parse(API_ADDRESS + "/" + API_VERSION + "/" + urlWithProjectID)
	if err != nil {
		return "", errors.New("Can not parse url address!")
	}
	
	return address.String(), nil	
}

func (o *Options) getFinalEndpointUrl(endpointUrl string, additionalArgs url.Values) (string, error) {
	address, err := url.Parse(fmt.Sprintf(endpointUrl, o.ProjectID))
	if err != nil {
		return "", err
	}
	hash, timestamp := o.getAuthHashAndTime();
	
	additionalArgs.Set("api_key", o.ApiKey)
	additionalArgs.Set("timestamp", timestamp)
	additionalArgs.Set("dev_hash", hash)
	
	return address.String() + "?" + additionalArgs.Encode(), nil
}