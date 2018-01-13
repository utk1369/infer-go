package ruler

import (
	"encoding/json"
	"reflect"
	"regexp"

	"fmt"

	"errors"

	"strings"

	"github.com/tj/go-debug"
	"github.com/utk1369/infer-go/enum"
	"github.com/utk1369/infer-go/util"
)

var ruleDebug = debug.Debug("ruler:rule")

// we'll use these values
// to avoid passing strings to our
// special comparison func for these comparators
const (
	eq        = iota
	neq       = iota
	gt        = iota
	gte       = iota
	lt        = iota
	lte       = iota
	exists    = iota
	nexists   = iota
	regex     = iota
	matches   = iota
	contains  = iota
	ncontains = iota
)

type Ruler struct {
	Name      *string        `json:"name"`
	Require   *enum.Matchers `json:"require"`
	IterateOn *string        `json:"iterateOn"`
	Rules     []*Rule        `json:"rules"`
}

func (r *Ruler) Test(o map[string]interface{}) (bool, map[string]interface{}, error) {

	r.ReplaceIteratorTokens(o)
	// Test All By Default for backward compatibility
	if r.Require == nil || *r.Require == enum.All {
		v, extra := testAll(r, o)
		return v, extra, nil
	} else if *r.Require == enum.Any {
		v, extra := testAny(r, o)
		return v, extra, nil
	}
	return false, nil, errors.New("invalid require condition: " + fmt.Sprint(*r.Require) + " for ruler: " + r.String())
}

func (ruler *Ruler) ReplaceIteratorTokens(param map[string]interface{}) error {
	if ruler.IterateOn == nil {
		return nil
	}
	obj := util.Pluck(param, *ruler.IterateOn)
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Slice {
		return errors.New("Invalid Path to iterate on: " + *ruler.IterateOn)
	}

	t := []*Rule{}
	for _, rule := range ruler.Rules {
		s := reflect.ValueOf(obj)
		for i := 0; i < s.Len(); i++ {
			replacedRule := new(Rule)
			util.DeepClone(rule, replacedRule)
			replacedRule.Path = strings.Replace(rule.Path, fmt.Sprint(enum.RuleIterator), fmt.Sprint(i), -1)
			t = append(t, replacedRule)
		}
	}
	ruler.Rules = t
	return nil
}

func (ruler *Ruler) String() string {
	rules := "Rules: "
	for _, rule := range ruler.Rules {
		rules = rules + "\n" + "\t" + rule.String() + ", "
	}
	return rules
}

func (ruler *Ruler) GetErrors() []error {
	return []error{}
}

// creates a new Ruler for you
// optionally accepts a pointer to a slice of filters
// if you have filters that you want to start with
func NewRuler(rules []*Rule) *Ruler {
	if rules != nil {
		return &Ruler{
			Require:   nil,
			IterateOn: nil,
			Rules:     rules,
		}
	}

	return &Ruler{}
}

// returns a new ruler with filters parsed from JSON data
// expects JSON as a slice of bytes and will parse your JSON for you!
func NewRulerWithJson(jsonstr []byte) (*Ruler, error) {
	var ruler *Ruler

	err := json.Unmarshal(jsonstr, &ruler)
	if err != nil {
		return nil, err
	}

	return ruler, nil
}

// adds a new rule for the property at `path`
// returns a RulerFilter that you can use to add conditions
// and more filters
func (r *Ruler) Rule(path string) *RulerRule {
	rule := &Rule{
		"",
		path,
		nil,
	}

	r.Rules = append(r.Rules, rule)

	return &RulerRule{
		r,
		rule,
	}
}

// tests all the Rules (i.e. filters) in your set of Rules,
// given a map that looks like a JSON object
// (map[string]interface{})
func testAll(r *Ruler, o map[string]interface{}) (bool, map[string]interface{}) {
	for _, f := range r.Rules {
		val := util.Pluck(o, f.Path)
		failureExtra := "Rule: [" + fmt.Sprint(f) + "] did not pass."

		if val != nil && f.Value != nil {
			// both the actual and expected value must be comparable
			a := reflect.TypeOf(val)
			e := reflect.TypeOf(f.Value)

			if !a.Comparable() || !e.Comparable() {
				return false, map[string]interface{}{"Ruler": failureExtra}
			}

			if !r.compare(f, val) {
				return false, map[string]interface{}{"Ruler": failureExtra}
			}
		} else if f.Value == nil && (f.Comparator == "exists" || f.Comparator == "nexists") {
			if !r.compare(f, val) {
				return false, map[string]interface{}{"Ruler": failureExtra}
			}
		} else if val == nil && f.Value != nil {
			return false, map[string]interface{}{"Ruler": failureExtra}
		} else {
			ruleDebug("did not find property (%s) on map", f.Path)
			// if we couldn't find the value on the map
			// and the comparator isn't exists/nexists, this fails
			return false, map[string]interface{}{"Ruler": failureExtra}
		}

	}
	return true, nil
}

