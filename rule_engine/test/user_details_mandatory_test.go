package test

import (
	"testing"

	"github.com/utk1369/infer-go/model/ruler"
	"github.com/utk1369/infer-go/rule_engine"
)

func TestEngineForRulerAllSuccess(t *testing.T) {
	rulerJson := getJsonByteArray("test_mandatory_details.json")
	data := &User{
		Details: &Details{
			Name: "Bob",
			Age:  23,
		},
		Verified: "Y",
	}

	r, err1 := ruler.NewRulerWithJson(rulerJson)

	if err1 != nil {
		panic("Could not create ruler from json: " + err1.Error())
	}

	param := getMap(data)
	engine := rule_engine.NewEngine(r, param)
	res, extras, err2 := engine.Run()

	if err2 != nil {
		panic("Could not execute condition: " + err2.Error())
	}

	if res == false {
		t.Error("Failed [Expected -> true, Returned -> false], Extra Info : ", extras)
		t.Failed()
	}
}

func TestEngineForRulerAllFailure(t *testing.T) {
	rulerJson := getJsonByteArray("test_mandatory_details.json")
	data := &User{
		Details: &Details{
			Name: "Bob",
			Age:  5,
		},
		Verified: "Y",
	}

	r, err1 := ruler.NewRulerWithJson(rulerJson)

	if err1 != nil {
		panic("Could not create ruler from json: " + err1.Error())
	}

	param := getMap(data)
	engine := rule_engine.NewEngine(r, param)
	res, extras, err2 := engine.Run()

	if err2 != nil {
		panic("Could not execute condition: " + err2.Error())
	}

	if res == true {
		t.Error("Failed: Expected -> false, Returned -> true : ", extras)
		t.Failed()
	}
}
