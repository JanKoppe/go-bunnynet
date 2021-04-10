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
	"testing"
)

func TestPullZoneEdgeRules(t *testing.T) {
	c, err := NewClient("")

	if err != nil {
		t.Errorf(err.Error())
	}

	// Create a fresh PullZone
	pullZone, err := c.CreatePullZone("go-bunnynet-testedgerules", "https://bunny.net", 0, PZTPremium)
	if err != nil {
		t.Errorf(err.Error())
	}

	rule := EdgeRule{
		ActionType: ERATForceSSL,
		Triggers: []EdgeRuleTrigger{
			{
				Type: ERTTRandomChance,
				PatternMatches: []string{
					"50",
				},
				PatternMatchingType: ERTPMTMatchAny,
			},
		},
	}

	// Insert a new EdgeRule
	ruleID, err := c.UpsertEdgeRule(pullZone.ID, rule)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Delete the EdgeRule again
	err = c.DeleteEdgeRule(pullZone.ID, ruleID)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Refresh the PullZone details
	updatedZone, err := c.GetPullZone(pullZone.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

	// We now should have no more PullZones, as we created and deleted only one.
	for _, r := range updatedZone.EdgeRules {
		t.Errorf("rule %v should have been deleted", r.Guid)
	}

	// Cleanup
	err = c.DeletePullZone(pullZone.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

}
