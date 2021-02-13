// Package urlreader open a connection to an URL and return an io.ReadCloser after Open
package urlreader

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// URLReader with default options
type URLReader struct {
	loc          string
	returnStatus int
	req          *http.Request
	trans        *http.Transport
}

// NewURLReader return a new URLReader-object
func NewURLReader(loc string) (*URLReader, error) {
	request, err := http.NewRequest(http.MethodGet, loc, nil)
	if err != nil {
		return nil, err
	}
	return &URLReader{loc: loc, returnStatus: http.StatusOK, req: request}, nil
}

// BasicAuth sets the basic auth header in the request
func (u *URLReader) BasicAuth(user string, pw string) *URLReader {
	u.req.SetBasicAuth(user, pw)
	return u
}

// OAuth2HeaderToken sets the Authorization HTTP-header to the given OAuth2 token.
// An already existing header will be overwritten.
func (u *URLReader) OAuth2HeaderToken(token string) *URLReader {
	u.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return u
}

// Header set a HTTP-header in the request.  An already existing header with the same name
// will be overwritten with the given value.
func (u *URLReader) Header(name string, value string) *URLReader {
	u.req.Header.Set(name, value)
	return u
}

// Proxy sets the proxy-URL, supported proxy-schemes are socks5, http and https.
func (u *URLReader) Proxy(fixedURL *url.URL) *URLReader {
	u.trans = &http.Transport{}
	u.trans.Proxy = http.ProxyURL(fixedURL)
	return u
}

// ReturnStatus sets the expected return status, in case a successful return status
// is not HTTP status-code OK.
func (u *URLReader) ReturnStatus(status int) *URLReader {
	u.returnStatus = status
	return u
}

// Open returns the body of the response as an io.ReadCloser, or nil and
// an error if the request fails or the return HTTP status-code
// is unequal to the expected status.
//
// The io.ReadCloser *must* be closed after the data in body has been consumed, in order
// to prevent leaks and/or preserve computer-resources.
//
// The context controls the entire lifetime of a request and its response: obtaining a
// connection, sending the request, and reading the response headers and body.
// As such, the context should have a time-limit to prevent the request from blocking too long.
func (u *URLReader) Open(ctx context.Context) (io.ReadCloser, error) {
	c := http.Client{}
	if u.trans != nil {
		c.Transport = u.trans
	}
	resp, err := c.Do(u.req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != u.returnStatus {
		lr := io.LimitReader(resp.Body, 8192)
		body, _ := ioutil.ReadAll(lr)
		defer func() {
			// In case of any remaining data
			_, _ = io.Copy(ioutil.Discard, resp.Body)
			_ = resp.Body.Close()
		}()
		return nil, fmt.Errorf("%s returned status %d; error = %s", u.loc, resp.StatusCode, body)
	}
	return resp.Body, nil
}
