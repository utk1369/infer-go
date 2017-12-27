package main

import (
	"fmt"

	"io/ioutil"

	"encoding/json"

	"bitbucket.org/grabpay/go-ruler/model"
	"bitbucket.org/grabpay/go-ruler/model/ruler"
	"bitbucket.org/grabpay/go-ruler/rule_engine"
)

func getJsonByteArray() []byte {

	byteArr, err := ioutil.ReadFile("./examples/test_expression.json")
	if err != nil {
		fmt.Println("Error loading file: ", err)
	}
	return byteArr
}

func TestNew() {
	ruleSet, _ := model.LoadRuleSetFromJson(getJsonByteArray())

	fmt.Println("RuleSet: ", ruleSet)
	engine := rule_engine.NewEngine(ruleSet)

	age := 17
	name := "test_name"
	city := "IXR"
	postalCode := "560102"

	type Details struct {
		Name *string `json:"name"`
	}

	type Address struct {
		City       *string `json:"city"`
		PostalCode *string `json:"postalcode"`
	}

	details := &Details{
		Name: &name,
	}

	type TestStruct struct {
		Age     *int     `json:"age"`
		Details *Details `json:"details"`
		Address *Address `json:"address"`
	}

	type paramStr struct {
		Test *TestStruct `json:"test"`
	}

	newParam := &paramStr{
		Test: &TestStruct{
			Age:     &age,
			Details: details,
			Address: &Address{
				City:       &city,
				PostalCode: &postalCode,
			},
		},
	}

	x, _ := json.Marshal(newParam)
	var o map[string]interface{}
	json.Unmarshal(x, &o)

	fmt.Println("o: ", o)

	result, extras, err := engine.Run(o)
	fmt.Println("Result: ", result)
	fmt.Println("Extras: ", extras)
	fmt.Println("Errors: ", err)
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
