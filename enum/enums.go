package enum

type Matchers string
type Iterators string

const (
	Any Matchers = "any"
	All Matchers = "all"

	RuleSetIterator    Iterators = "$$$$"
	ExpressionIterator Iterators = "$$$"
	RulerIterator      Iterators = "$$"
	RuleIterator       Iterators = "$"
)
