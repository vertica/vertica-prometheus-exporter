package main

import (
	// "gopkg.in/yaml.v2"

	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/vertica/vertica-exporter/config"
	"k8s.io/klog/v2"
)

func Updatelog(configFile string) {
	// const Retention = "5days"
	c, err1 := config.Load(configFile)
	Retention := strings.Join(c.Retention, "")
	re := regexp.MustCompile(`^[0-9]+`)
	Rdays, _ := strconv.ParseInt(re.FindString(Retention), 0, 64)
	fmt.Println("inside update", reflect.TypeOf(Rdays).Kind())
	// yfile, err1 := ioutil.ReadFile("examples/vertica_exporter.yml")

	if err1 != nil {
		klog.Fatal(err1)
	}
	// data := make(map[string]interface{})
	// err := yaml.Unmarshal(yfile, &data)
	// if err != nil {
	// 	klog.Fatal(err)
	// }

	f, err := os.Open("../../LogFile/myfile.log")
	if err != nil {
		klog.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Log file created at") {
			re := regexp.MustCompile(`\d{4}/\d{2}/\d{2}`)
			date := re.FindString(scanner.Text())
			logdate, _ := time.Parse("2006/01/02", date)
			currentdate := time.Now()
			difference := currentdate.Sub(logdate)
			days := int64(difference.Hours() / 24)

			if days >= Rdays {
				fmt.Println("delete log")
			}
		}

	}
}
