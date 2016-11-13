// Copyright (C) 2016 Nicolas Lamirault <nicolas.lamirault@gmail.com>

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

type domoticzserver struct {
	*httptest.Server
}

func handler(server *domoticzserver, uri string, filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(b)
	}
}

func newDomoticzServer(uri string, filename string) *domoticzserver {
	h := &domoticzserver{}
	h.Server = httptest.NewServer(handler(h, uri, filename))
	return h
}

func getClientAndServer(t *testing.T, uri string, username string, password string, filename string) (*domoticzserver, *Client) {
	h := newDomoticzServer(uri, filename)
	client, err := NewClient(h.Listener.Addr().String(), username, password) // h.URL, username, password)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return h, client
}

func TestDomoticzGetAllDevicesWithEmptyResult(t *testing.T) {
	server, client := getClientAndServer(t, "", "", "", "no_devices.json")
	defer server.Close()
	devices, err := client.GetAllDevices()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Infof("Devices response: %s", devices)
	if devices.Status != "OK" || devices.Title != "Devices" {
		log.Fatalf("Invalid devices response: %s", devices)
	}
}

func TestDomoticzGetDevice(t *testing.T) {
	server, client := getClientAndServer(t, "", "", "", "device.json")
	defer server.Close()
	device, err := client.GetDevice("123")
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Infof("Devices response: %s", device)
	if device.Status != "OK" || device.Title != "Devices" {
		log.Fatalf("Invalid device response: %s", device)
	}
	if len(device.Result) != 1 {
		log.Fatalf("Invalid number of device: %s", device)
	}
	if device.Result[0].TypeImg != "temperature" ||
		device.Result[0].Temp != 20.20 {
		log.Fatalf("Invalid device response: %s", device)
	}
}
