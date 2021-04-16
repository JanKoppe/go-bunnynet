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
	"testing"
)

func TestPullZoneHostnames(t *testing.T) {
	testHostname := "test.bunny.net"

	c, err := NewClient("")

	if err != nil {
		t.Errorf(err.Error())
	}

	// Create a fresh PullZone
	pullZone, err := c.CreatePullZone("go-bunnynet-testhostnames", "https://bunny.net", 0, PZTPremium)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c.AddPullZoneHostname(pullZone.ID, testHostname)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c.SetPullZoneHostnameForceSSL(pullZone.ID, testHostname, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Refresh the PullZone details
	updatedZone, err := c.GetPullZone(pullZone.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

	found := false
	forceSSL := false
	for _, h := range updatedZone.Hostnames {
		if h.Value == testHostname {
			found = true
			forceSSL = h.ForceSSL
		}
	}
	if !found {
		t.Error("hostname was not added, but not found in refreshed zone details")
	}
	if !forceSSL {
		t.Error("forceSSL was enabled, but not set to enabled in refreshed zone details")
	}

	// Remove hostname again
	err = c.RemovePullZoneHostname(pullZone.ID, testHostname)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Refresh the PullZone details
	updatedZone, err = c.GetPullZone(pullZone.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, h := range updatedZone.Hostnames {
		if h.Value == testHostname {
			t.Error("hostname was not removed, found in refreshed zone details")
		}
	}

	// Cleanup
	err = c.DeletePullZone(pullZone.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

}
