package model

import (
	"errors"

	"fmt"

	"reflect"

	"strings"

	"bitbucket.org/grabpay/infer-go/enum"
	"bitbucket.org/grabpay/infer-go/util"
)

type When struct {
	Require     *enum.Matchers `json:"require"`
	Name        *string        `json:"name"`
	Description *string        `json:"desc"`
	IterateOn   *string        `json:"iterateOn"`
	Expressions []Expression   `json:"expressions"`
}

func (when *When) Test(param map[string]interface{}) (isSatisifed bool, extras map[string]interface{}, err error) {

	expressionMatcher := *when.Require
	when.ReplaceIteratorTokens(param)
	for _, expression := range when.Expressions {
		isSatisfied, extras, err := expression.Test(param)

		if err != nil {
			return false, nil, err
		} else if expressionMatcher == enum.All && isSatisfied == false {
			extras[reflect.TypeOf(when).String()] = *when.Name
			return false, extras, err
		} else if expressionMatcher == enum.Any && isSatisfied == true {
			return true, nil, nil
		}
	}
	if expressionMatcher == enum.Any {
		return false, map[string]interface{}{"condition": *when.Name}, nil
	} else if expressionMatcher == enum.All {
		return true, nil, nil
	}
	return false, nil, errors.New("invalid require condition: " + fmt.Sprint(expressionMatcher) + " for when: " + *when.Name)
}

func (w *When) ReplaceIteratorTokens(param map[string]interface{}) error {
	if w.IterateOn == nil {
		return nil
	}
	obj := util.Pluck(param, *w.IterateOn)
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Slice {
		return errors.New("Invalid Path to iterate on: " + *w.IterateOn)
	}

	t := []Expression{}
	for _, e := range w.Expressions {
		s := reflect.ValueOf(obj)
		for i := 0; i < s.Len(); i++ {
			replacedExpression := new(Expression)
			util.DeepClone(e, replacedExpression)
			iteratorTokenFound := false
			for _, ruler := range replacedExpression.Rulers {
				if ruler.IterateOn != nil {
					*ruler.IterateOn = strings.Replace(*ruler.IterateOn, fmt.Sprint(enum.ExpressionIterator), fmt.Sprint(i), -1)
				}
				for _, rule := range ruler.Rules {
					if strings.Contains(rule.Path, fmt.Sprint(enum.ExpressionIterator)) {
						iteratorTokenFound = true
					}
					rule.Path = strings.Replace(rule.Path, fmt.Sprint(enum.ExpressionIterator), fmt.Sprint(i), -1)
				}
			}
			if !iteratorTokenFound {
				t = append(t, e)
				break
			}
			t = append(t, *replacedExpression)
		}
	}
	w.Expressions = t
	return nil
}

func (when *When) String() string {
	w := "When: "
	if when.Description != nil {
		w = w + *when.Description + ". "
	}

	for _, expression := range when.Expressions {
		w = w + "\nRequire: " + fmt.Sprint(*when.Require) + ", \n\tExpressions: " + expression.String()
	}
	return w
}

func (w *When) GetErrors() []error {
	return []error{}
}

func (w *When) ElementToIterateOn(iterateOn *string) *When {
	w.IterateOn = iterateOn
	return w
}

func (w *When) Matcher(require enum.Matchers) *When {
	w.Require = &require
	return w
}

func (w *When) Expression(expression Expression) *When {
	w.Expressions = append(w.Expressions, expression)
	return w
}

func NewWhenCondition(name string, desc string) *When {
	return &When{
		Name:        &name,
		Description: &desc,
	}
}
