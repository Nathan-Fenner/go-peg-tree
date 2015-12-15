package core

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

func NewState() State {
	return State{
		UID:         0,
		Exported:    map[string]string{},
		Definitions: map[string]Definition{},
	}
}

type Definition struct {
	Result string
	Body   string
}

type State struct {
	UID         int
	Exported    map[string]string
	Definitions map[string]Definition
}

func (state *State) DefineWithName(template string, name string, returns string, detail string) {
	state.Definitions[name] = Definition{
		Result: returns,
		Body: `
var memoWhere` + name + ` = map[int]Result{}
var memoWhat` + name + ` = map[int]` + returns + `{}

func match` + name + `(input []rune, here int) (Result, ` + returns + `) {
	if result, ok := memoWhere` + name + `[here]; ok {
		return result, memoWhat` + name + `[here]
	}
	result, value := matchDetail` + name + `(input, here)
	memoWhere` + name + `[here] = result
	memoWhat` + name + `[here] = value
	return result, value
}

// ` + detail + `
func matchDetail` + name + `(input []rune, here int) (Result, ` + returns + `) {` +
			strings.Replace(template, "\n", "\n\t", -1) + `
}`}
}

func (state *State) DefineRoot(name string, peg Peg) {
	state.DefineWithName(peg.Template(state), "root"+name, peg.TypeName(), peg.String())
}

func (state *State) Define(peg Peg) string {
	template := peg.Template(state)
	id := state.UID
	state.UID++
	state.DefineWithName(template, fmt.Sprintf("%d", id), peg.TypeName(), peg.String())
	return fmt.Sprintf("%d", id)
}
func (state *State) DefineIn(peg Peg, source string) string {
	id := state.Define(peg)
	return fmt.Sprintf(source, id)
}

func (state *State) Generate(name string) string {
	file := `package ` + name + `

import "fmt"
`

	exported := []string{}
	for name := range state.Definitions {
		if unicode.IsUpper([]rune(name)[0]) {
			exported = append(exported, name)
		}
	}
	for _, name := range exported {
		definition := state.Definitions[name]
		if unicode.IsUpper([]rune(name)[0]) {
			file += `
func Parse` + name + `(string input) (` + definition.Result + `) {
	check, value := matchroot%s([]rune(input), 0)
	if check.Ok {
		return value, nil
	}
	var zero ` + definition.Result + `
	return zero, fmt.Errorf("%s", check.Expected.Explain())
}`
		}
	}

	file += `
// Below is the internal generated parse structure.
// It's not very efficient right now, but is accomplishes parsing in linear time.
// Currently, there's no way to parse multiple inputs, due to the fact that the
// state of the parse is stored in global variables.

type Result struct {
	Ok       bool
	At       int
	Expected []Reject
}

type Reject interface {
	Reason() string
}

func (r Result) Explain() string {
	if r.Ok {
		return fmt.Sprintf("Okay: %d characters parsed", r.At)
	}
	s := "Failed to parse."
	for _, v := range r.Expected {
		s += "\n\t" + v.Reason()
	}
	return s
}

type Expected struct {
	At    int
	Token string
}

func (e Expected) Reason() string {
	return fmt.Sprintf("expected %q at %d", e.Token, e.At)
}

func Failure(tokens ...Reject) Result {
	return Result{
		Ok:       false,
		Expected: tokens,
	}
}
func FailureCombined(first []Reject, second []Reject) Result {
	return Result{
		Ok:       false,
		Expected: append(append([]Reject{}, first...), second...),
	}
}
func Success(at int) Result {
	return Result{
		Ok: true,
		At: at,
	}
}

type Exclude struct {
	At      int
	Message string
}

func (e Exclude) Reason() string {
	return fmt.Sprintf("%s is not allowed (at %d)", e.Message, e.At)
}
`

	names := []string{}
	for key := range state.Definitions {
		names = append(names, key)
	}
	sort.Strings(names)
	for _, key := range names {
		file += state.Definitions[key].Body + "\n"
	}

	return file
}
