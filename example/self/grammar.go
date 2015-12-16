package self

import "fmt"
import "regexp"

func NewParser(input string) Parser {
	return Parser{
		input:           []byte(input),
		wherem1:         map[int]Result{},
		whatm1:          map[int]string{},
		resourcem1Regex: regexp.MustCompile("[a-zA-Z_][a-zA-Z_0-9]*"),
		wherem0:         map[int]Result{},
		whatm0:          map[int]string{},
	}
}

type Parser struct {
	input []byte
	// Internal memoization tables
	wherem1 map[int]Result
	whatm1  map[int]string
	m1Regex *regexp.Regexp
	wherem0 map[int]Result
	whatm0  map[int]string
}

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

func (parser Parser) m0(input []byte, here int) (Result, string) {
	return parser.m1(input, here)
}

var wherem1 = map[int]Result{}
var whatm1 = map[int]string{}

func (parser Parser) m1(input []byte, here int) (Result, string) {
	if result, ok := parser.wherem1[here]; ok {
		return result, parser.whatm1[here]
	}
	result, value := parser.dm1(input, here)
	parser.wherem1[here] = result
	parser.whatm1[here] = value
	return result, value
}

// regex "[a-zA-Z_][a-zA-Z_0-9]*"
func (parser Parser) dm1(input []byte, here int) (Result, string) {
	match := parser.resourcem1Regex.Index(input[here:])
	if match == nil || match[0] != 0 {
		return Failure(Expected{Token: "regex " + "[a-zA-Z_][a-zA-Z_0-9]*"})
	}
	end := match[1]
	return Success(here + end), string(input[here : here+end])

}
