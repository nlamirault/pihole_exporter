// Copyright (C) 2016, 2017 Nicolas Lamirault <nicolas.lamirault@gmail.com>

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
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type piholeserver struct {
	*httptest.Server
}

func handler(server *piholeserver, uri string, filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(b)
	}
}

func newPiholeServer(uri string, filename string) *piholeserver {
	h := &piholeserver{}
	h.Server = httptest.NewServer(handler(h, uri, filename))
	return h
}

func getClientAndServer(t *testing.T, uri string, username string, password string, filename string) (*piholeserver, *Client) {
	h := newPiholeServer(uri, filename)
	client, err := NewClient(h.URL) //h.Listener.Addr().String()) // h.URL, username, password)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return h, client
}

func TestPiholeGetMetrics(t *testing.T) {
	server, client := getClientAndServer(t, "", "", "", "stats.json")
	defer server.Close()
	metrics, err := client.GetMetrics()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Infof("Metrics response: %s", metrics)
	if metrics.DomainsBeingBlocked != 101934 ||
		metrics.DNSQueriesToday != 2593 {
		log.Fatalf("Invalid metrics response: %s", metrics)
	}
}
