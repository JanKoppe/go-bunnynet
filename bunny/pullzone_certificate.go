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
	"fmt"
	"net/url"
)

func (c *Client) LoadFreeCertificate(hostname string) error {
	// why is this a GET, bunny?
	v := url.Values{}
	v.Set("hostname", hostname)

	req, err := c.newRequest("GET", "/pullzone/loadFreeCertificate", v.Encode(), nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) AddCustomCertificate(zoneID int64, hostname string, certificate string, key string) error {
	opts := map[string]string{
		"Hostname":       hostname,
		"Certificate":    certificate,
		"CertificateKey": key,
	}

	req, err := c.newRequest("POST", fmt.Sprintf("/pullzone/%v/addCertificate", zoneID), "", opts)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) DeleteCustomCertificate(zoneID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}

	req, err := c.newRequest("DELETE", fmt.Sprintf("/pullzone/%v/removeCertificate", zoneID), "", opts)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}
