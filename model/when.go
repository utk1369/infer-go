package model

import "errors"

type When struct {
	Require     *string      `json:"require"`
	Description *string      `json:"desc"`
	Expressions []Expression `json:"expressions"`
}

func (when *When) IsSatisfied(params ...map[string]interface{}) (isSatisifed bool, extras map[string]interface{}, err error) {

	expressionMatcher := *when.Require
	for _, expression := range when.Expressions {
		isSatisfied, extras, err := expression.Evaluate(params[0])

		if err != nil {
			return false, nil, err
		} else if expressionMatcher == "all" && isSatisfied == false {
			extras["when"] = when
			return false, extras, err
		} else if expressionMatcher == "any" && isSatisfied == true {
			return true, nil, nil
		}
	}
	if expressionMatcher == "any" {
		return false, map[string]interface{}{"when": *when}, nil
	} else if expressionMatcher == "all" {
		return true, nil, nil
	}
	return false, nil, errors.New("invalid require condition: " + expressionMatcher + " for when: " + *when.Description)
}

func (when *When) String() string {
	w := "When: "
	if when.Description != nil {
		w = w + *when.Description + ". "
	}

	for _, expression := range when.Expressions {
		w = w + "\nCondition: " + *when.Require + ", \n\tExpressions: " + expression.String()
	}
	return w
}
