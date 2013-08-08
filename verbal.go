package verbal

import (
	"regexp"
	"strings"
)

/*
TODO: Work with strings in a more efficient way.
TODO: Add an option to have a lazy regex that is only compiled when computing.
TODO: Documentation
*/

type Flag uint

const (
	MULTILINE   Flag = 1 << iota
	IGNORE_CASE Flag = 1 << iota
	DOTALL      Flag = 1 << iota
	UNGREEDY    Flag = 1 << iota
)

type Expression struct {
	prefixes string
	source   string
	suffixes string
	pattern  string
	flags    Flag
	compiled *regexp.Regexp
}

func NewExpression() *Expression {
	return &Expression{flags: MULTILINE}
}

func sanitize(value string) string {
	return regexp.QuoteMeta(value)
}

func (e *Expression) getFlags() string {
	flags := "misU"
	result := []rune{}

	for i, flag := range flags {
		if e.flags&1<<uint(i) != 0 {
			result = append(result, flag)
		}
	}

	return string(result)
}

func (e *Expression) update() {
	compile := strings.Join([]string{"(?", e.getFlags(), ")", e.prefixes, e.source, e.suffixes}, "")
	e.compiled = regexp.MustCompile(compile)
}

func (e *Expression) Add(values ...string) *Expression {
	value := strings.Join(values, "")
	e.source = strings.Join([]string{e.source, value}, "")
	e.update()

	return e
}

func (e *Expression) HasStartOfLine(enable bool) *Expression {
	if enable {
		e.prefixes = "^"
	} else {
		e.prefixes = ""
	}

	e.update()
	return e
}

func (e *Expression) StartOfLine() *Expression {
	return e.HasStartOfLine(true)
}

func (e *Expression) HasEndOfLine(enable bool) *Expression {
	if enable {
		e.suffixes = "$"
	} else {
		e.suffixes = ""
	}

	e.update()
	return e
}

func (e *Expression) EndOfLine() *Expression {
	return e.HasEndOfLine(true)
}

func (e *Expression) AnyCase() *Expression {
	e.AddModifier(IGNORE_CASE)
	return e
}

func (e *Expression) OneLine() *Expression {
	e.RemoveModifier(MULTILINE)
	return e
}

func (e *Expression) MatchAllWithDot() *Expression {
	e.AddModifier(DOTALL)
	return e
}

func (e *Expression) AddModifier(f Flag) *Expression {
	e.flags |= f
	return e
}

func (e *Expression) RemoveModifier(f Flag) *Expression {
	e.flags &= ^f
	return e
}

func (e *Expression) Then(value string) *Expression {
	v := sanitize(value)
	return e.Add("(", v, ")")
}

func (e *Expression) Find(value string) *Expression {
	return e.Then(value)
}

func (e *Expression) Maybe(value string) *Expression {
	v := sanitize(value)
	return e.Add("(", v, "?)")
}

func (e *Expression) AtLeastOne(value string) *Expression {
	v := sanitize(value)
	return e.Add("(", v, "+)")
}

func (e *Expression) AnyNumberOf(value string) *Expression {
	v := sanitize(value)
	return e.Add("(", v, "*)")
}

func (e *Expression) Multiple(value string) *Expression {
	v := sanitize(value)
	return e.Add("(", v, "{2,})")
}

func (e *Expression) Anything() *Expression {
	return e.Add("(.*")
}

func (e *Expression) AnythingBut(value string) *Expression {
	v := sanitize(value)
	return e.Add("([^", v, "]?)")
}

func (e *Expression) Something() *Expression {
	return e.Add("(.+")
}

func (e *Expression) SomethingBut(value string) *Expression {
	v := sanitize(value)
	return e.Add("([^", v, "]+)")
}

func (e *Expression) LineBreak() *Expression {
	return e.Add("(\\n|(\\r\\n))")
}

func (e *Expression) Br() *Expression {
	return e.Add("(\\n|(\\r\\n))")
}

func (e *Expression) Tab() *Expression {
	return e.Add("\\t")
}

func (e *Expression) Word() *Expression {
	return e.Add("\\w+")
}

func (e *Expression) AnyOf(value string) *Expression {
	v := sanitize(value)
	return e.Add("[", v, "]")
}

func (e *Expression) Any(value string) *Expression {
	return e.AnyOf(value)
}

func (e *Expression) Or(value string) *Expression {
	v := sanitize(value)
	e.source = strings.Join([]string{"(?:(", source, ")|(?:", v, "))"}, "")
	return e
}

func (e *Expression) Test(test string) bool {
	e.update()
	return e.compiled.MatchString(test)
}

func (e *Expression) Replace(source string, repl string) string {
	e.update()
	return e.compiled.ReplaceAllString(source, repl)
}

func (e *Expression) String() string {
	e.update()
	return e.compiled.String()
}

func (e *Expression) Range(objs []interface{}) *Expression {
	return nil
}
