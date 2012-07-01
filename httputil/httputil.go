// Package httputil implements some http utility functions.
package httputil

import "bytes"
import "crypto/tls"
import "exp/html"
import "io/ioutil"
import "net/http"
import "strings"

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
func Post(rawUrl, bodyType, data string) (s string, err error) {
   buf, err := PostRaw(rawUrl, bodyType, data)
   if err != nil {
      return "", err
   }
   return string(buf), nil
}

// PostRaw issues a POST to the specified URL and returns the raw response.
func PostRaw(rawUrl, bodyType, data string) (buf []byte, err error) {
   res, err := client.Post(rawUrl, bodyType, strings.NewReader(data))
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf, err = ioutil.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return buf, nil
}

// Get issues a GET to the specified URL.
func Get(rawUrl string) (s string, err error) {
   buf, err := GetRaw(rawUrl)
   if err != nil {
      return "", err
   }
   return string(buf), nil
}

// GetRaw issues a GET to the specified URL and returns the raw response.
func GetRaw(rawUrl string) (buf []byte, err error) {
   res, err := client.Get(rawUrl)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf, err = ioutil.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return buf, nil
}

// GetDoc issues a GET request, parses it and returns a *html.Node.
func GetDoc(rawUrl string) (doc *html.Node, err error) {
   buf, err := GetRaw(rawUrl)
   if err != nil {
      return nil, err
   }
   doc, err = html.Parse(bytes.NewReader(buf))
   if err != nil {
      return nil, err
   }
   return doc, nil
}
