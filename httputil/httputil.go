// Package httputil implements some http utility functions.
package httputil

import "crypto/tls"
import "io/ioutil"
import "net/http"
import "strings"
import "bytes"
import "exp/html"

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

// GetNode issues a GET request, parses it and returns a *html.Node.
func GetNode(rawUrl string) (htmlNode *html.Node, err error) {

   // Gets the content from the HTTP url
   htmlBuf, err := GetRaw(rawUrl)
   if err != nil {
      return nil, err
   }

   // Make the content into a *html.Node
   htmlNode, err = html.Parse(bytes.NewReader(htmlBuf))
   if err != nil {
      return nil, err
   }
   return htmlNode, nil
}
