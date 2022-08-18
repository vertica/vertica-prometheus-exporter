package Test

import (
	"strings"
	"fmt"	
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// fmt.Println(filepath.Abs("./examples/*.collector.yml"))
// txt, _ := ioutil.ReadFile(path)
func WalkMatch(root, pattern string) ([]string, error) {
    var matches []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
            return err
        } else if matched {
            matches = append(matches, path)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return matches, nil
}

// TestExporter function is checking vertica_exporter.yml file
func TestExporter(t *testing.T) {
	yfile,err1:= ioutil.ReadFile("../examples/vertica_exporter.yml")

	if err1 != nil {
		fmt.Println(fmt.Errorf("read: %w", err1))
	}
	data := make(map[string]interface{})
	err := yaml.Unmarshal(yfile, &data)
	if err != nil {
		fmt.Println(fmt.Errorf("read: %w", err))
	}
	//var a, b1, b2, b3, b4, b5, c, d bool
	var cp1, cp2, cp3, cp4, cp5, cp6, cp7,cp8 bool

	for key, value := range data {
		if key == "collector_files" && reflect.TypeOf(value).Kind() == reflect.Slice {
			cp1 = true
		} else if key == "global" {
			for key2, value2 := range value.(map[interface{}]interface{}) {
				if key2 == "scrape_timeout_offset" && reflect.TypeOf(value2).Kind() == reflect.String {
					cp2 = true

				} else if key2 == "min_interval" && reflect.TypeOf(value2).Kind() == reflect.String {
					cp3 = true

				} else if key2 == "max_connection_lifetime" && reflect.TypeOf(value2).Kind() == reflect.String {
					cp4 = true

				} else if key2 == "max_connections" && reflect.TypeOf(value2).Kind() == reflect.Int {
					cp5 = true

				} else if key2 == "max_idle_connections" && reflect.TypeOf(value2).Kind() == reflect.Int {
					cp6 = true
				}
			}
		} else if key == "target" {
			for keys, values := range value.(map[interface{}]interface{}) {
				if keys == "collectors" && reflect.TypeOf(values).Kind() == reflect.Slice {
					cp7 = true
				} else if keys == "data_source_name" && reflect.TypeOf(values).Kind() == reflect.String {
					s1 := values.(string)
					if strings.HasPrefix(s1, "vertica://") {
						cp8 = true
					}
				}
			}
		}
	}
	cps := []bool{cp1, cp2, cp3, cp4, cp5, cp6, cp7,cp8}
	for _, cp := range cps {
		switch cp {
		case cp1:
			if !cp1 {
				fmt.Println("collector_files not configured properly")
				t.Fail()
			}

		case cp2:
			if !cp2 {
				fmt.Println("global:scrape_timeout_offset not configured properly")
				t.Fail()
			}
		case cp3:
			if !cp3 {
				fmt.Println("global:min_interval not configured properly")
				t.Fail()
			}
		case cp4:
			if !cp4 {
				fmt.Println("global:max_connection_lifetime not configured properly")
				t.Fail()
			}
		case cp5:
			if !cp5 {
				fmt.Println("global:max_connections not configured properly")
				t.Fail()
			}
		case cp6:
			if !cp6 {
				fmt.Println("global:max_idle_connections not configured properly")
				t.Fail()
			}
		case cp7:
			if !cp7 {
				fmt.Println("target:collectors not configured properly")
				t.Fail()
			}
		case cp8:
			if !cp8 {
				fmt.Println("target:data_source_name not configured properly")
				t.Fail()
			}
		}
	}

}





func TestSamp(t *testing.T){
	var path = "../examples/"
	files, err := WalkMatch(path,"*.collector.yml")

	sfile,err2 := ioutil.ReadFile("sample.yml")
	if err2 != nil {
		fmt.Println(fmt.Errorf("read: %w", err2))
	}
	sample := make(map[string]interface{})
	errs := yaml.Unmarshal(sfile, &sample)
	if errs != nil {
		fmt.Println(fmt.Errorf("read: %w", errs))
	}
	var Keywords []string
	Types := make(map[string]interface{})
	for std,stdVals := range sample{
		if std == "keywords"{
			for _,item := range stdVals.([]interface{}) {
				Keywords = append(Keywords,item.(string))
			}
			
		}else if std == "types"{
			for item,val := range stdVals.(map[interface{}]interface{}) {
				Types[item.(string)] = val
			} 	
		}
	}
	for _, file := range files {
		yfile, err1 := ioutil.ReadFile(file)
		
		if err1 != nil {
			fmt.Println(fmt.Errorf("read: %w", err1))
			// t.Fail()
		}
		data := make(map[string]interface{})
		erry := yaml.Unmarshal(yfile, &data)
		
		if err != nil {
			fmt.Println(fmt.Errorf("read: %w", err))
		}
		if erry != nil {
			fmt.Println(fmt.Errorf("read: %w", erry))
		}
		
		
		for key, value := range data {
			count := 0
			for i := range Keywords{
				count=count+1
				if Keywords[i] == key{
					count=0
					if _,ok := Types[key];ok{
						if Types[key]=="Slice"{
							// fmt.Println(value)
							if reflect.TypeOf(value).Kind() != reflect.Slice{
								fmt.Println(key, "not configured properly")
								t.Fail()
								break
							}
							
							for _, mvalue := range value.([]interface{}) {
								for m_key, m_value := range mvalue.(map[interface{}]interface{}) {
									// fmt.Println(m_key,":",m_value)
									count2 :=0
									for j := range Keywords{
										// fmt.Println(m_key,"===",Keywords[j])
										count2=count2+1
										if Keywords[j] == m_key{
											count2 = 0
											if m_key == "metric_name" && reflect.TypeOf(m_value).Kind() == reflect.String {
												s1 := m_value.(string)
												if !strings.HasPrefix(s1, "vertica_") {
													fmt.Println(key,":",m_key, "not configured properly")
													t.Fail()
													break
												}
																							
											}else if Types[m_key.(string)]=="Slice"{
												if reflect.TypeOf(m_value).Kind() != reflect.Slice{
											
													fmt.Println(key, "not configured properly")
													t.Fail()
													break
												}
											}else if reflect.TypeOf(m_value).Kind() != reflect.String{
												fmt.Println(key,":",m_key, "not configured properly")
												t.Fail()
												break
											}
											break

										}

										
									}
									if count2 == len(Keywords){
										fmt.Println(key,":",m_key,"key not found")
										t.Fail()
										break
									}
								}
							}
									
						}
						break
					}
					if reflect.TypeOf(value).Kind() != reflect.String{
						fmt.Println(key, "not configured properly")
						t.Fail()	
					}
					break
				}
				
			}
			if count == len(Keywords){
				fmt.Println(key,"key not found")
				t.Fail()
				break
			}
		}
	}
	
}

	

// Test_Verticastandard function is checking vertica_standard.collector.yml file
// func Test_Verticastandard(t *testing.T) {
// 	yfile, err1 := ioutil.ReadFile("vertica_standard.collector.yml")

// 	if err1 != nil {
// 		fmt.Println(fmt.Errorf("read: %w", err1))
// 	}
// 	data := make(map[string]interface{})
// 	err := yaml.Unmarshal(yfile, &data)
// 	if err != nil {
// 		fmt.Println(fmt.Errorf("read: %w", err))
// 	}
// 	var cp1, cp2, cp3, cp4, cp5, cp6 bool
// 	for key, value := range data {
// 		if key == "collector_name" && value == "vertica_standard" {
// 			cp1 = true
// 		} else if key == "metrics" {
// 			for _, mvalue := range value.([]interface{}) {
// 				for m_key, m_value := range mvalue.(map[interface{}]interface{}) {
// 					if m_key == "metric_name" && reflect.TypeOf(m_value).Kind() == reflect.String {
// 						s1 := m_value.(string)
// 						if strings.HasPrefix(s1, "vertica_") {
// 							cp2 = true
// 						}
// 					} else if m_key == "type" && reflect.TypeOf(m_value).Kind() == reflect.String {
// 						cp3 = true
// 					} else if m_key == "help" && reflect.TypeOf(m_value).Kind() == reflect.String {
// 						cp4 = true
// 					//key_labels
// 					} else if m_key == "values" && reflect.TypeOf(m_value).Kind() == reflect.Slice {
// 						cp5 = true
// 					} else if m_key == "query" || m_key == "query_ref"{
// 						cp6 = true
// 					}
// 					//&& reflect.TypeOf(m_value).Kind() == reflect.String {
						
					

// 				}
// 			}

// 		}
// 	}

// 	cps := []bool{cp1, cp2, cp3, cp4, cp5, cp6}
// 	for _, cp := range cps {
// 		switch cp {
// 		case cp1:
// 			if !cp1 {
// 				fmt.Println("collector_name not configured properly")
// 				t.Fail()
// 			}

// 		case cp2:
// 			if !cp2 {
// 				fmt.Println("metrics:metric_name not configured properly")
// 				t.Fail()
// 			}
// 		case cp3:
// 			if !cp3 {
// 				fmt.Println("metrics:type not configured properly")
// 				t.Fail()
// 			}
// 		case cp4:
// 			if !cp4 {
// 				fmt.Println("metrics:help not configured properly")
// 				t.Fail()
// 			}
// 		case cp5:
// 			if !cp5 {
// 				fmt.Println("metrics:values not configured properly")
// 				t.Fail()
// 			}
// 		case cp6:
// 			if !cp6 {
// 				fmt.Println("metrics:query not configured properly")
// 				t.Fail()
// 			}
// 		}
// 	}
// }