// tests any of the Rules in your set of Rules,
// given a map that looks like a JSON object
// (map[string]interface{})
func testAny(r *Ruler, o map[string]interface{}) (bool, map[string]interface{}) {
	for _, f := range r.Rules {
		val := util.Pluck(o, f.Path)

		if val != nil && f.Value != nil {
			// both the actual and expected value must be comparable
			a := reflect.TypeOf(val)
			e := reflect.TypeOf(f.Value)

			if !a.Comparable() || !e.Comparable() {
				continue
			}

			if r.compare(f, val) {
				return true, nil
			}
		} else if f.Value == nil && (f.Comparator == "exists" || f.Comparator == "nexists") {
			if r.compare(f, val) {
				return true, nil
			}
		} else if val == nil && f.Value != nil {
			continue
		} else {
			ruleDebug("did not find property (%s) on map", f.Path)
			// if we couldn't find the value on the map
			// and the comparator isn't exists/nexists, this fails
			continue
		}

	}
	return false, map[string]interface{}{"Ruler": *r}
}

// compares real v. actual values
func (r *Ruler) compare(f *Rule, actual interface{}) bool {
	ruleDebug("beginning comparison")
	expected := f.Value
	switch f.Comparator {
	case "eq":
		return actual == expected

	case "neq":
		return actual != expected

	case "gt":
		return r.inequality(gt, actual, expected)

	case "gte":
		return r.inequality(gte, actual, expected)

	case "lt":
		return r.inequality(lt, actual, expected)

	case "lte":
		return r.inequality(lte, actual, expected)

	case "exists":
		// not sure this makes complete sense
		return actual != nil

	case "nexists":
		return actual == nil

	case "regex":
		fallthrough
	case "contains":
		fallthrough
	case "matches":
		return r.regexp(actual, expected)

	case "ncontains":
		return !r.regexp(actual, expected)
	default:
		//should probably return an error or something
		//but this is good for now
		//if comparator is not implemented, return false
		ruleDebug("unknown comparator %s", f.Comparator)
		return false
	}
}

// runs equality comparison
// separated in a different function because
// we need to do another type assertion here
// and some other acrobatics
func (r *Ruler) inequality(op int, actual, expected interface{}) bool {
	// need some variables for these deals
	ruleDebug("entered inequality comparison")
	var cmpStr [2]string
	var cmpUint [2]uint64
	var cmpInt [2]int64
	var cmpFloat [2]float64

	for idx, i := range []interface{}{actual, expected} {
		switch t := i.(type) {
		case uint8:
			cmpUint[idx] = uint64(t)
		case uint16:
			cmpUint[idx] = uint64(t)
		case uint32:
			cmpUint[idx] = uint64(t)
		case uint64:
			cmpUint[idx] = t
		case uint:
			cmpUint[idx] = uint64(t)
		case int8:
			cmpInt[idx] = int64(t)
		case int16:
			cmpInt[idx] = int64(t)
		case int32:
			cmpInt[idx] = int64(t)
		case int64:
			cmpInt[idx] = t
		case int:
			cmpInt[idx] = int64(t)
		case float32:
			cmpFloat[idx] = float64(t)
		case float64:
			cmpFloat[idx] = t
		case string:
			cmpStr[idx] = t
		default:
			ruleDebug("invalid type for inequality comparison")
			return false
		}
	}

	// whichever of these works, we're happy with
	// but if you're trying to compare a string to an int, oh well!
	switch op {
	case gt:
		return cmpStr[0] > cmpStr[1] ||
			cmpUint[0] > cmpUint[1] ||
			cmpInt[0] > cmpInt[1] ||
			cmpFloat[0] > cmpFloat[1]
	case gte:
		return cmpStr[0] >= cmpStr[1] ||
			cmpUint[0] >= cmpUint[1] ||
			cmpInt[0] >= cmpInt[1] ||
			cmpFloat[0] >= cmpFloat[1]
	case lt:
		return cmpStr[0] < cmpStr[1] ||
			cmpUint[0] < cmpUint[1] ||
			cmpInt[0] < cmpInt[1] ||
			cmpFloat[0] < cmpFloat[1]
	case lte:
		return cmpStr[0] <= cmpStr[1] ||
			cmpUint[0] <= cmpUint[1] ||
			cmpInt[0] <= cmpInt[1] ||
			cmpFloat[0] <= cmpFloat[1]
	}

	return false
}

func (r *Ruler) regexp(actual, expected interface{}) bool {
	ruleDebug("beginning regexp")
	// regexps must be strings
	var streg string
	var ok bool
	if streg, ok = expected.(string); !ok {
		ruleDebug("expected value not actually a string, bailing")
		return false
	}

	var astring string
	if astring, ok = actual.(string); !ok {
		ruleDebug("actual value not actually a string, bailing")
		return false
	}

	reg, err := regexp.Compile(streg)
	if err != nil {
		ruleDebug("regexp is bad, bailing")
		return false
	}

	return reg.MatchString(astring)
}
