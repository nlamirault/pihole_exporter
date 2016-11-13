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

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/nlamirault/pihole_exporter/pihole"
	exporter_version "github.com/nlamirault/pihole_exporter/version"
)

const (
	banner = "pihole_exporter - %s\n"

	namespace = "pihole"
)

var (
	debug         bool
	version       bool
	listenAddress string
	metricsPath   string
	endpoint      string
	username      string
	password      string
	ids           string

	domainsBeingBlocked = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "domains_being_blocked"),
		"Domains being blocked.",
		nil, nil,
	)
	dnsQueries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "dns_queries_today"),
		"DNS Queries today.",
		nil, nil,
	)
	adsBlocked = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ads_blocked_today"),
		"Ads blocked today.",
		nil, nil,
	)

	adsPercentage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ads_percentage_today"),
		"Ads percentage today.",
		nil, nil,
	)

	domainsOverTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "domains_over_time"),
		"Domains over time.",
		nil, nil,
	)

	adsOverTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ads_over_time"),
		"Ads over time.",
		nil, nil,
	)

	topQueries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_queries"),
		"Top queries.",
		nil, nil,
	)

	topAds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_ads"),
		"Top Ads.",
		nil, nil,
	)
	topSources = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_sources"),
		"Top sources.",
		nil, nil,
	)
	queryTypes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "query_types"),
		"DNS Query types.",
		[]string{"type"}, nil,
	)
)

// Exporter collects Pihole stats from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
	Pihole *pihole.Client
}

// NewExporter returns an initialized Exporter.
func NewExporter(endpoint string) (*Exporter, error) {
	log.Infoln("Setup Pihole exporter using URL: %s", endpoint)
	pihole, err := pihole.NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	return &Exporter{
		Pihole: pihole,
	}, nil
}

// Describe describes all the metrics ever exported by the Pihole exporter.
// It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- domainsBeingBlocked
	ch <- dnsQueries
	ch <- adsBlocked
	ch <- adsPercentage
	ch <- domainsOverTime
	ch <- adsOverTime
	ch <- topQueries
	ch <- topAds
	ch <- topSources
	ch <- queryTypes
}

// Collect the stats from channel and delivers them as Prometheus metrics.
// It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Infof("Pihole exporter starting")
	resp, err := e.Pihole.GetMetrics()
	if err != nil {
		log.Errorf("Pihole error: %s", err.Error())
		return
	}
	if val, err := strconv.ParseFloat(resp.DomainsBeingBlocked, 64); err == nil {
		ch <- prometheus.MustNewConstMetric(
			domainsBeingBlocked, prometheus.GaugeValue, val)
	}
	if val, err := strconv.ParseFloat(resp.DNSQueriesToday, 64); err == nil {
		ch <- prometheus.MustNewConstMetric(
			dnsQueries, prometheus.GaugeValue, val)
	}
	if val, err := strconv.ParseFloat(resp.AdsBlockedToday, 64); err == nil {
		ch <- prometheus.MustNewConstMetric(
			adsBlocked, prometheus.GaugeValue, val)
	}
	if val, err := strconv.ParseFloat(resp.AdsPercentageToday, 64); err == nil {
		ch <- prometheus.MustNewConstMetric(
			adsPercentage, prometheus.GaugeValue, val)
	}
	if val, err := strconv.ParseFloat(resp.QueryA, 64); err == nil {
		ch <- prometheus.MustNewConstMetric(
			queryTypes, prometheus.GaugeValue, val, "A")
	}
	if val, err := strconv.ParseFloat(resp.QueryAAAA, 64); err == nil {
		ch <- prometheus.MustNewConstMetric(
			queryTypes, prometheus.GaugeValue, val, "AAAA")
	}
	if val, err := strconv.ParseFloat(resp.QueryPTR, 64); err == nil {
		ch <- prometheus.MustNewConstMetric(
			queryTypes, prometheus.GaugeValue, val, "PTR")
	}
	if val, err := strconv.ParseFloat(resp.QuerySOA, 64); err == nil {
		ch <- prometheus.MustNewConstMetric(
			queryTypes, prometheus.GaugeValue, val, "SOA")
	}

	log.Infof("Pihole exporter finished")
}

func init() {
	// parse flags
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.StringVar(&listenAddress, "web.listen-address", ":9311", "Address to listen on for web interface and telemetry.")
	flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	flag.StringVar(&endpoint, "pihole", "", "Endpoint of Pihole")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(banner, exporter_version.Version))
		flag.PrintDefaults()
	}

	flag.Parse()

	if version {
		fmt.Printf("%s", exporter_version.Version)
		os.Exit(0)
	}

	if len(endpoint) == 0 {
		usageAndExit("Pihole endpoint cannot be empty.", 1)
	}
}

func main() {
	exporter, err := NewExporter(endpoint)
	if err != nil {
		log.Errorf("Can't create exporter : %s", err)
		os.Exit(1)
	}
	log.Infoln("Register exporter")
	prometheus.MustRegister(exporter)

	http.Handle(metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Pihole Exporter</title></head>
             <body>
             <h1>Pihole Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.Infoln("Listening on", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Usage()
	os.Exit(exitCode)
}
