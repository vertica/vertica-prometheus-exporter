package Test

import (
	"strings"
	"io/ioutil"
    "gopkg.in/yaml.v2"
	"fmt"
	// "container/list"
    "testing"
	"reflect"
)

// type global struct {
// 	scrape_timeout_offset int
// 	min_interval int
// 	max_connections int
// 	max_idle_connections int
// 	max_connection_lifetime int
// }

// type target struct{
// 	data_source_name string
// 	collectors *list.List
	
// }
// var collector_files string

func TestDemo(t *testing.T){
	yfile, err1 := ioutil.ReadFile("vertica_exporter.yml")
	
	if err1 != nil {
			fmt.Println(fmt.Errorf("read: %w", err1))
		}
	data := make(map[string]interface{})
	err := yaml.Unmarshal(yfile, &data)
	if err != nil {
			fmt.Println(fmt.Errorf("read: %w", err))
		}
	var a,b1,b2,b3,b4,b5,c,d bool
	
	for key, value := range data {
		if key == "collector_files" && reflect.TypeOf(value).Kind() == reflect.Slice{
			a = true
		}else if key == "global" {
			for key2, value2 := range value.(map[interface {}]interface{}) {
				if (key2 == "scrape_timeout_offset" && reflect.TypeOf(value2).Kind() == reflect.String){
					b1 = true
				}else if (key2 =="min_interval" && reflect.TypeOf(value2).Kind() == reflect.String){
					b2 = true
				}else if key2 =="max_connection_lifetime" && reflect.TypeOf(value2).Kind() == reflect.String{
					b3 = true
				}else if key2 =="max_connections" && reflect.TypeOf(value2).Kind() == reflect.Int{
					b4 = true
				}else if key2 =="max_idle_connections" && reflect.TypeOf(value2).Kind() == reflect.Int {
					b5 = true
				}
			}
		}else if key =="target"{
			for keys, values := range value.(map[interface {}]interface{}) {
				if keys == "collectors" && reflect.TypeOf(values).Kind() == reflect.Slice{
					c = true
				}else if keys == "data_source_name" && reflect.TypeOf(values).Kind() == reflect.String { 
					s1:=values.(string)
					if strings.HasPrefix(s1,"vertica://"){
						d = true
					}
				}
			}
		}
	}
	if a&&b1&&b2&&b3&&b4&&b5&&c&&d{
		fmt.Println("Valid format")
	}else{
		fmt.Println("Invalid format")
	}
}