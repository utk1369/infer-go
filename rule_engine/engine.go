package rule_engine

import (
	"errors"

	"bitbucket.org/grabpay/go-ruler/model"
)

type Engine struct {
	ruleSet *model.RuleSet
}

func NewEngine(r *model.RuleSet) *Engine {
	return &Engine{
		ruleSet: r,
	}
}

func (engine *Engine) Run(params ...map[string]interface{}) (bool, map[string]interface{}, error) {
	ruleSet := engine.ruleSet
	whenMatcher := *ruleSet.Require
	for _, when := range ruleSet.When {
		truth, extras, err := when.IsSatisfied(params[0])
		if err != nil {
			return false, nil, err
		} else if whenMatcher == "all" && truth == false {
			extras["ruleset"] = ruleSet
			return false, extras, err
		} else if whenMatcher == "any" && truth == true {
			return true, nil, nil
		}
	}

	if whenMatcher == "all" {
		return true, nil, nil
	} else if whenMatcher == "any" {
		return false, map[string]interface{}{"ruleSet": ruleSet}, nil
	}
	return false, nil, errors.New("invalid require condition: " + whenMatcher + " for ruleSet: " + *ruleSet.Name)
}
