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

// APIVersion is OneSky API version which will be used
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
	"getFile":               apiEndpoint{"projects/%d/translations", "GET"},
	"postFile":              apiEndpoint{"projects/%d/files", "POST"},
	"deleteFile":            apiEndpoint{"projects/%d/files", "DELETE"},
	"listFiles":             apiEndpoint{"projects/%d/files", "GET"},
	"importTasks":           apiEndpoint{"projects/%d/import-tasks", "GET"},
	"importTask":            apiEndpoint{"projects/%d/import-tasks/%d", "GET"},
	"getTranslationsStatus": apiEndpoint{"projects/%d/translations/status", "GET"},
	"getLanguages":          apiEndpoint{"projects/%d/languages", "GET"},
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
type getLanguagesResponse struct {
	Data []Language `json:"data"`
}

// TaskData is a struct which contains informations about import task
type TaskData struct {
	ID                  int64
	OriginalID          interface{} `json:"id"`
	File                TaskFile    `json:"file"`
	StringCount         int         `json:"string_count"`
	WordCount           int         `json:"word_count"`
	Status              string      `json:"status"`
	CreateddAt          string      `json:"created_at"`
	CreateddAtTimestamp int         `json:"created_at_timestamp"`
}

// Language is a struct which contains informations about locale
type Language struct {
	Code                string `json:"code"`
	EnglishName         string `json:"english_name"`
	LocalName           string `json:"local_name"`
	CustomLocale        string `json:"custom_locale"`
	Locale              string `json:"locale"`
	Region              string `json:"region"`
	TranslationProgress string `json:"translation_progress"`
}

// TaskFile is a struct which contains informations about file of import task
type TaskFile struct {
	Name   string   `json:"name"`
	Format string   `json:"format"`
	Locale Language `json:"locale"`
}

// ImportTasksResponse is a struct which contains informations about the response from list import tasks API
type ImportTasksResponse struct {
	Data []TaskData `json:"data"`
}

// ImportTaskResponse is a struct which contains informations about the response from show an import task API
type ImportTaskResponse struct {
	Data TaskData `json:"data"`
}

// UploadData is a struct which contains informations about uploaded file
type UploadData struct {
	Name     string   `json:"name"`
	Format   string   `json:"format"`
	Language Language `json:"language"`
	Import   TaskData `json:"import"`
}

// UploadResponse is a struct which contains informations about the response from upload file API
type UploadResponse struct {
	Data UploadData `json:"data"`
}

// TranslationsStatus is a struct which contains information about a project's translation status
type TranslationsStatus struct {
	FileName    string   `json:"file_name"`
	Locale      Language `json:"locale"`
	Progress    string   `json:"progress"`
	StringCount int64    `json:"string_count"`
	WordCount   int64    `json:"word_count"`
}
type getTranslationsStatusResponse struct {
	Data TranslationsStatus `json:"data"`
}

// convertToInt64 : Convert interface{} to int64
func convertToInt64(in interface{}) (int64, error) {
	switch in.(type) {
	case string:
		return strconv.ParseInt(in.(string), 10, 64)
	case int:
		return int64(in.(int)), nil
	case int16:
		return int64(in.(int16)), nil
	case int32:
		return int64(in.(int32)), nil
	case int64:
		return in.(int64), nil
	case uint:
		return int64(in.(uint)), nil
	case uint16:
		return int64(in.(uint16)), nil
	case uint32:
		return int64(in.(uint32)), nil
	case uint64:
		return int64(in.(uint64)), nil
	case float32:
		return int64(in.(float32)), nil
	case float64:
		return int64(in.(float64)), nil
	default:
		return 0, fmt.Errorf("%s: %v", "Unable to convert value", in)
	}
}

