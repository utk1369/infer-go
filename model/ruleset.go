package model

import (
	"encoding/json"
	"fmt"

	"errors"

	"reflect"
	"strings"

	"github.com/utk1369/infer-go/enum"
	"github.com/utk1369/infer-go/util"
)

type RuleSet struct {
	ID        *string        `json:"id"`
	Name      *string        `json:"name"`
	IterateOn *string        `json:"iterateOn"`
	Require   *enum.Matchers `json:"require"`
	When      []When         `json:"when"`
}

func LoadRuleSetFromJson(jsonRuleSet []byte) (*RuleSet, error) {
	var ruleSet *RuleSet
	err := json.Unmarshal(jsonRuleSet, &ruleSet)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return ruleSet, nil
}

func (rs *RuleSet) Test(param map[string]interface{}) (result bool, extras map[string]interface{}, err error) {
	whenMatcher := *rs.Require
	rs.ReplaceIteratorTokens(param)
	for _, when := range rs.When {
		truth, extras, err := when.Test(param)
		if err != nil {
			return false, nil, err
		} else if whenMatcher == enum.All && truth == false {
			extras["ruleset"] = rs.Name
			return false, extras, err
		} else if whenMatcher == enum.Any && truth == true {
			return true, nil, nil
		}
	}

	if whenMatcher == enum.All {
		return true, nil, nil
	} else if whenMatcher == enum.Any {
		return false, map[string]interface{}{"ruleSet": *rs.Name}, nil
	}
	return false, nil, errors.New("invalid require condition: " + fmt.Sprint(whenMatcher) + " for ruleSet: " + *rs.Name)
}

func (rs *RuleSet) ReplaceIteratorTokens(param map[string]interface{}) error {
	if rs.IterateOn == nil {
		return nil
	}
	obj := util.Pluck(param, *rs.IterateOn)
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Slice {
		return errors.New("Invalid Path to iterate on: " + *rs.IterateOn)
	}

	t := []When{}
	for _, w := range rs.When {
		s := reflect.ValueOf(obj)
		for i := 0; i < s.Len(); i++ {
			replacedWhen := new(When)
			util.DeepClone(w, replacedWhen)
			iteratorTokenFound := false
			for _, expression := range replacedWhen.Expressions {
				if expression.IterateOn != nil {
					*expression.IterateOn = strings.Replace(*expression.IterateOn, fmt.Sprint(enum.RuleSetIterator), fmt.Sprint(i), -1)
				}
				for _, r := range expression.Rulers {
					for _, rule := range r.Rules {
						if strings.Contains(rule.Path, fmt.Sprint(enum.RuleSetIterator)) {
							iteratorTokenFound = true
						}
						rule.Path = strings.Replace(rule.Path, fmt.Sprint(enum.RuleSetIterator), fmt.Sprint(i), -1)
					}
				}
			}
			if !iteratorTokenFound {
				t = append(t, w)
				break
			}
			t = append(t, *replacedWhen)
		}
	}
	rs.When = t
	return nil
}

func (rs *RuleSet) String() string {
	ruleSet := ""
	if rs.ID != nil {
		ruleSet = ruleSet + "ID: " + *rs.ID + ". "
	}

	for _, when := range rs.When {
		ruleSet = ruleSet + "\nRequire: " + fmt.Sprint(*rs.Require) + ", \n\tWhen: " + when.String()
	}

	return ruleSet
}

func (rs *RuleSet) GetErrors() []error {
	return []error{}
}
