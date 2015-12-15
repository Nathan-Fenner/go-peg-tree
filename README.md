# go-peg-tree
A type-safe PEG parser generator for Go

`go get github.com/nathan-fenner/go-peg-tree`

What's this?
============
This project is a parser-generator for the [parsing expression grammar (PEG)](https://en.wikipedia.org/wiki/Parsing_expression_grammar).

It takes a PEG as input, with some Go code mixed in, like this:

```
import "strconv"

alias number <- regex{ [0-9]+ } go float64 { strconv.ParseFloat(arg, 64) };

atom <- "(" Expression ")" go float64 { arg.V1 } / number;

product <-
  (atom "*" product go float64 { arg.V0 * arg.V2 })
  /
  atom;

sum <-
  (product "+" sum go float64 { arg.V0 + arg.V2 })
  /
  product;

Expression <- sum;

```

and it generates a file that contains

```
func NewParser(input string) Parser {
  ...
}
func (parser Parser) Expression() (float64, error) {
  ...
}
```

(it is not yet possible to actually provide PEG input; you must currently build
the grammar out of a collection of PEG types)

(regexes are not yet available)

(imports are also not yet available)

Efficiency
==========
It's basically a recursive-descent parser with memoization. In particular, it
it may require a great deal of space on the stack in order to run, since it
performs the parsing through a recursive matching procedure.

Do you support left-recursion?
==============================
No. Left-recursion makes the "memoized" aspect of the parser hard to write,
which forces you to give up the linear-time property of parsing.

At the moment, however, the parser generator does not even warn you about left-
recursive grammars. This will cause the generated parser to go into an infinite
loop (and quickly overflow the stack).

Error messages?
===============
It produces simple error messages; it will tell you a list of all the allowable
tokens that could follow wherever you are. In addition, you can "alias" certain
nodes so that they'll be identified by a name, rather than listing the tokens
that they expect to come next. This is useful, for example, for abstracting the
contents of an "identifier" definition away from its actual presence.

Type safe?
==========
Every PEG expression results in a value of a specific type. These can be
combined by Go expressions into new values. For example, we can define some of
our types like this (the rest are left to the reader's imagination):

```
type Argument {
  Name string
  Type Type
}

type Func struct {
  Name      string      // the name of the function
  Arguments []Argument  // the arguments to the function
  Return    *Type       // the (optional) return type
  Body      []Statement // the body of the function
}

```

Now our parser for `func` will look something like:

```
argument Argument <- identifier type go { Argument{arg.V0, arg.V1} };

arguments []Argument <-
  ( argument "," arguments go { append(argument, arguments...) } )
  /
  go { []Argument{} };


func Func <-
  "func" identifier "(" arguments ")" type? "{"
    statement*
  "}"
  go { buildFunc(arg) };
```

Then our `buildFunc` function has the following signature:

```
func buildFunc(arg struct{
  V0 string      // "func"
  V1 string      // identifier
  V2 string      // "("
  V3 []Argument  // arguments
  V4 string      // ")"
  V5 *Type       // optional return type
  V6 string      // "{"
  V7 []Statement // the body of the function
  V8 string      // "}"
}) Func {
  return Func {
    Name:      arg.V1,
    Arguments: arg.V3,
    Return:    arg.V5,
    Body:      arg.V7,
  }
}
```

We can make this a lot more beautiful with a few alterations to the syntax in
the grammar definition:

```
argument Argument <- identifier type go { Argument{arg.V0, arg.V1} };

arguments []Argument <-
  ( argument "," arguments go { append(argument, arguments...) } )
  /
  go { []Argument{} };


func Func <-
  "func" name:identifier "(" arguments:arguments ")" returns:(type?) "{"
    body:statement*
  "}"
  go { buildFunc(arg) };

// so now we have

func buildFunc(arg struct{
  name      string      // identifier
  arguments []Argument  // arguments
  returns   *Type       // optional return type
  body      []Statement // the body of the function
}) Func {
  return Func {
    Name:      arg.name,
    Arguments: arg.arguments,
    Return:    arg.returns,
    Body:      arg.body,
  }
}
```
