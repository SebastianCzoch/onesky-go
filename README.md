# onesky-go [![Build Status](https://travis-ci.org/SebastianCzoch/onesky-go.svg?branch=master)](https://travis-ci.org/SebastianCzoch/onesky-go) [![GoDoc](https://godoc.org/github.com/SebastianCzoch/onesky-go?status.svg)](https://godoc.org/github.com/SebastianCzoch/onesky-go) [![Issue Stats](http://issuestats.com/github/SebastianCzoch/onesky-go/badge/pr?style=flat-square)](http://issuestats.com/github/SebastianCzoch/onesky-go)



Go utils for working with [OneSky](http://www.oneskyapp.com/) translation service.

## Examples
### Example 1 - Download file

```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	fmt.Println(onesky.DownloadFile("filename", "locale"))
}
```

### Example 2 - Upload file

```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	err := onesky.UploadFile("messages.yml", "YAML", "en-US")
	if err != nil {
		fmt.Println("Can not upload file")
	}
}
```

### Example 3 - Delete file

```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	err := onesky.DeleteFile("messages.yml")
	if err != nil {
		fmt.Println("Can not delete file")
	}
}
```

## Install

```
$ go get github.com/SebastianCzoch/onesky-go
````

## API

### (c *Client) DownloadFile(fileName, locale string) (string, error)
Downloads translation file from OneSky.

Returns file content via string.

### (c *Client) UploadFile(file, fileFormat, locale string) error
Upload translation file to OneSky.
* `file` should be a full path to file

### (c *Client) DeleteFile(fileName string) error
Permanently remove file from OneSky service (with translations)!

## Tests

```
$ go test
````

## License

[GNU v2](./LICENSE)
