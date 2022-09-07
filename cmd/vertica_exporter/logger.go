package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	// "path/filepath"
	// "time"

	"io/ioutil"

	"gopkg.in/yaml.v2"


	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger(configFile string) {
	yfile, err1 := ioutil.ReadFile("examples/vertica_exporter.yml")
	if err1 != nil {
		log.Fatal(err1)
	}
	data := make(map[string]interface{})
	err := yaml.Unmarshal(yfile, &data)
	if err != nil {
		log.Fatal(err)
	}
	var Rdays int64
	for key, value := range data {
		if key == "Retention" {
			Retention:=fmt.Sprint(value)
			re := regexp.MustCompile(`^[0-9]+`)
			// re := regexp.MustCompile("[0-9]+")
			// s := re.FindAllString(Retention,-1)
			// t := strconv.Atoi(s)
			
			// fmt.Println(t)
			// Rdays,_ = 
			
			Rdays,_=strconv.ParseInt(re.FindString(Retention),0,64)
			
		}
	}
	
	retention_days := int(Rdays)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./Logfile/logfile.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     retention_days, //days
		Compress:   true, // disabled by default
	  }
	  mWriter := io.MultiWriter(os.Stderr, lumberjackLogger)
	  log.SetOutput(mWriter)
}