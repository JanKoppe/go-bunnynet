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

import "net/url"

type BillingRecordType int

const (
	BRTPayment       BillingRecordType = 0
	BRTUnknown1      BillingRecordType = 1
	BRTUnknown2      BillingRecordType = 2
	BRTCharge        BillingRecordType = 3
	BRTUnknown4      BillingRecordType = 4
	BRTCouponApplied BillingRecordType = 5
)

type BillingRecord struct {
	ID               int64  `json:"Id"`
	PaymentID        string `json:"PaymentId"`
	Amount           float32
	Payer            string
	Timestamp        BunnyTime
	InvoiceAvailable bool
	Type             BillingRecordType
}

type BillingDetails struct {
	Balance                   float32
	ThisMonthCharges          float32
	BillingRecords            []BillingRecord
	MonthlyChargesStorage     float32
	MonthlyChargesEUTraffic   float32
	MonthlyChargesUSTraffic   float32
	MonthlyChargesASIATraffic float32
	MonthlyChargesSATraffic   float32
}

type BillingSummaryReport struct {
	PullZoneID           int64 `json:"PullZoneId"`
	MonthlyUsage         float32
	MonthlyBandwidthUsed int64
}

type BillingSummary struct {
	PullZones []BillingSummaryReport
}

func (c *Client) GetBillingDetails() (BillingDetails, error) {

	req, err := c.newRequest("GET", "/billing", "", nil)
	if err != nil {
		return BillingDetails{}, err
	}

	var details BillingDetails
	_, err = c.do(req, &details)
	return details, err
}

func (c *Client) GetBillingSummary() (BillingSummary, error) {

	req, err := c.newRequest("GET", "/billing/summary", "", nil)
	if err != nil {
		return BillingSummary{}, err
	}

	var summary BillingSummary
	_, err = c.do(req, &summary)
	return summary, err
}

func (c *Client) ApplyPromoCode(code string) (ErrorResponse, error) {
	// why is this a GET, bunny?
	// why is a errorresponse returned for a 200?
	v := url.Values{}
	v.Set("CouponCode", code)

	req, err := c.newRequest("GET", "/billing/applycode", v.Encode(), nil)
	if err != nil {
		return ErrorResponse{}, err
	}

	var msg ErrorResponse
	_, err = c.do(req, &msg)
	return msg, err
}
