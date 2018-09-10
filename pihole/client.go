// Copyright (C) 2016-2018 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pihole

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/nlamirault/pihole_exporter/version"
)

const (
	acceptHeader = "application/json"
	mediaType    = "application/json"
)

var (
	userAgent = fmt.Sprintf("pihole-exporter/%s", version.Version)
)

type Client struct {
	Endpoint string
	// Username string
	// Password string
}

func NewClient(endpoint string) (*Client, error) {
	url, err := url.Parse(endpoint)
	if err != nil || url.Scheme != "http" {
		return nil, fmt.Errorf("Invalid PiHole address: %s", err)
	}
	log.Debugf("PiHole client creation")
	return &Client{
		Endpoint: url.String(),
		// Username: username,
		// Password: password,
	}, nil
}

func (c *Client) setupHeaders(request *http.Request) {
	request.Header.Add("Content-Type", mediaType)
	request.Header.Add("Accept", acceptHeader)
	request.Header.Add("User-Agent", userAgent)
	// request.SetBasicAuth(c.Username, c.Password)
}

func (client *Client) GetMetrics() (*Metrics, error) {
	log.Infof("Get metrics")
	resp, err := http.Get(fmt.Sprintf("%s/admin/api.php?summaryRaw&overTimeData&topItems&recentItems&getQueryTypes&getForwardDestinations&getQuerySources&jsonForceObject", client.Endpoint))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Debugf("Metrics response: %s", body)
	var metrics Metrics
	dec := json.NewDecoder(bytes.NewBuffer(body))
	if err := dec.Decode(&metrics); err != nil {
		return nil, err
	}
	return &metrics, nil
}
