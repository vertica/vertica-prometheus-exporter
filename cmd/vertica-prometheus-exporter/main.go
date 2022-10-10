package main

// (c) Copyright [2018-2022] Micro Focus or one of its affiliates.
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
// http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// 
// MIT license brought forward from the sql-exporter repo by burningalchemist
// 
// MIT License
// 
// Copyright (c) 2017 Alin Sinpalean
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"

	_ "net/http/pprof"

	_ "github.com/kardianos/minwinsvc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	vertica_prometheus_exporter "github.com/vertica/vertica-prometheus-exporter"

	log "github.com/sirupsen/logrus"
	// "gopkg.in/natefinch/lumberjack.v2"
)

const (
	envConfigFile = "VERTICAEXPORTER_CONFIG"
	envDebug      = "VERTICAEXPORTER_DEBUG"
)

var (
	showVersion   = flag.Bool("version", false, "Print version, license, and copyright information")
	listenAddress = flag.String("web.listen-address", ":9968", "Address to listen on for web interface and telemetry")
	metricsPath   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics")
	enableReload  = flag.Bool("web.enable-reload", false, "Enable reload collector data handler")
	configFile    = flag.String("config.file", "metrics/vertica-prometheus-exporter.yml", "Vertica Prometheus Exporter configuration filename")
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		// DisableColors: true,
		FullTimestamp: true,
	})

	// log.SetLevel(log.InfoLevel)
	prometheus.MustRegister(version.NewCollector("vertica_prometheus_exporter"))
}

func main() {
	
	if os.Getenv(envDebug) != "" {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(1)
	}

	// Override --alsologtostderr default value.
	if alsoLogToStderr := flag.Lookup("alsologtostderr"); alsoLogToStderr != nil {
		alsoLogToStderr.DefValue = "true"
		_ = alsoLogToStderr.Value.Set("true")
	}
	// Override the config.file default with the verticaEXPORTER_CONFIG environment variable if set.
	if val, ok := os.LookupEnv(envConfigFile); ok {
		*configFile = val
	}

	flag.Parse()

	if *showVersion {
		fmt.Println(version.Print("vertica-prometheus-exporter, Licensed under the Apache License, Version 2.0, Copyright [2018-2022] Micro Focus or one of its affiliates"))
		os.Exit(0)
	}

	log.Infof("Starting vertica exporter %s %s", version.Info(), version.BuildContext())

	exporter, err := vertica_prometheus_exporter.NewExporter(*configFile)
	if err != nil {
		log.Fatalf("Error creating exporter: %s", err)
	}
	SetupLogger(*configFile)
	// Setup and start webserver.
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "OK", http.StatusOK) })
	http.HandleFunc("/", HomeHandlerFunc(*metricsPath))
	http.HandleFunc("/config", ConfigHandlerFunc(*metricsPath, exporter))
	http.Handle(*metricsPath, promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, ExporterHandlerFor(exporter)))
	// Expose exporter metrics separately, for debugging purposes.
	http.Handle("/vertica-prometheus-exporter-metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

	// Expose refresh handler to reload query collections
	if *enableReload {
		http.HandleFunc("/reload", reloadCollectors(exporter))
	}
	log.Infof("Listening on %s", *listenAddress)
	//SetupLogger(*configFile)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
	// SetupLogger()

}

func reloadCollectors(e vertica_prometheus_exporter.Exporter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Reloading the collectors...")
		config := e.Config()
		if err := config.ReloadCollectorFiles(); err != nil {
			log.Errorf("Error reloading collector configs - %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		target, err := vertica_prometheus_exporter.NewTarget("", "", string(config.Target.DSN), config.Collectors, nil, config.Globals)
		if err != nil {
			log.Errorf("Error creating a new target - %v", err)
		}
		e.UpdateTarget([]vertica_prometheus_exporter.Target{target})

		log.Infof("Query collectors have been successfully reloaded")
		w.WriteHeader(http.StatusNoContent)
	}
}

// LogFunc is an adapter to allow the use of any function as a promhttp.Logger. If f is a function, LogFunc(f) is a
// promhttp.Logger that calls f.
type LogFunc func(args ...interface{})

// Println implements promhttp.Logger.
func (log LogFunc) Println(args ...interface{}) {
	log(args)
}
