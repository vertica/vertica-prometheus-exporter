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

import (
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

// Logging function to log to the file
// it will take max file size and retention_day from vertica-prometheus-exporter.yml file .

func SetupLogger(configFile string) {
	yfile, err1 := ioutil.ReadFile("metrics/vertica-prometheus-exporter.yml")
	if err1 != nil {
		log.Fatal(err1)
	}
	data := make(map[interface{}]interface{})

	err := yaml.Unmarshal(yfile, &data)
	if err != nil {
		log.Fatal(err)
	}

	var Rdays int
	var maxFilesize int
	for key, value := range data {
		if key == "Log" {
			for key, v := range value.(map[interface{}]interface{}) {
				if key == "retention_day" {
					Rdays = v.(int)
				}
				if key == "max_log_filesize" {
					maxFilesize = v.(int)
				}
			}

		}
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logfile/vertica-prometheus-exporter.log",
		MaxSize:    maxFilesize, // megabytes
		MaxBackups: 3,
		MaxAge:     Rdays, //days
		Compress:   true,  // disabled by default
	}
	mWriter := io.MultiWriter(os.Stderr, lumberjackLogger)
	log.SetOutput(mWriter)

}
