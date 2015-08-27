// Go utils for working with OneSky translation service
package onesky

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// APIAddress is https address to OneSky API
const APIAddress = "https://platform.api.onesky.io"

// API Version is OneSky API version which will be used
const APIVersion = "1"

// Client is basics struct for this package
type Client struct {
	Secret    string
	APIKey    string
	ProjectID int
}

type apiEndpoint struct {
	path   string
	method string
}

type api struct {
	getFile apiEndpoint
}

var apiEndpoints = map[string]apiEndpoint{
	"getFile": apiEndpoint{"projects/%d/translations", "GET"},
}

// DownloadFile is method on Client struct which download form OneSky service choosen file as string
func (c *Client) DownloadFile(fileName, locale string) (string, error) {
	_, err := c.getURLForEndpoint("getFile")
	if err != nil {
		return "", err
	}

	endpointURL, err := c.getURLForEndpoint("getFile")
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Set("locale", locale)
	v.Set("source_file_name", fileName)
	address, err := c.getFinalEndpointURL(endpointURL, v)
	res, err := getFileAsString(address)
	if err != nil {
		return "", nil
	}
	return res, nil
}

func getFileAsString(address string) (string, error) {
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

func (c *Client) getAuthHashAndTime() (string, string) {
	hasher := md5.New()
	time := strconv.Itoa(int(time.Now().Unix()))
	hasher.Write([]byte(time + c.Secret))

	return hex.EncodeToString(hasher.Sum(nil)), time
}

func (c *Client) getURLForEndpoint(endpointName string) (string, error) {
	if _, ok := apiEndpoints[endpointName]; !ok {
		return "", errors.New("endpoint not found")
	}

	urlWithProjectID := fmt.Sprintf(apiEndpoints[endpointName].path, c.ProjectID)
	address, err := url.Parse(APIAddress + "/" + APIVersion + "/" + urlWithProjectID)
	if err != nil {
		return "", errors.New("can not parse url address")
	}

	return address.String(), nil
}

func (c *Client) getFinalEndpointURL(endpointURL string, additionalArgs url.Values) (string, error) {
	address, err := url.Parse(endpointURL)
	if err != nil {
		return "", err
	}
	hash, timestamp := c.getAuthHashAndTime()

	additionalArgs.Set("api_key", c.APIKey)
	additionalArgs.Set("timestamp", timestamp)
	additionalArgs.Set("dev_hash", hash)

	return address.String() + "?" + additionalArgs.Encode(), nil
}
