package main

import (
	"encoding/json"
	"fmt"
)

func printJson(data interface{}) {
	dataJson, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(dataJson))
}
