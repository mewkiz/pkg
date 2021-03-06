// Package httputil implements some http utility functions.
package httputil

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// client is the default http client used by httputil requests.
var client = http.DefaultClient

// InsecureClient is a http client which allows https connections with invalid
// certificates.
var InsecureClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// SetClient sets the default http client used by httputil requests.
func SetClient(c *http.Client) {
	client = c
}

// Post issues a POST to the specified URL.
func Post(rawURL, bodyType, data string) (buf []byte, err error) {
	resp, err := client.Post(rawURL, bodyType, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// PostString issues a POST to the specified URL and returns the response as a
// string.
func PostString(rawURL, bodyType, data string) (s string, err error) {
	buf, err := Post(rawURL, bodyType, data)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// Get issues a GET to the specified URL and returns the raw response.
func Get(rawURL string) (buf []byte, err error) {
	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// GetString issues a GET to the specified URL and returns the response as a
// string.
func GetString(rawURL string) (s string, err error) {
	buf, err := Get(rawURL)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// GetDoc issues a GET request, parses it and returns an HTML node.
func GetDoc(rawURL string) (doc *html.Node, err error) {
	buf, err := Get(rawURL)
	if err != nil {
		return nil, err
	}
	doc, err = html.Parse(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// A Session contains the cookies and User-Agent for a series of requests.
type Session struct {
	// Cookies used for each session request.
	Cookies []*http.Cookie
	// The User-Agent of each session request. The default Go http User-Agent is
	// used if UserAgent is empty.
	UserAgent string
}

// Get issues a GET to the specified URL and returns the raw response. The
// request uses the session's cookies and User-Agent.
func (sess *Session) Get(rawURL string) (buf []byte, err error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	for _, cookie := range sess.Cookies {
		req.AddCookie(cookie)
	}
	if len(sess.UserAgent) != 0 {
		req.Header.Set("User-Agent", sess.UserAgent)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// GetString issues a GET to the specified URL and returns the response as a
// string. The request uses the session's cookies and User-Agent.
func (sess *Session) GetString(rawURL string) (s string, err error) {
	buf, err := sess.Get(rawURL)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// Post issues a POST to the specified URL. The request uses the session's
// cookies and User-Agent.
func (sess *Session) Post(rawURL, bodyType, data string) (buf []byte, err error) {
	req, err := http.NewRequest("POST", rawURL, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	for _, cookie := range sess.Cookies {
		req.AddCookie(cookie)
	}
	if len(sess.UserAgent) != 0 {
		req.Header.Set("User-Agent", sess.UserAgent)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// PostString issues a POST to the specified URL and returns the response as a
// string. The request uses the session's cookies and User-Agent.
func (sess *Session) PostString(rawURL, bodyType, data string) (s string, err error) {
	buf, err := sess.Post(rawURL, bodyType, data)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
