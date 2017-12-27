package model

import (
	"encoding/json"
	"fmt"
)

type RuleSet struct {
	ID      *string `json:"id"`
	Name    *string `json:"name"`
	Policy  *string `json:"policy"`
	Country *string `json:"country"`
	Require *string `json:"require"`
	When    []When  `json:"when"`
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

func (rs *RuleSet) String() string {
	ruleSet := ""
	if rs.ID != nil {
		ruleSet = ruleSet + "ID: " + *rs.ID + ". "
	}

	for _, when := range rs.When {
		ruleSet = ruleSet + "\nCondition: " + *rs.Require + ", \n\tWhen: " + when.String()
	}

	return ruleSet
}
