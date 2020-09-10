package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadJSONAsString() (feedLogsJsonString string, err error) {
	jsonFile, err := os.Open("resources/feedlogs.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	feedLogsJsonString = string(byteValue)

	return
}

func ReadJSONAsMap() (feedLogsJsonMap map[string]interface{}, err error) {
	jsonFile, err := os.Open("resources/feedlogs.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	json.Unmarshal([]byte(byteValue), &feedLogsJsonMap)

	return
}
