# urlreader

A simple module for opening and returning a stream from an URL.

Supports HTTP Basic Auth, header token and Oauth2-token.  Other headers may also be set.

A custom return status may be specified if a successful return status is unequal
to the HTTP status code OK.

The stream **must** be closed after the data has been consumed.

## Dependencies

* [Go](https://golang.org/) >= 1.15

## Example

```go

package main

import (
    "log"
    "time"
    "bufio"
    "context"
    "encoding/json"
    "net/url"
    "git.it.ntnu.no/go/urlreader"
)

type Obj struct {
    Somestr string `json:"somestring"`
    Someint int    `json:"someint"`
}

func main() {
    ur, err := urlreader.NewURLReader("https://...")
    if err != nil {
        log.Fatal(err)
    }
    proxyURL, err := url.Parse("socks5://...")
    if err != nil {
        log.Fatal(err)
    }
    ur.Proxy(proxyURL).Header("Accept", "application/json;version=1")
    ur.Header("Accept-Charset", "utf-8")
    ur.Header("Cache-Control", "private, no-cache, no-store, must-revalidate, max-age=0")

    ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
    defer cancel()
    // Never ever, use nil as a context!
    // The context should have a time-limit to prevent the request from
    // blocking too long.
    rdr, err := ur.BasicAuth("...", "...").Open(ctx)
    // rdr, err := ur.OAuth2HeaderToken("...").Open(ctx)
    // rdr, err := ur.Header("X-SomeAuthToken", "...").Open(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer rdr.Close()

    br := bufio.NewReaderSize(rdr, 1024*8)
    obj := make([]Obj, 0)
    if err := json.NewDecoder(br).Decode(&obj); err != nil {
        log.Fatal(err)
    }
    // ...
}
```
