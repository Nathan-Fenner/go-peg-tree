package core

import (
	"fmt"
	"strings"
)

type Resource struct {
	Name       string
	Type       string
	Expression string
}

type Context struct {
	Resources []Resource
	Imports   []string
}

type Peg interface {
	Template(*State) string
	String() string
	TypeName() string
	Context() Context
}

type Literal string

func (l Literal) Template(state *State) string {
	return fmt.Sprintf(`
if here+%d > len(input) || string(input[here:here+%d]) != %q {
	return Failure(Expected{Token: %q}), ""
}
return Success(here + %d), %q`, len(string(l)), len(string(l)), string(l), string(l), len(string(l)), string(l))
}
func (l Literal) String() string {
	return fmt.Sprintf("%q", string(l))
}
func (l Literal) TypeName() string {
	return "string"
}
func (l Literal) Context() Context {
	return Context{}
}

type Sequence []Peg

func (s Sequence) Template(state *State) string {
	template := "\nresult := " + s.TypeName() + "{}"
	for i := range s {
		template += state.DefineIn(s[i], `
if next, value := %s(input, here); next.Ok {
	here = next.At
	result.V`+fmt.Sprintf("%d", i)+` = value
} else {
	return next, `+s.TypeName()+`{}
}`)
	}
	return template + `
	return Success(here), result`
}
func (s Sequence) String() string {
	pieces := make([]string, len(s))
	for i := range s {
		pieces[i] = s[i].String()
	}
	return strings.Join(pieces, " ")
}
func (s Sequence) TypeName() string {
	name := "struct{"
	for i, p := range s {
		name += fmt.Sprintf("V%d %s;", i, p.TypeName())
	}
	return name + "}"
}
func (s Sequence) Context() Context {
	return Context{}
}

type Alternate []Peg

func (a Alternate) Template(state *State) string {
	template := "\nnotes := []Reject{}"
	for i := range a {
		template += state.DefineIn(a[i], `
if next, value := %s(input, here); next.Ok {
	return next, value
} else {
	notes = append(notes, next.Expected...)
}`)
	}
	template += "\nvar zero " + a.TypeName() + "\nreturn Failure(notes...), zero"
	return template
}
func (a Alternate) String() string {
	pieces := make([]string, len(a))
	for i := range a {
		pieces[i] = a[i].String()
	}
	return "(" + strings.Join(pieces, " / ") + ")"
}
func (a Alternate) TypeName() string {
	return a[0].TypeName()
}
func (a Alternate) Context() Context {
	return Context{}
}

type Star struct {
	Argument Peg
}

func (s Star) Template(state *State) string {
	return state.DefineIn(s.Argument, `
result := []`+s.Argument.TypeName()+`{}
for {
	next, value := %s(input, here)
	if !next.Ok {
		return Success(here), result
	}
	here = next.At
	result = append(result, value)
}`)
}
func (s Star) String() string {
	return "(" + s.String() + ")*"
}
func (s Star) TypeName() string {
	return "[]" + s.Argument.TypeName()
}
func (s Star) Context() Context {
	return Context{}
}

type Not struct {
	Argument Peg
}

func (n Not) Template(state *State) string {
	return state.DefineIn(n.Argument, `
check, _ := %s(input, here)
if !check.Ok {
  return Success(here), struct{}{}
}
return Failure(Exclude{`+n.Argument.String()+`}), struct{}{}`)
}
func (n Not) String() string {
	return "not (" + n.Argument.String() + ")"
}
func (n Not) TypeName() string {
	return "struct{}"
}
func (n Not) Context() Context {
	return Context{}
}

type And struct {
	Argument Peg
}

func (and And) Template(state *State) string {
	return state.DefineIn(and.Argument, `
check, value := %s(input, here)
if !check.Ok {
	var zero `+and.Argument.TypeName()+`
	return check, zero
}
return Success(here), value`)
}
func (and And) String() string {
	return "&(" + and.Argument.String() + ")"
}
func (and And) TypeName() string {
	return and.Argument.TypeName()
}
func (and And) Context() Context {
	return Context{}
}

type Root struct {
	Name string
	Type string
}

func (root Root) Template(state *State) string {
	return `/* illegal - roots should not be generated */`
}
func (root Root) String() string {
	return "root " + root.Name
}
func (root Root) TypeName() string {
	return root.Type
}
func (root Root) Context() Context {
	return Context{}
}

type Go struct {
	Argument   Peg
	Returns    string
	Expression string
}

func (g Go) Template(state *State) string {
	return state.DefineIn(g.Argument, `
check, value := %s(input, here)
if !check.Ok {
	var zero `+g.Returns+`
	return check, zero
}
answer := func(arg `+g.Argument.TypeName()+`) `+g.Returns+` {
	return `+g.Expression+`
}(value)
return check, answer`)
}
func (g Go) String() string {
	return fmt.Sprintf("%s go { %s }", g.Argument.TypeName(), g.Returns)
}
func (g Go) TypeName() string {
	return g.Returns
}
func (g Go) Context() Context {
	return Context{}
}

type Regex struct {
	Regex string
}

func (r Regex) Template(state *State) string {
	return fmt.Sprintf(`
match := parser.resource%sRegex.Index(input[here:])
if match == nil || match[0] != 0 {
	return Failure(Expected{Token: "regex " + %q})
}
end := match[1]
return Success(here + end), string(input[here : here+end])
`, state.Current, r.Regex)
}
func (r Regex) String() string {
	return fmt.Sprintf("regex %q", r.Regex)
}
func (r Regex) TypeName() string {
	return "string"
}
func (r Regex) Context() Context {
	return Context{
		Imports:   []string{"regexp"},
		Resources: []Resource{{Name: "Regex", Type: "*regexp.Regexp", Expression: fmt.Sprintf(`regexp.MustCompile(%q)`, r.Regex)}},
	}
}
