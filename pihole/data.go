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

// This the package for the Pi HOLE API
// See: https://github.com/pi-hole/AdminLTE

package pihole

type DomainsOverTime struct {
	Stats map[string]string
}

type AdsOverTime struct {
	Stats map[string]string
}

// Metrics define PiHome Prometheus metrics
type Metrics struct {
	DomainsBeingBlocked float64            `json:"domains_being_blocked"`
	DNSQueriesToday     float64            `json:"dns_queries_today"`
	AdsBlockedToday     float64            `json:"ads_blocked_today"`
	AdsPercentageToday  float64            `json:"ads_percentage_today"`
	DomainsOverTime     DomainsOverTime    `json:"domains_over_time"`
	AdsOverTime         AdsOverTime        `json:"ads_over_time"`
	TopQueries          map[string]float64 `json:"top_queries"`
	TopAds              map[string]float64 `json:"top_ads"`
	TopSources          map[string]float64 `json:"top_sources"`
	QueryA              float64            `json:"query[A]"`
	QueryAAAA           float64            `json:"query[AAAA]"`
	QueryPTR            float64            `json:"query[PTR]"`
	QuerySOA            float64            `json:"query[SOA]"`
	Eight844            float64            `json:"8.8.4.4"`
	Eight888            float64            `json:"8.8.8.8"`
}
