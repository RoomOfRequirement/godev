package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	client    *http.Client
	transport *http.Transport
)

func init() {
	transport = &http.Transport{MaxIdleConnsPerHost: 100}
	client = &http.Client{
		Transport:     transport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
}

// Get ...
func Get(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

// Post ...
func Post(url string, values url.Values) ([]byte, error) {
	resp, err := client.PostForm(url, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
