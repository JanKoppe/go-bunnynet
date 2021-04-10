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
	"strconv"
)

type StorageZone struct {
	ID                 int64  `json:"Id"`
	UserID             string `json:"UserId"`
	Name               string
	Password           string
	DateModified       BunnyTime
	Deleted            bool
	StorageUsed        int64
	FilesStored        int64
	Region             string
	ReplicationRegions []string
	PullZones          []PullZone
	ReadOnlyPassword   string
}

func (c *Client) ListStorageZones() ([]StorageZone, error) {
	req, err := c.newRequest("GET", "/storagezone", "", nil)
	if err != nil {
		return []StorageZone{}, err
	}

	var storageZones []StorageZone
	_, err = c.do(req, &storageZones)
	return storageZones, err
}

func (c *Client) GetStorageZone(zoneID int64) (StorageZone, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/storagezone/%v", zoneID), "", nil)
	if err != nil {
		return StorageZone{}, err
	}

	var storageZone StorageZone
	_, err = c.do(req, &storageZone)
	return storageZone, err
}

func (c *Client) AddStorageZone(originURL string, name string, region string, replicationRegions []string) (StorageZone, error) {
	opts := map[string]interface{}{
		"OriginUrl":          originURL,
		"Name":               name,
		"Region":             region,
		"ReplicationRegions": replicationRegions,
	}

	req, err := c.newRequest("POST", "/storagezone", "", opts)
	if err != nil {
		return StorageZone{}, err
	}

	var storageZone StorageZone
	_, err = c.do(req, &storageZone)
	return storageZone, err
}

func (c *Client) UpdateStorageZone(zoneID int64, originURL string, replicationRegions []string) error {
	opts := map[string]interface{}{
		"OriginUrl":        originURL,
		"ReplicationZones": replicationRegions, // sic
	}

	req, err := c.newRequest("POST", fmt.Sprintf("/storagezone/%v", zoneID), "", opts)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) DeleteStorageZone(zoneID int64) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/storagezone/%v", zoneID), "", nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) ResetStorageZonePassword(zoneID int64) error {
	v := url.Values{}
	v.Set("id", strconv.FormatInt(zoneID, 10))

	req, err := c.newRequest("POST", "/storagezone/resetPassword", v.Encode(), nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) ResetStorageZoneReadOnlyPassword(zoneID int64) error {
	v := url.Values{}
	v.Set("id", strconv.FormatInt(zoneID, 10))

	req, err := c.newRequest("POST", "/storagezone/resetReadOnlyPassword", v.Encode(), nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}
