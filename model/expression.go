package model

import (
	"errors"

	"fmt"

	"reflect"
	"strings"

	"github.com/utk1369/infer-go/enum"
	"github.com/utk1369/infer-go/model/ruler"
	"github.com/utk1369/infer-go/util"
)

type Expression struct {
	Require     *enum.Matchers `json:"require"`
	Name        *string        `json:"name"`
	Description *string        `json:"desc"`
	IterateOn   *string        `json:"iterateOn"`
	Rulers      []ruler.Ruler  `json:"rulers"`
}

func (expression *Expression) Test(param map[string]interface{}) (bool, map[string]interface{}, error) {
	expression.ReplaceIteratorTokens(param)
	ruleMatcher := *expression.Require
	for _, ruler := range expression.Rulers {
		eval, extras, err := ruler.Test(param)

		if err != nil {
			return false, nil, err
		} else if ruleMatcher == enum.All && !eval {
			extras["expression"] = *expression.Name
			return false, extras, nil
		} else if ruleMatcher == enum.Any && eval {
			return true, nil, nil
		}
	}

	if ruleMatcher == enum.Any {
		return false, map[string]interface{}{"expression": *expression.Name}, nil
	} else if ruleMatcher == enum.All {
		return true, nil, nil
	}
	return false, nil, errors.New("invalid require condition: " + fmt.Sprint(ruleMatcher) + " for expression: " + *expression.Name)
}

func (expression *Expression) ReplaceIteratorTokens(param map[string]interface{}) error {
	if expression.IterateOn == nil {
		return nil
	}
	obj := util.Pluck(param, *expression.IterateOn)
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Slice {
		return errors.New("Invalid Path to iterate on: " + *expression.IterateOn)
	}

	t := []ruler.Ruler{}
	for _, r := range expression.Rulers {
		s := reflect.ValueOf(obj)
		for i := 0; i < s.Len(); i++ {
			replacedRuler := new(ruler.Ruler)
			util.DeepClone(r, replacedRuler)
			iteratorTokenFound := false
			for _, rule := range replacedRuler.Rules {
				if strings.Contains(rule.Path, fmt.Sprint(enum.RulerIterator)) {
					iteratorTokenFound = true
				}
				rule.Path = strings.Replace(rule.Path, fmt.Sprint(enum.RulerIterator), fmt.Sprint(i), -1)

			}
			if !iteratorTokenFound {
				t = append(t, r)
				break
			}
			t = append(t, *replacedRuler)
		}
	}
	expression.Rulers = t
	return nil
}

func (expression *Expression) String() string {
	expr := "\nExpression: "
	if expression.Description != nil {
		expr = expr + *expression.Description + ". "
	}

	if expression.IterateOn != nil {
		expr = expr + "\n\tIterate On: " + *expression.IterateOn + ". "
	}

	for _, rulers := range expression.Rulers {
		expr = expr + "\nCondition: " + fmt.Sprint(*expression.Require) + ", \n\tRulers: " + rulers.String()
	}
	return expr
}

func (expression *Expression) GetErrors() []error {
	return []error{}
}
