package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"path"
)

type CachedHttpClient struct {
	Client *http.Client
	CacheDir string
}

func NewCachedHttpClient(cacheDir string) (*CachedHttpClient, error) {
	if cacheDir != "" {
		err := DirExists(cacheDir, true)
		if err != nil {
			return nil, err
		}
	}
	c := &CachedHttpClient {
		Client: &http.Client{},
		CacheDir: cacheDir,
	}
	return c, nil
}

func (c *CachedHttpClient) GetPageRaw(url string, headers map[string] string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *CachedHttpClient) GetPage(url string, headers map[string] string) (string, error) {
	id := ""
	if c.CacheDir != "" {
		h := md5.New()
		h.Write([]byte(url))
		id = hex.EncodeToString(h.Sum(nil))

		page, err := ioutil.ReadFile(path.Join(c.CacheDir, id))
		if err == nil {
			return string(page), nil
		}
	}

	res, err := c.GetPageRaw(url, headers)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if c.CacheDir != "" {
		_ = ioutil.WriteFile(path.Join(c.CacheDir, id), data, 0644)
	}

	return string(data), nil
}
