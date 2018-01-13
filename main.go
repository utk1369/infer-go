package main

import (
	"fmt"

	"github.com/utk1369/infer-go/model/ruler"
)

func TestNew() {
}

func TestDefault() {
	rules := []byte(`[
		{"comparator": "eq", "path": "library.name", "value": "go-ruler"},
		{"comparator": "gte", "path": "library.age", "value": 1.2}
	  ]`)

	// supports loading rules from JSON data
	engine, _ := ruler.NewRulerWithJson(rules)

	result, _, _ := engine.Test(map[string]interface{}{
		"library": map[string]interface{}{
			"name": "go-ruler",
			"age":  0.9,
		},
	})
	fmt.Println("Result: ", result)
}

func main() {
	TestNew()
}
