package model

type Condition interface {
	// takes in the parameter as a map and returns the result along with any extra info, errors if any
	Test(param map[string]interface{}) (result bool, extras map[string]interface{}, err error)
	// takes in the parameter as a map and replaces all the iterator tokens of the containing conditions with the indices
	ReplaceIteratorTokens(param map[string]interface{}) error
	// toString() method for the condition
	String() string
	// validates the condition
	GetErrors() []error
}
