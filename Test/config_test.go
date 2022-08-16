package Test

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	// "container/list"
	"reflect"
	"testing"
)

//TestExporter function is checking vertica_exporter.yml file
func Test_Exporter(t *testing.T) {
	yfile, err1 := ioutil.ReadFile("vertica_exporter.yml")

	if err1 != nil {
		fmt.Println(fmt.Errorf("read: %w", err1))
	}
	data := make(map[string]interface{})
	err := yaml.Unmarshal(yfile, &data)
	if err != nil {
		fmt.Println(fmt.Errorf("read: %w", err))
	}
	var a, b1, b2, b3, b4, b5, c, d bool

	for key, value := range data {
		if key == "collector_files" && reflect.TypeOf(value).Kind() == reflect.Slice {
			a = true
		} else if key == "global" {
			for key2, value2 := range value.(map[interface{}]interface{}) {
				if key2 == "scrape_timeout_offset" && reflect.TypeOf(value2).Kind() == reflect.String {
					b1 = true

				} else if key2 == "min_interval" && reflect.TypeOf(value2).Kind() == reflect.String {
					b2 = true

				} else if key2 == "max_connection_lifetime" && reflect.TypeOf(value2).Kind() == reflect.String {
					b3 = true

				} else if key2 == "max_connections" && reflect.TypeOf(value2).Kind() == reflect.Int {
					b4 = true

				} else if key2 == "max_idle_connections" && reflect.TypeOf(value2).Kind() == reflect.Int {
					b5 = true
				}
			}
		} else if key == "target" {
			for keys, values := range value.(map[interface{}]interface{}) {
				if keys == "collectors" && reflect.TypeOf(values).Kind() == reflect.Slice {
					c = true
				} else if keys == "data_source_name" && reflect.TypeOf(values).Kind() == reflect.String {
					s1 := values.(string)
					if strings.HasPrefix(s1, "vertica://") {
						d = true
					}
				}
			}
		}
	}
	cps := []bool{a, b1, b2, b3, b4, b5, c, d}
	for _, cp := range cps {
		switch cp {
		case a:
			if !a {
				fmt.Println("collector_files not configured properly")
				t.Fail()
			}

		case b1:
			if !b1 {
				fmt.Println("global:scrape_timeout_offset not configured properly")
				t.Fail()
			}
		case b2:
			if !b2 {
				fmt.Println("global:min_interval not configured properly")
				t.Fail()
			}
		case b3:
			if !b3 {
				fmt.Println("global:max_connection_lifetime not configured properly")
				t.Fail()
			}
		case b4:
			if !b4 {
				fmt.Println("global:max_connections not configured properly")
				t.Fail()
			}
		case b5:
			if !b5 {
				fmt.Println("global:max_idle_connections not configured properly")
				t.Fail()
			}
		case c:
			if !c {
				fmt.Println("target:collectors not configured properly")
				t.Fail()
			}
		case d:
			if !d {
				fmt.Println("target:data_source_name not configured properly")
				t.Fail()
			}
		}
	}

}

// Test_Verticastandard function is checking vertica_standard.collector.yml file
func Test_Verticastandard(t *testing.T) {
	yfile, err1 := ioutil.ReadFile("vertica_standard.collector.yml")

	if err1 != nil {
		fmt.Println(fmt.Errorf("read: %w", err1))
	}
	data := make(map[string]interface{})
	err := yaml.Unmarshal(yfile, &data)
	if err != nil {
		fmt.Println(fmt.Errorf("read: %w", err))
	}
	var cp1, cp2, cp3, cp4, cp5, cp6 bool
	for key, value := range data {
		if key == "collector_name" && value == "vertica_standard" {
			cp1 = true
		} else if key == "metrics" {
			for _, mvalue := range value.([]interface{}) {
				for m_key, m_value := range mvalue.(map[interface{}]interface{}) {
					if m_key == "metric_name" && reflect.TypeOf(m_value).Kind() == reflect.String {
						s1 := m_value.(string)
						if strings.HasPrefix(s1, "vertica_") {
							cp2 = true
						}
					} else if m_key == "type" && reflect.TypeOf(m_value).Kind() == reflect.String {
						cp3 = true
					} else if m_key == "help" && reflect.TypeOf(m_value).Kind() == reflect.String {
						cp4 = true
					} else if m_key == "values" && reflect.TypeOf(m_value).Kind() == reflect.Slice {
						cp5 = true
					} else if m_key == "query" && reflect.TypeOf(m_value).Kind() == reflect.String {
						cp6 = true
					}

				}
			}

		}
	}

	cps := []bool{cp1, cp2, cp3, cp4, cp5, cp6}
	for _, cp := range cps {
		switch cp {
		case cp1:
			if !cp1 {
				fmt.Println("collector_name not configured properly")
				t.Fail()
			}

		case cp2:
			if !cp2 {
				fmt.Println("metrics:metric_name not configured properly")
				t.Fail()
			}
		case cp3:
			if !cp3 {
				fmt.Println("metrics:type not configured properly")
				t.Fail()
			}
		case cp4:
			if !cp4 {
				fmt.Println("metrics:help not configured properly")
				t.Fail()
			}
		case cp5:
			if !cp5 {
				fmt.Println("metrics:values not configured properly")
				t.Fail()
			}
		case cp6:
			if !cp6 {
				fmt.Println("metrics:query not configured properly")
				t.Fail()
			}
		}
	}
}