// ImportTask : Show an import task. Parameters: import_id
func (c *Client) ImportTask(importID int64) (TaskData, error) {
	endpoint, err := getEndpoint("importTask")
	if err != nil {
		return TaskData{}, err
	}
	values := url.Values{}
	urlStr, err := endpoint.full(c, values, importID)
	if err != nil {
		return TaskData{}, err
	}

	res, err := makeRequest(endpoint.method, urlStr, nil, "")
	if err != nil {
		return TaskData{}, err
	}

	if res.StatusCode != http.StatusOK {
		return TaskData{}, fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := getResponseBodyAsString(res)
	if err != nil {
		return TaskData{}, err
	}
	aux := ImportTaskResponse{}
	err = json.Unmarshal([]byte(body), &aux)
	if err != nil {
		return TaskData{}, err
	}
	if i, err := convertToInt64(aux.Data.OriginalID); err == nil {
		// fmt.Printf("%T, %v", i, i)
		aux.Data.ID = i
	}

	return aux.Data, nil
}

// ImportTasks : List import tasks. Parameters: page: 1, per_page: 50, status: [all|completed|in-progress|failed]
// tasks, err := onesky.ImportTasks(map[string]interface{}{"per_page": 50, "status": "in-progress"})
func (c *Client) ImportTasks(params map[string]interface{}) ([]TaskData, error) {
	endpoint, err := getEndpoint("importTasks")
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("page", "1")
	values.Set("per_page", "50")
	values.Set("status", "all")
	for k, v := range params {
		values.Set(k, fmt.Sprintf("%v", v))
	}
	urlStr, err := endpoint.full(c, values)
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
	aux := ImportTasksResponse{}
	err = json.Unmarshal([]byte(body), &aux)
	if err != nil {
		return nil, err
	}
	for i := range aux.Data {
		task := &aux.Data[i]
		if id, err := convertToInt64(task.OriginalID); err == nil {
			// fmt.Printf("\n%T, %v, %T\n", id, id, task.OriginalID)
			task.ID = id
		}
	}

	return aux.Data, nil
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
func (c *Client) UploadFile(file, fileFormat, locale string, keepStrings bool) (UploadData, error) {
	endpoint, err := getEndpoint("postFile")
	if err != nil {
		return UploadData{}, err
	}

	v := url.Values{}
	v.Set("locale", locale)
	v.Set("file_format", fileFormat)
	v.Set("is_keeping_all_strings", strconv.FormatBool(keepStrings))
	urlStr, err := endpoint.full(c, v)
	if err != nil {
		return UploadData{}, err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(file)
	if err != nil {
		return UploadData{}, err
	}
	defer w.Close()

	fw, err := w.CreateFormFile("file", file)
	if err != nil {
		return UploadData{}, err
	}

	if _, err = io.Copy(fw, f); err != nil {
		return UploadData{}, err
	}

	w.Close()

	res, err := makeRequest(endpoint.method, urlStr, &b, w.FormDataContentType())
	if err != nil {
		return UploadData{}, err
	}

	if res.StatusCode != http.StatusCreated {
		return UploadData{}, fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := getResponseBodyAsString(res)
	if err != nil {
		return UploadData{}, err
	}

	aux := UploadResponse{}
	err = json.Unmarshal([]byte(body), &aux)
	if err != nil {
		return UploadData{}, err
	}
	if i, err := convertToInt64(aux.Data.Import.OriginalID); err == nil {
		aux.Data.Import.ID = i
	}

	return aux.Data, nil
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

// GetTranslationsStatus returns information about a project's translation status
func (c *Client) GetTranslationsStatus(fileName, locale string) (TranslationsStatus, error) {
	endpoint, err := getEndpoint("getTranslationsStatus")
	if err != nil {
		return TranslationsStatus{}, err
	}

	v := url.Values{}
	v.Set("file_name", fileName)
	v.Set("locale", locale)
	urlStr, err := endpoint.full(c, v)
	if err != nil {
		return TranslationsStatus{}, err
	}

	res, err := makeRequest(endpoint.method, urlStr, nil, "")
	if err != nil {
		return TranslationsStatus{}, err
	}

	if res.StatusCode != http.StatusOK {
		return TranslationsStatus{}, fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := getResponseBodyAsString(res)
	if err != nil {
		return TranslationsStatus{}, err
	}

	aux := getTranslationsStatusResponse{}
	err = json.Unmarshal([]byte(body), &aux)
	if err != nil {
		return TranslationsStatus{}, err
	}

	return aux.Data, nil
}

// GetLanguages is method on Client struct which download from OneSky service information about available languages in project
func (c *Client) GetLanguages() ([]Language, error) {
	endpoint, err := getEndpoint("getLanguages")
	if err != nil {
		return nil, err
	}

	v := url.Values{}
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

	aux := getLanguagesResponse{}
	err = json.Unmarshal([]byte(body), &aux)
	if err != nil {
		return nil, err
	}

	return aux.Data, nil
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

func (e *apiEndpoint) full(c *Client, additionalArgs url.Values, extends ...interface{}) (string, error) {
	extends = append([]interface{}{c.ProjectID}, extends...)
	urlWithProjectID := fmt.Sprintf(e.path, extends...)
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
