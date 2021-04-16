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
	"fmt"
)

type EdgeRuleTriggerPatternMatchingType int32

const (
	ERTPMTMatchAny   EdgeRuleTriggerPatternMatchingType = 0
	ERTPMTMatchAll   EdgeRuleTriggerPatternMatchingType = 1
	ERTPMTMatchANone EdgeRuleTriggerPatternMatchingType = 2
)

type EdgeRuleTriggerType int32

const (
	ERTTUrl            EdgeRuleTriggerType = 0
	ERTTRequestHeader  EdgeRuleTriggerType = 1
	ERTTResponseHeader EdgeRuleTriggerType = 2
	ERTTUrlExtension   EdgeRuleTriggerType = 3
	ERTTCountryCode    EdgeRuleTriggerType = 4
	ERTTRemoteIP       EdgeRuleTriggerType = 5
	ERTTUrlQueryString EdgeRuleTriggerType = 6
	ERTTRandomChance   EdgeRuleTriggerType = 7
)

type EdgeRuleTrigger struct {
	Type                EdgeRuleTriggerType
	PatternMatches      []string
	PatternMatchingType EdgeRuleTriggerPatternMatchingType
	Parameter1          string
}

type EdgeRuleActionType int32

const (
	ERATForceSSL                   EdgeRuleActionType = 0
	ERATRedirect                   EdgeRuleActionType = 1
	ERATOriginURL                  EdgeRuleActionType = 2
	ERATOverrideCacheTime          EdgeRuleActionType = 3
	ERATBlockRequest               EdgeRuleActionType = 4
	ERATSetResponseHeader          EdgeRuleActionType = 5
	ERATSetRequestHeader           EdgeRuleActionType = 6
	ERATForceDownload              EdgeRuleActionType = 7
	ERATDisableTokenAuthentication EdgeRuleActionType = 8
	ERATEnableTokenAuthentication  EdgeRuleActionType = 9
	ERATOverrideCacheTimePublic    EdgeRuleActionType = 10
	ERATIgnoreQueryString          EdgeRuleActionType = 11
	ERATDisableOptimizer           EdgeRuleActionType = 12
	ERATForceCompression           EdgeRuleActionType = 13
)

type EdgeRuleTriggerMatchingType int32

const (
	ERTMTMatchAny   EdgeRuleTriggerMatchingType = 0
	ERTMTMatchAll   EdgeRuleTriggerMatchingType = 1
	ERTMTMatchANone EdgeRuleTriggerMatchingType = 2
)

type EdgeRule struct {
	Guid                string `json:",omitempty"`
	ActionType          EdgeRuleActionType
	ActionParameter1    string
	ActionParameter2    string
	Triggers            []EdgeRuleTrigger
	TriggerMatchingType EdgeRuleTriggerMatchingType
	Description         string
	Enabled             bool
}

func (c *Client) UpsertEdgeRule(zoneID int64, r EdgeRule) (string, error) {
	// Because the bunny.net API does not reply the guid of a newly created EdgeRule,
	// we have to compare the list of EdgeRules in the PullZone before and after the
	// API call. The diff should contain exactly one new EdgeRule, from which we can
	// read the Guid.
	//
	// This really is only a workaround, because it could be that another rule is
	// inserted by a third party before we refresh the PullZone details. The set would
	// then miss both of the new EdgeRule Guids, and the behaviour is indeterminate.
	//
	// That probability in the real world is pretty low though, so, this_is_fine.jpg

	guid := ""

	if r.Guid != "" {
		guid = r.Guid
	}

	// Get PullZone details before we modify anything
	zoneBefore, err := c.GetPullZone(zoneID)
	if err != nil {
		return "", err
	}
	// save all before-Edgerule-Guids in a "set"
	set := make(map[string]bool)
	for _, r := range zoneBefore.EdgeRules {
		set[r.Guid] = true
	}

	// upsert EdgeRule
	err = c.doRequest("POST", fmt.Sprintf("/pullzone/%v/edgerules/addOrUpdate", zoneID), "", r, nil)
	if err != nil {
		return "", err
	}

	// Refresh PullZone details
	zoneAfter, err := c.GetPullZone(zoneID)
	if err != nil {
		return "", err
	}

	// if any of the after-EdgeRule-Guids is not in the "set", it's the Guid
	// of the freshly created EdgeRule. Return that!
	for _, r := range zoneAfter.EdgeRules {
		if !set[r.Guid] {
			guid = r.Guid
		}
	}

	return guid, err
}

func (c *Client) DeleteEdgeRule(zoneID int64, ruleID string) error {
	return c.doRequest("DELETE", fmt.Sprintf("/pullzone/%v/edgerules/%v", zoneID, ruleID), "", nil, nil)
}

// This API call seems to be incomplete. This function is effectively a stub, and doesn't actually do anything.
// I suggest you update the edge rule and toggle the Enabled field there.
//func (c *Client) SetEdgeRuleEnabled(zoneID int64, ruleID string, enabled bool) error {
//	opts := map[string]interface{}{
//		"Id":    zoneID,
//		"Value": enabled,
//	}
//	req, err := c.newRequest("DELETE", fmt.Sprintf("/pullzone/%v/edgerules/%v/setEdgeRuleEnabled", zoneID, ruleID), "", opts)
//
//	if err != nil {
//		return err
//	}
//
//	_, err = c.do(req, nil)
//	return err
//}
