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

type PullZoneType int32

const (
	PZTPremium PullZoneType = 0
	PZTVolume  PullZoneType = 1
)

type PullZone struct {
	ID                                   int64  `json:"Id,omitempty"`
	Name                                 string `json:",omitempty"`
	OriginURL                            string `json:"OriginUrl"`
	Enabled                              bool
	Hostnames                            []PullZoneHostname
	StorageZoneID                        int64 `json:"StoragezoneId,omitempty"`
	AllowedReferrers                     []string
	BlockedReferrers                     []string
	BlockedIps                           []string
	EnableGeoZoneUS                      bool
	EnableGeoZoneEU                      bool
	EnableGeoZoneASIA                    bool
	EnableGeoZoneSA                      bool
	EnableGeoZoneAF                      bool
	ZoneSecurityEnabled                  bool
	ZoneSecurityKey                      string
	ZoneSecurityIncludeHashRemoteIP      bool
	IgnoreQueryStrings                   bool
	MonthlyBandwidthLimit                int
	MonthlyBandwidthUsed                 int     `json:",omitempty"`
	MonthlyCharges                       float32 `json:",omitempty"`
	AddHostHeader                        bool
	Type                                 PullZoneType
	AccessControlOrigionHeaderExtensions []string
	EnableAccessControlOriginHeader      bool
	DisableCookies                       bool
	BudgetRedirectedCountries            []string
	BlockedCountries                     []string
	EnableOriginShield                   bool
	CacheControlMaxAgeOverride           int
	CacheControlPublicMaxAgeOverride     int
	BurstSize                            int
	RequestLimit                         int
	BlockRootPathAccess                  bool
	BlockPostRequests                    bool
	LimitRatePerSecond                   float32
	LimitRateAfter                       float32
	ConnectionLimitPerIPCount            int
	PriceOverride                        float32
	AddCanonicalHeader                   bool
	EnableLogging                        bool
	EnableCacheSlice                     bool
	EdgeRules                            []EdgeRule
	EnableWebPVary                       bool
	EnableCountryCodeVary                bool
	EnableMobileVary                     bool
	EnableHostnameVary                   bool
	CnameDomain                          string
	AWSSigningEnabled                    bool
	AWSSigningKey                        string
	AWSSigningSecret                     string
	AWSSigningRegionName                 string
	LoggingIPAnonymizationEnabled        bool
	EnableTLS1                           bool
	EnableTLS1_1                         bool
	VerifyOriginSSL                      bool
	OriginShieldZoneCode                 string
	LogForwardingEnabled                 bool
	LogForwardingHostname                string
	LogForwardingPort                    int
	LogForwardingToken                   string
	LoggingStorageZoneID                 int64 `json:"LoggingStorageZoneId,omitempty"`
	FollowRedirects                      bool
	VideoLibraryID                       int64 `json:"VideoLibraryId,omitempty"`
}

func (c *Client) ListPullZones() ([]PullZone, error) {
	req, err := c.newRequest("GET", "/pullzone", "", nil)
	if err != nil {
		return []PullZone{}, err
	}

	var pullZones []PullZone
	_, err = c.do(req, &pullZones)
	return pullZones, err
}

func (c *Client) GetPullZone(zoneID int64) (PullZone, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/pullzone/%v", zoneID), "", nil)
	if err != nil {
		return PullZone{}, err
	}

	var pullZone PullZone
	_, err = c.do(req, &pullZone)
	return pullZone, err
}

func (c *Client) CreatePullZone(name string, origin string, storageZoneID int64, pzt PullZoneType) (PullZone, error) {
	opts := map[string]interface{}{
		"Name":          name,
		"OriginUrl":     origin,
		"StorageZoneId": storageZoneID,
		"Type":          pzt,
	}

	req, err := c.newRequest("POST", "/pullzone", "", opts)
	if err != nil {
		return PullZone{}, err
	}

	var pullZone PullZone
	_, err = c.do(req, &pullZone)

	return pullZone, err
}

func (c *Client) DeletePullZone(zoneID int64) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/pullzone/%v", zoneID), "", nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) UpdatePullZone(pz PullZone) error {
	pzID := pz.ID

	// null non-settable fields so they get omitted
	pz.ID = 0
	pz.MonthlyBandwidthUsed = 0
	pz.MonthlyCharges = 0.0

	req, err := c.newRequest("POST", fmt.Sprintf("/pullzone/%v", pzID), "", pz)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) ResetPullZoneToken(zoneID int64) error {
	req, err := c.newRequest("POST", fmt.Sprintf("/pullzone/%v/resetSecurityKey", zoneID), "", nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) AddPullZoneAllowedReferrer(zoneID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/pullzone/%v/addAllowedReferrer", zoneID), "", opts)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) RemovePullZoneAllowedReferrer(zoneID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/pullzone/%v/removeAllowedReferrer", zoneID), "", opts)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) AddPullZoneBlockedReferrer(zoneID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/pullzone/%v/addBlockedReferrer", zoneID), "", opts)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) RemovePullZoneBlockedReferrer(zoneID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/pullzone/%v/removeBlockedReferrer", zoneID), "", opts)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) AddPullZoneBlockedIP(zoneID int64, blockedIP string) error {
	opts := map[string]string{
		"BlockedIp": blockedIP,
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/pullzone/%v/addBlockedIp", zoneID), "", opts)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

func (c *Client) RemovePullZoneBlockedIP(zoneID int64, blockedIP string) error {
	opts := map[string]string{
		"BlockedIp": blockedIP,
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/pullzone/%v/removeBlockedIp", zoneID), "", opts)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}
