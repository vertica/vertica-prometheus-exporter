package main

import (
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

// Logging function to log to the file
// it will take max file size and retention_day from vertica_prometheus_exporter.yml file .

func SetupLogger(configFile string) {
	yfile, err1 := ioutil.ReadFile("metrices/vertica_prometheus_exporter.yml")
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
		Filename:   "./Logfile/logfile.log",
		MaxSize:    maxFilesize, // megabytes
		MaxBackups: 3,
		MaxAge:     Rdays, //days
		Compress:   true,  // disabled by default
	}
	mWriter := io.MultiWriter(os.Stderr, lumberjackLogger)
	log.SetOutput(mWriter)

}
