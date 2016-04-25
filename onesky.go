// Package onesky is go utils for working with OneSky translation service
// Copyright (c) 2015 Sebastian Czoch <sebastian@czoch.eu>. All rights reserved.
// Use of this source code is governed by a GNU v2 license found in the LICENSE file.
package onesky

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// APIAddress is https address to OneSky API
const APIAddress = "https://platform.api.onesky.io"

// API Version is OneSky API version which will be used
const APIVersion = "1"

// Client is basics struct for this package contains Secret, APIKey and ProjectID which is needed to authorize in OneSky service
type Client struct {
	Secret    string
	APIKey    string
	ProjectID int
}

type apiEndpoint struct {
	path   string
	method string
}

var apiEndpoints = map[string]apiEndpoint{
	"getFile":    apiEndpoint{"projects/%d/translations", "GET"},
	"postFile":   apiEndpoint{"projects/%d/files", "POST"},
	"deleteFile": apiEndpoint{"projects/%d/files", "DELETE"},
	"listFiles":  apiEndpoint{"projects/%d/files", "GET"},
}

// FileData is a struct which contains informations about file uploaded to OneSky service
type FileData struct {
	Name                string     `json:"name"`
	FileName            string     `json:"file_name"`
	StringCount         int        `json:"string_count"`
	LastImport          LastImport `json:"last_import"`
	UpoladedAt          string     `json:"uploaded_at"`
	UpoladedAtTimestamp int        `json:"uploaded_at_timestamp"`
}

// LastImport is a struct which contains informations about last upload
type LastImport struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type listFilesResponse struct {
	Data []FileData `json:"data"`
}

// ListFiles is method on Client struct which download form OneSky service informations about uploaded files
func (c *Client) ListFiles(page, perPage int) ([]FileData, error) {
	endpoint, err := getEndpoint("listFiles")
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("page", strconv.Itoa(page))
	v.Set("per_page", strconv.Itoa(perPage))
	urlStr, err := endpoint.full(c, v)
	if err != nil {
		return nil, err
	}

	res, err := makeRequest(endpoint.method, urlStr, nil, "")
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := getResponseBodyAsString(res)
	if err != nil {
		return nil, err
	}

	aux := listFilesResponse{}
	err = json.Unmarshal([]byte(body), &aux)
	if err != nil {
		return nil, err
	}

	return aux.Data, nil
}

// DownloadFile is method on Client struct which download form OneSky service choosen file as string
func (c *Client) DownloadFile(fileName, locale string) (string, error) {
	endpoint, err := getEndpoint("getFile")
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Set("locale", locale)
	v.Set("source_file_name", fileName)
	urlStr, err := endpoint.full(c, v)
	if err != nil {
		return "", err
	}

	res, err := makeRequest(endpoint.method, urlStr, nil, "")
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := getResponseBodyAsString(res)
	if err != nil {
		return "", err
	}

	return body, nil
}

// UploadFile is method on Client struct which upload file to OneSky service
func (c *Client) UploadFile(file, fileFormat, locale string) error {
	endpoint, err := getEndpoint("postFile")
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Set("locale", locale)
	v.Set("file_format", fileFormat)
	urlStr, err := endpoint.full(c, v)
	if err != nil {
		return err
	}

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

	res, err := makeRequest(endpoint.method, urlStr, &b, w.FormDataContentType())
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	return nil
}

// DeleteFile is method on Client struct which remove file from OneSky service
func (c *Client) DeleteFile(fileName string) error {
	endpoint, err := getEndpoint("deleteFile")
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Set("file_name", fileName)
	urlStr, err := endpoint.full(c, v)
	if err != nil {
		return err
	}

	res, err := makeRequest(endpoint.method, urlStr, nil, "")
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	return nil
}

func makeRequest(method, urlStr string, body io.Reader, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getResponseBodyAsString(response *http.Response) (string, error) {
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

func (e *apiEndpoint) full(c *Client, additionalArgs url.Values) (string, error) {
	urlWithProjectID := fmt.Sprintf(e.path, c.ProjectID)
	address, err := url.Parse(APIAddress + "/" + APIVersion + "/" + urlWithProjectID)
	if err != nil {
		return "", err
	}

	hash, timestamp := c.getAuthHashAndTime()
	additionalArgs.Set("api_key", c.APIKey)
	additionalArgs.Set("timestamp", timestamp)
	additionalArgs.Set("dev_hash", hash)

	return address.String() + "?" + additionalArgs.Encode(), nil
}

func getEndpoint(name string) (apiEndpoint, error) {
	endpoint, ok := apiEndpoints[name]
	if !ok {
		return apiEndpoint{}, fmt.Errorf("endpoint %s not found", name)
	}

	return endpoint, nil
}
