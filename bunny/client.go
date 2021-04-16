// Copyright (c) 2021 Jan Koppe
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package bunny

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	BaseURL    *url.URL
	AccessKey  string
	httpClient *http.Client
}

type ErrorResponse struct {
	ErrorKey string
	Field    string
	Message  string
	Err      error
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v, ErrorKey: %v, Field: %v", r.Message, r.ErrorKey, r.Field)
}

func NewClient(key string) (*Client, error) {

	if key == "" {
		if kenv := os.Getenv("BUNNYCDN_ACCESSKEY"); kenv == "" {
			return nil, errors.New("required access key not provided")
		} else {
			key = kenv
		}
	}

	baseurl := "https://api.bunny.net/"

	if envurl := os.Getenv("BUNNYCDN_URL"); envurl != "" {
		baseurl = envurl
	}

	u, err := url.Parse(baseurl)
	if err != nil {
		return nil, err
	}

	h := &http.Client{
		Timeout: 60 * time.Second,
	}

	c := &Client{
		AccessKey:  key,
		BaseURL:    u,
		httpClient: h,
	}

	return c, nil
}

func (c *Client) newRequest(method, path string, rawquery string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	u.RawQuery = rawquery

	var buf io.ReadWriter

	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("User-Agent", "go-bunnynet/dev")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("AccessKey", c.AccessKey)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 400:
		msg := ErrorResponse{}
		if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
			switch resp.StatusCode {
			case 400:
				return resp, errors.New("client error")
			case 401:
				return resp, errors.New("unauthorized")
			case 404:
				return resp, errors.New("not found")
			case 500:
				return resp, errors.New("server error")
			}
		}
		return resp, &msg
	case resp.StatusCode == 200 || resp.StatusCode == 201:
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF { // some 201s don't return the created resource.
			return resp, nil
		} else if err != nil { // might be a malformed response or unfitting interface, pass error
			return resp, err
		}
	default:
		return resp, nil
	}
	return resp, nil
}

func (c *Client) doRequest(method, path string, rawquery string, body interface{}, v interface{}) error {
	req, err := c.newRequest(method, path, rawquery, body)
	if err != nil {
		return err
	}

	_, err = c.do(req, v)
	return err
}
