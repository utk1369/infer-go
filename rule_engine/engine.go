package rule_engine

import (
	"bitbucket.org/grabpay/infer-go/model"
)

type Engine struct {
	condition model.Condition
	param     map[string]interface{}
}

func NewEngine(c model.Condition, param map[string]interface{}) *Engine {
	return &Engine{
		condition: c,
		param:     param,
	}
}

func (engine *Engine) Run() (bool, map[string]interface{}, error) {
	if engine.condition.GetErrors() != nil && len(engine.condition.GetErrors()) > 0 {
		return false, nil, engine.condition.GetErrors()[0]
	}
	return engine.condition.Test(engine.param)
}
