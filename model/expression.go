package model

import (
	"errors"

	"bitbucket.org/grabpay/go-ruler/model/ruler"
)

type Expression struct {
	Require     *string       `json:"require"`
	Description *string       `json:"desc"`
	Rulers      []ruler.Ruler `json:"rulers"`
}

func (expression *Expression) Evaluate(params ...map[string]interface{}) (bool, map[string]interface{}, error) {
	param := params[0]
	ruleMatcher := *expression.Require
	for _, ruler := range expression.Rulers {
		eval, extras, err := ruler.Test(param)

		if err != nil {
			return false, nil, err
		} else if ruleMatcher == "all" && !eval {
			extras["expression"] = expression
			return false, extras, nil
		} else if ruleMatcher == "any" && eval {
			return true, nil, nil
		}
	}

	if ruleMatcher == "any" {
		return false, map[string]interface{}{"expression": expression}, nil
	} else if ruleMatcher == "all" {
		return true, nil, nil
	}
	return false, nil, errors.New("invalid require condition: " + ruleMatcher + " for expression: " + *expression.Description)
}

func (expression *Expression) String() string {
	expr := "\nExpression: "
	if expression.Description != nil {
		expr = expr + *expression.Description + ". "
	}

	for _, rulers := range expression.Rulers {
		expr = expr + "\nCondition: " + *expression.Require + ", \n\tRulers: " + rulers.String()
	}
	return expr
}
