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
)

type VideoLibrary struct {
	ID                               int64 `json:"Id,omitempty"`
	Name                             string
	VideoCount                       int       `json:",omitempty"`
	DateCreated                      BunnyTime `json:",omitempty"`
	ReplicationRegions               []string
	APIKey                           string `json:"ApiKey"`
	ReadOnlyAPIKey                   string `json:"ReadOnlyApiKey"`
	HasWatermark                     bool
	WatermarkPositionLeft            int
	WatermarkPositionTop             int
	WatermarkWidth                   int
	WatermarkHeight                  int
	EnabledResolutions               string
	ViAiPublisherID                  string `json:"ViAiPublisherId"`
	VastTagURL                       string `json:"VastTagUrl"`
	WebhookURL                       string `json:"WebhookUrl"`
	CaptionsFontSize                 int
	CaptionsFontColor                string
	CaptionsBackground               string
	UILanguage                       string
	AllowEarlyPlay                   bool
	PlayerTokenAuthenticationEnabled bool
	AllowedReferrers                 []string
	BlockedReferrers                 []string
	BlockNoneReferrer                bool
	EnableMP4Fallback                bool
	KeepOriginalFiles                bool
	AllowDirectPlay                  bool
	EnableDRM                        bool
	Bitrate240p                      int
	Bitrate360p                      int
	Bitrate480p                      int
	Bitrate720p                      int
	Bitrate1080p                     int
	Bitrate1440p                     int
	Bitrate2160p                     int
}

func (c *Client) ListVideoLibraries() (*[]VideoLibrary, error) {
	var videoLibrary []VideoLibrary
	return &videoLibrary, c.doRequest("GET", "/videolibrary", "", nil, &videoLibrary)
}

func (c *Client) GetVideoLibrary(libraryID int64) (*VideoLibrary, error) {
	var videoLibrary VideoLibrary
	return &videoLibrary, c.doRequest("GET", fmt.Sprintf("/videolibrary/%v", libraryID), "", nil, &videoLibrary)
}

func (c *Client) AddVideoLibrary(name string, replicationRegions []string) (*VideoLibrary, error) {
	opts := map[string]interface{}{
		"Name":               name,
		"ReplicationRegions": replicationRegions,
	}
	var videoLibrary VideoLibrary
	return &videoLibrary, c.doRequest("POST", "/videolibrary", "", opts, &videoLibrary)
}

func (c *Client) UpdateVideoLibrary(library VideoLibrary) (*VideoLibrary, error) {
	libraryID := library.ID

	// null non-settable fields so they get omitted in marshal
	library.ID = 0
	library.VideoCount = 0
	library.DateCreated = BunnyTime{}
	var videoLibrary VideoLibrary
	return &videoLibrary, c.doRequest("POST", fmt.Sprintf("/videolibrary/%v", libraryID), "", library, &videoLibrary)
}

func (c *Client) DeleteVideoLibrary(libraryID int64) error {
	return c.doRequest("DELETE", fmt.Sprintf("/videolibrary/%v", libraryID), "", nil, nil)
}

func (c *Client) AddVideoLibraryAllowedReferrer(libraryID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}
	return c.doRequest("POST", fmt.Sprintf("/videolibrary/%v/addAllowedReferrer", libraryID), "", opts, nil)
}

func (c *Client) RemoveVideoLibraryAllowedReferrer(libraryID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}
	return c.doRequest("POST", fmt.Sprintf("/videolibrary/%v/removeAllowedReferrer", libraryID), "", opts, nil)
}

func (c *Client) AddVideoLibraryBlockedReferrer(libraryID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}
	return c.doRequest("POST", fmt.Sprintf("/videolibrary/%v/addBlockedReferrer", libraryID), "", opts, nil)

}

func (c *Client) RemoveVideoLibraryBlockedReferrer(libraryID int64, hostname string) error {
	opts := map[string]string{
		"Hostname": hostname,
	}
	return c.doRequest("POST", fmt.Sprintf("/videolibrary/%v/removeBlockedReferrer", libraryID), "", opts, nil)
}
