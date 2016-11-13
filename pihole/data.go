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

// This the package for the Pi HOLE API
// See: https://github.com/pi-hole/AdminLTE

package pihole

type DomainsOverTime struct {
	Stats map[string]string
}

type AdsOverTime struct {
	Stats map[string]string
}

type TopQueries struct {
	Stats map[string]string
}

type TopAds struct {
	Stats map[string]string
}

type TopSources struct {
	Stats map[string]string
}

type Metrics struct {
	DomainsBeingBlocked string          `json:"domains_being_blocked"`
	DNSQueriesToday     string          `json:"dns_queries_today"`
	AdsBlockedToday     string          `json:"ads_blocked_today"`
	AdsPercentageToday  string          `json:"ads_percentage_today"`
	DomainsOverTime     DomainsOverTime `json:"domains_over_time"`
	AdsOverTime         AdsOverTime     `json:"ads_over_time"`
	TopQueries          TopQueries      `json:"top_queries"`
	TopAds              TopAds          `json:"top_ads"`
	QueryA              string          `json:"query[A]"`
	QueryAAAA           string          `json:"query[AAAA]"`
	QueryPTR            string          `json:"query[PTR]"`
	QuerySOA            string          `json:"query[SOA]"`
	Eight844            string          `json:"8.8.4.4"`
	Eight888            string          `json:"8.8.8.8"`
	TopSources          TopSources      `json:"top_sources"`
}
