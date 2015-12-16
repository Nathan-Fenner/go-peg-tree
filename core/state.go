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
		Imports:     []string{"fmt"},
		Definitions: map[string]Definition{},
	}
}

type Definition struct {
	Resources []Resource
	Result    string
	Body      string
}

type State struct {
	UID         int                   // For assigning unique identifiers
	Roots       map[string]string     // The names roots
	Imports     []string              // The imports collectively required
	Definitions map[string]Definition // Definitions (from UID, not name)
	Current     string                // Name of current unit (for resource handling)
}

func (state *State) AddImport(name string) {
	for i := range state.Imports {
		if state.Imports[i] == name {
			return
		}
	}
	state.Imports = append(state.Imports, name)
}

func (state *State) AddImports(names []string) {
	for i := range names {
		state.AddImport(names[i])
	}
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

func (state *State) DefineWithName(template string, name string, returns string, detail string, resources []Resource) {
	state.Definitions[name] = Definition{
		Resources: resources,
		Result:    returns,
		Body: `
var where` + name + ` = map[int]Result{}
var what` + name + ` = map[int]` + returns + `{}

func (parser Parser) ` + name + `(input []byte, here int) (Result, ` + returns + `) {
	if result, ok := parser.where` + name + `[here]; ok {
		return result, parser.what` + name + `[here]
	}
	result, value := parser.d` + name + `(input, here)
	parser.where` + name + `[here] = result
	parser.what` + name + `[here] = value
	return result, value
}

// ` + detail + `
func (parser Parser) d` + name + `(input []byte, here int) (Result, ` + returns + `) {` +
			strings.Replace(template, "\n", "\n\t", -1) + `
}`}
}

func (state *State) DefineRoot(root string, peg Peg) {
	name := state.GetRootID(root)
	state.Definitions[name] = Definition{
		Result: peg.TypeName(),
		Body: `
func (parser Parser)` + name + ` (input []byte, here int) (Result, ` + peg.TypeName() + `) {
  return parser.` + state.Define(peg) + `(input, here)
}
`,
	}
}

func (state *State) Define(peg Peg) string {
	context := peg.Context()
	state.AddImports(context.Imports)
	if root, ok := peg.(Root); ok {
		return state.GetRootID(root.Name)
	}
	id := state.UniqueID()
	state.Current = id
	template := peg.Template(state)
	state.DefineWithName(template, id, peg.TypeName(), peg.String(), context.Resources)
	return id
}
func (state *State) DefineIn(peg Peg, source string) string {
	id := state.Define(peg)
	return fmt.Sprintf(source, "parser."+id)
}

func (state *State) Generate(packageName string) string {
	file := `package ` + packageName + `

`

	sort.Strings(state.Imports)
	for _, name := range state.Imports {
		file += fmt.Sprintf("\nimport %q", name)
	}

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
func (parser Parser) ` + root + `(input string) (` + definition.Result + `, error) {
	check, value := parser.` + id + `([]byte(input), 0)
	if check.Ok {
		return value, nil
	}
	var zero ` + definition.Result + `
	return zero, fmt.Errorf("%s", check.Explain())
}`
	}

	//////////////////////////////

	file += `

func NewParser(input string) Parser {
	return Parser {
		input: []byte(input),`

	for i, definition := range state.Definitions {
		file += `
		where` + i + `: map[int]Result{},
		what` + i + `:  map[int]` + definition.Result + `{},`
		for _, resource := range definition.Resources {
			file += "\n\t\tresource" + i + resource.Name + ": " + resource.Expression + ","
		}
	}
	file += "\n\t}\n}\n\n"

	///////////////////////////
	file += `

type Parser struct {
	input []byte
	// Internal memoization tables`

	for i, definition := range state.Definitions {
		file += `
	where` + i + ` map[int]Result
	what` + i + `  map[int]` + definition.Result
		for _, resource := range definition.Resources {
			file += "\n" + i + resource.Name + " " + resource.Type
		}
	}
	file += "}\n"

	/////////////////////////////

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
