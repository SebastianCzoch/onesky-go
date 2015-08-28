// Package onesky is go utils for working with OneSky translation service
// Copyright (c) 2015 Sebastian Czoch <sebastian@czoch.eu>. All rights reserved.
// Use of this source code is governed by a GNU v2 license found in the LICENSE file.
package onesky

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"bytes"
	"mime/multipart"
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
	"postFile": apiEndpoint{"projects/%d/files", "POST"},
}

// DownloadFile is method on Client struct which download form OneSky service choosen file as string
func (c *Client) DownloadFile(fileName, locale string) (string, error) {
	v := url.Values{}
	v.Set("locale", locale)
	v.Set("source_file_name", fileName)
	address, err := c.getFinalEndpointURL("getFile", v)
	response, err := getFileAsString(address)
	if err != nil {
		return "", err
	}

	return response, nil
}

// UploadFile is method on Client struct which upload file to OneSky service
func (c *Client) UploadFile(file, fileFormat, locale string) error {
	v := url.Values{}
	v.Set("locale", locale)
	v.Set("file_format", fileFormat)
	address, err := c.getFinalEndpointURL("postFile", v)
	
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    f, err := os.Open(file)
    if err != nil {
        return err
    }
    defer w.Close()

    fw, err := w.CreateFormFile("file", file)
    if err != nil {
        return err
    }

    if _, err = io.Copy(fw, f); err != nil {
        return err
    }
    w.Close()

    req, err := http.NewRequest("POST", address, &b)
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", w.FormDataContentType())
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        return err
    }

    if res.StatusCode <= 200 || res.StatusCode > 299 {
        return fmt.Errorf("bad status: %s", res.Status)
    }

    return nil
}

func getFileAsString(address string) (string, error) {
	response, err := http.Get(address)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	res, err := ioutil.ReadAll(response.Body)
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

func (c *Client) getURL(endpointName string) (string, error) {
	_, err := c.getURLForEndpoint(endpointName)
	if err != nil {
		return "", err
	}

	endpointURL, err := c.getURLForEndpoint(endpointName)
	if err != nil {
		return  "", err
	}

	return endpointURL, nil
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

func (c *Client) getFinalEndpointURL(endpointName string, additionalArgs url.Values) (string, error) {
	endpointURL, err := c.getURL(endpointName)
	if err != nil {
		return "", err
	}

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
