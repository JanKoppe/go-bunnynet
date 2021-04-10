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

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

func NewClient(key string) (Client, error) {
	c := Client{}
	h := &http.Client{Timeout: 60 * time.Second}

	if key == "" {
		env := os.Getenv("BUNNYCDN_ACCESSKEY")
		if env == "" {
			return c, errors.New("required access key not provided")
		} else {
			key = env
		}
	}

	baseurl := "https://api.bunny.net/"

	if envurl := os.Getenv("BUNNYCDN_URL"); envurl != "" {
		baseurl = envurl
	}

	u, err := url.Parse(baseurl)
	if err != nil {
		return c, err
	}

	c.AccessKey = key
	c.BaseURL = u
	c.httpClient = h

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode == 400 {
		msg := ErrorResponse{}

		if len(body) < 1 {
			return resp, errors.New("client error")
		}

		if err = json.Unmarshal(body, &msg); err != nil {
			return resp, err
		}
		return resp, &msg
	}

	if resp.StatusCode == 401 {
		msg := ErrorResponse{}

		if len(body) < 1 {
			return resp, errors.New("unauthorized")
		}

		if err = json.Unmarshal(body, &msg); err != nil {
			return resp, err
		}
		return resp, &msg
	}

	if resp.StatusCode == 404 {
		msg := ErrorResponse{}

		if len(body) < 1 {
			return resp, errors.New("not found")
		}

		if err = json.Unmarshal(body, &msg); err != nil {
			return resp, err
		}
		return resp, &msg
	}

	if resp.StatusCode == 500 {
		msg := ErrorResponse{}

		if len(body) < 1 {
			return resp, errors.New("server error")
		}

		if err = json.Unmarshal(body, &msg); err != nil {
			return resp, err
		}
		return resp, &msg
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {

		// TODO workaround for 200 responses without content
		// return early without trying to unmarshal
		if len(body) < 1 {
			return resp, nil
		}

		if err = json.Unmarshal(body, v); err != nil {
			return resp, err
		}
	}

	return resp, err
}
