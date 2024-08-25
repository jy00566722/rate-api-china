package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
)

const (
	Reset      = "\033[0m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Purple     = "\033[35m"
	Cyan       = "\033[36m"
	Gray       = "\033[37m"
	White      = "\033[97m"
	RedBold    = "\033[31;1m"
	GreenBold  = "\033[32;1m"
	YellowBold = "\033[33;1m"
	BlueBold   = "\033[34;1m"
	PurpleBold = "\033[35;1m"
	CyanBold   = "\033[36;1m"
)

func Test() {
	dir1 := "/Users/pengfeng/vscode/rate-all-with-wxt/internation-wxt-vue-rate/assets/locales"
	dir2 := "/Users/pengfeng/vscode/rate-all-with-wxt/internation-wxt-vue-rate/assets/locales1"

	files1, err := ioutil.ReadDir(dir1)
	if err != nil {
		panic(err)
	}

	for _, f := range files1 {
		file1 := filepath.Join(dir1, f.Name())
		file2 := filepath.Join(dir2, f.Name())
		content1, err := ioutil.ReadFile(file1)
		if err != nil {
			panic(err)
		}

		content2, err := ioutil.ReadFile(file2)
		if err != nil {
			panic(err)
		}

		var data1, data2 map[string]interface{}

		if err = json.Unmarshal(content1, &data1); err != nil {
			panic(fmt.Sprintf("error decoding json from file: %s\nerror: %s", file1, err))
		}

		if err = json.Unmarshal(content2, &data2); err != nil {
			panic(fmt.Sprintf("error decoding json from file: %s\nerror: %s", file2, err))
		}

		// fmt.Printf("比较文件: %s\n", f.Name())
		compareJSON(data1, data2, f.Name())

	}
}

// compareJSON compares the JSON data map field by field and prints the differences found.
func compareJSON(json1, json2 map[string]interface{}, fileName string) {
	allKeys := make(map[string]bool)
	for k := range json1 {
		allKeys[k] = true
	}
	for k := range json2 {
		allKeys[k] = true
	}

	var hasDifference bool
	for key := range allKeys {
		v1, ok1 := json1[key]
		v2, ok2 := json2[key]

		// Check if both JSONs contain the key
		if !ok1 {
			fmt.Printf("%sKey '%s' 在文件中没有  '%s' on the first folder.%s\n", Purple, key, fileName, Reset)
			hasDifference = true
			continue
		}
		if !ok2 {
			fmt.Printf("%sKey '%s' 在文件中没有 '%s' on the second folder.%s\n", Cyan, key, fileName, Reset)
			hasDifference = true
			continue
		}

		// If both files have the key, check if the values are the same
		if !reflect.DeepEqual(v1, v2) {
			fmt.Printf("%sValue of key '%s' 【值不同】 in the file '%s'.\nFirst folder: %v\nSecond folder: %v%s\n", BlueBold, key, fileName, v1, v2, Reset)
			hasDifference = true
		}
	}

	if !hasDifference {
		fmt.Printf("%s完全相同!!!! '%s'.%s\n", YellowBold, fileName, Reset)
	}
}
