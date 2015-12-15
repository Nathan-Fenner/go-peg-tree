package core

import (
	"fmt"
	"strings"
)

type Peg interface {
	Template(*State) string
	String() string
	TypeName() string
}

type Literal string

func (l Literal) Template(state *State) string {
	return fmt.Sprintf(`
if here+%d >= len(input) || string(input[here:here+%d]) != %q {
	return Failure(Expected{Token: %q, At: here}), ""
}
return Success(here + %d), %q`, len([]rune(string(l))), len([]rune(string(l))), string(l), string(l), len([]rune(string(l))), string(l))
}
func (l Literal) String() string {
	return fmt.Sprintf("%q", string(l))
}
func (l Literal) TypeName() string {
	return "string"
}

type Sequence []Peg

func (s Sequence) Template(state *State) string {
	template := "\nresult := " + s.TypeName() + "{}"
	for i := range s {
		template += state.DefineIn(s[i], `
if next, value := match%s(input, here); next.Ok {
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

type Alternate []Peg

func (a Alternate) Template(state *State) string {
	template := "\nnotes := []Reject{}"
	for i := range a {
		template += state.DefineIn(a[i], `
if next, value := match%s(input, here); next.Ok {
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

type Star struct {
	Argument Peg
}

func (s Star) Template(state *State) string {
	return state.DefineIn(s.Argument, `
result := []`+s.Argument.TypeName()+`{}
for {
	next, value := match%s(input, here)
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

type Not struct {
	Argument Peg
}

func (n Not) Template(state *State) string {
	return state.DefineIn(n.Argument, `
check, _ := match%s(input, here)
if !check.Ok {
  return Success(here), struct{}{}
}
return Failure(Exclude{
  At:      here,
  Message: `+n.Argument.String()+`,
}), struct{}{}`)
}
func (n Not) String() string {
	return "not (" + n.Argument.String() + ")"
}
func (n Not) TypeName() string {
	return "struct{}"
}

type And struct {
	Argument Peg
}

func (and And) Template(state *State) string {
	return state.DefineIn(and.Argument, `
check, value := match%s(input, here)
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

type Root struct {
	Name string
	Type string
}

func (root Root) Template(state *State) string {
	return `
return matchroot` + root.Name + "(input, here)"
}
func (root Root) String() string {
	return root.Name
}
func (root Root) TypeName() string {
	return root.Type
}

type Go struct {
	Argument   Peg
	Returns    string
	Expression string
}

func (g Go) Template(state *State) string {
	return state.DefineIn(g.Argument, `
check, value := match%s(input, here)
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
