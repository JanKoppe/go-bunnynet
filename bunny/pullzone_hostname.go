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

import "fmt"

type PullZoneHostname struct {
	ID               int `json:"Id,omitempty"`
	Value            string
	ForceSSL         bool
	IsSystemHostname bool
	HasCertificate   bool
}

func (c *Client) RemovePullZoneHostname(zoneID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}
	return c.doRequest("DELETE", fmt.Sprintf("/pullzone/%v/removeHostname", zoneID), "", opts, nil)
}

func (c *Client) AddPullZoneHostname(zoneID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}
	return c.doRequest("POST", fmt.Sprintf("/pullzone/%v/addHostname", zoneID), "", opts, nil)
}

func (c *Client) SetPullZoneHostnameForceSSL(zoneID int64, hostname string, forceSSL bool) error {
	opts := map[string]interface{}{
		"Hostname": hostname,
		"ForceSSL": forceSSL,
	}
	return c.doRequest("POST", fmt.Sprintf("/pullzone/%v/setForceSSL", zoneID), "", opts, nil)
}
