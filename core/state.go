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
		Roots:       map[string]string{},
		Definitions: map[string]Definition{},
	}
}

type Definition struct {
	Result string
	Body   string
}

type State struct {
	UID         int
	Roots       map[string]string
	Definitions map[string]Definition
}

func (state *State) UniqueID() string {
	s := fmt.Sprintf("m%d", state.UID)
	state.UID++
	return s
}

func (state *State) GetRootID(root string) string {
	name, ok := state.Roots[root]
	if ok {
		return name
	}
	state.Roots[root] = state.UniqueID()
	return state.Roots[root]
}

func (state *State) DefineWithName(template string, name string, returns string, detail string) {
	state.Definitions[name] = Definition{
		Result: returns,
		Body: `
var where` + name + ` = map[int]Result{}
var what` + name + ` = map[int]` + returns + `{}

func ` + name + `(input []rune, here int) (Result, ` + returns + `) {
	if result, ok := where` + name + `[here]; ok {
		return result, what` + name + `[here]
	}
	result, value := d` + name + `(input, here)
	where` + name + `[here] = result
	what` + name + `[here] = value
	return result, value
}

// ` + detail + `
func d` + name + `(input []rune, here int) (Result, ` + returns + `) {` +
			strings.Replace(template, "\n", "\n\t", -1) + `
}`}
}

func (state *State) DefineRoot(root string, peg Peg) {
	name := state.GetRootID(root)
	state.Definitions[name] = Definition{
		Result: peg.TypeName(),
		Body: `
func ` + name + ` (input []rune, here int) (Result, ` + peg.TypeName() + `) {
  return ` + state.Define(peg) + `(input, here)
}
`,
	}
}

func (state *State) Define(peg Peg) string {
	if root, ok := peg.(Root); ok {
		return state.GetRootID(root.Name)
	}
	template := peg.Template(state)
	id := state.UniqueID()
	state.DefineWithName(template, id, peg.TypeName(), peg.String())
	return id
}
func (state *State) DefineIn(peg Peg, source string) string {
	id := state.Define(peg)
	return fmt.Sprintf(source, id)
}

func (state *State) Generate(packageName string) string {
	file := `package ` + packageName + `

import "fmt"
`

	exported := []string{}
	for root, _ := range state.Roots {
		if unicode.IsUpper([]rune(root)[0]) {
			exported = append(exported, root)
		}
	}
	for _, root := range exported {
		definition := state.Definitions[state.Roots[root]]
		id := state.GetRootID(root)
		file += `
func Parse` + root + `(input string) (` + definition.Result + `, error) {
	check, value := ` + id + `([]rune(input), 0)
	if check.Ok {
		return value, nil
	}
	var zero ` + definition.Result + `
	return zero, fmt.Errorf("%s", check.Explain())
}`
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
	s := "Failed to parse. Expected at " + fmt.Sprintf("%d", r.At) + " one of:"
	for _, v := range r.Expected {
		s += "\n\t" + v.Reason()
	}
	return s
}

type Expected struct {
	Token string
}

func (e Expected) Reason() string {
	return fmt.Sprintf("%q", e.Token)
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
	Message string
}

func (e Exclude) Reason() string {
	return fmt.Sprintf("but not %s", e.Message)
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
