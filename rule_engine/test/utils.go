package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func getJsonByteArray(fileName string) []byte {

	byteArr, err := ioutil.ReadFile("./data/" + fileName)
	if err != nil {
		fmt.Println("Error loading file: ", err)
	}
	return byteArr
}

func getMap(param interface{}) map[string]interface{} {
	x, _ := json.Marshal(param)
	var o map[string]interface{}
	json.Unmarshal(x, &o)
	return o
}
