package arithmetic

import "fmt"

func (parser Parser) Expression(input string) (float64, error) {
	check, value := parser.m19([]rune(input), 0)
	if check.Ok {
		return value, nil
	}
	var zero float64
	return zero, fmt.Errorf("%s", check.Explain())
}

func NewParser(input string) Parser {
	return Parser{
		input:    []rune(input),
		wherem13: map[int]Result{},
		whatm13:  map[int]float64{},
		wherem18: map[int]Result{},
		whatm18:  map[int]float64{},
		wherem9:  map[int]Result{},
		whatm9:   map[int]float64{},
		wherem12: map[int]Result{},
		whatm12:  map[int]float64{},
		wherem16: map[int]Result{},
		whatm16: map[int]struct {
			V0 float64
			V1 string
			V2 float64
		}{},
		wherem2:  map[int]Result{},
		whatm2:   map[int]float64{},
		wherem0:  map[int]Result{},
		whatm0:   map[int]float64{},
		wherem3:  map[int]Result{},
		whatm3:   map[int]float64{},
		wherem10: map[int]Result{},
		whatm10:  map[int]string{},
		wherem11: map[int]Result{},
		whatm11:  map[int]float64{},
		wherem19: map[int]Result{},
		whatm19:  map[int]float64{},
		wherem4:  map[int]Result{},
		whatm4:   map[int]string{},
		wherem7:  map[int]Result{},
		whatm7:   map[int]string{},
		wherem8:  map[int]Result{},
		whatm8:   map[int]float64{},
		wherem15: map[int]Result{},
		whatm15:  map[int]string{},
		wherem17: map[int]Result{},
		whatm17:  map[int]float64{},
		wherem14: map[int]Result{},
		whatm14:  map[int]float64{},
		wherem1:  map[int]Result{},
		whatm1:   map[int]string{},
		wherem5:  map[int]Result{},
		whatm5:   map[int]float64{},
		wherem6:  map[int]Result{},
		whatm6:   map[int]float64{},
	}
}

type Parser struct {
	input []rune
	// Internal memoization tables
	wherem13 map[int]Result
	whatm13  map[int]float64
	wherem18 map[int]Result
	whatm18  map[int]float64
	wherem2  map[int]Result
	whatm2   map[int]float64
	wherem0  map[int]Result
	whatm0   map[int]float64
	wherem3  map[int]Result
	whatm3   map[int]float64
	wherem9  map[int]Result
	whatm9   map[int]float64
	wherem12 map[int]Result
	whatm12  map[int]float64
	wherem16 map[int]Result
	whatm16  map[int]struct {
		V0 float64
		V1 string
		V2 float64
	}
	wherem4  map[int]Result
	whatm4   map[int]string
	wherem7  map[int]Result
	whatm7   map[int]string
	wherem8  map[int]Result
	whatm8   map[int]float64
	wherem10 map[int]Result
	whatm10  map[int]string
	wherem11 map[int]Result
	whatm11  map[int]float64
	wherem19 map[int]Result
	whatm19  map[int]float64
	wherem1  map[int]Result
	whatm1   map[int]string
	wherem5  map[int]Result
	whatm5   map[int]float64
	wherem6  map[int]Result
	whatm6   map[int]float64
	wherem15 map[int]Result
	whatm15  map[int]string
	wherem17 map[int]Result
	whatm17  map[int]float64
	wherem14 map[int]Result
	whatm14  map[int]float64
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

func (parser Parser) m0(input []rune, here int) (Result, float64) {
	return parser.m2(input, here)
}

var wherem1 = map[int]Result{}
var whatm1 = map[int]string{}

func (parser Parser) m1(input []rune, here int) (Result, string) {
	if result, ok := parser.wherem1[here]; ok {
		return result, parser.whatm1[here]
	}
	result, value := parser.dm1(input, here)
	parser.wherem1[here] = result
	parser.whatm1[here] = value
	return result, value
}

// "one"
func (parser Parser) dm1(input []rune, here int) (Result, string) {
	if here+3 > len(input) || string(input[here:here+3]) != "one" {
		return Failure(Expected{Token: "one"}), ""
	}
	return Success(here + 3), "one"
}

var wherem10 = map[int]Result{}
var whatm10 = map[int]string{}

func (parser Parser) m10(input []rune, here int) (Result, string) {
	if result, ok := parser.wherem10[here]; ok {
		return result, parser.whatm10[here]
	}
	result, value := parser.dm10(input, here)
	parser.wherem10[here] = result
	parser.whatm10[here] = value
	return result, value
}

// "four"
func (parser Parser) dm10(input []rune, here int) (Result, string) {
	if here+4 > len(input) || string(input[here:here+4]) != "four" {
		return Failure(Expected{Token: "four"}), ""
	}
	return Success(here + 4), "four"
}

var wherem11 = map[int]Result{}
var whatm11 = map[int]float64{}

func (parser Parser) m11(input []rune, here int) (Result, float64) {
	if result, ok := parser.wherem11[here]; ok {
		return result, parser.whatm11[here]
	}
	result, value := parser.dm11(input, here)
	parser.wherem11[here] = result
	parser.whatm11[here] = value
	return result, value
}

// string go { float64 }
func (parser Parser) dm11(input []rune, here int) (Result, float64) {
	check, value := parser.m10(input, here)
	if !check.Ok {
		var zero float64
		return check, zero
	}
	answer := func(arg string) float64 {
		return 4
	}(value)
	return check, answer
}

func (parser Parser) m12(input []rune, here int) (Result, float64) {
	return parser.m13(input, here)
}

var wherem13 = map[int]Result{}
var whatm13 = map[int]float64{}

func (parser Parser) m13(input []rune, here int) (Result, float64) {
	if result, ok := parser.wherem13[here]; ok {
		return result, parser.whatm13[here]
	}
	result, value := parser.dm13(input, here)
	parser.wherem13[here] = result
	parser.whatm13[here] = value
	return result, value
}

// (root one / root two / root three / root four)
func (parser Parser) dm13(input []rune, here int) (Result, float64) {
	notes := []Reject{}
	if next, value := parser.m0(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	if next, value := parser.m3(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	if next, value := parser.m6(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	if next, value := parser.m9(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	var zero float64
	return Failure(notes...), zero
}

func (parser Parser) m14(input []rune, here int) (Result, float64) {
	return parser.m18(input, here)
}

var wherem15 = map[int]Result{}
var whatm15 = map[int]string{}

func (parser Parser) m15(input []rune, here int) (Result, string) {
	if result, ok := parser.wherem15[here]; ok {
		return result, parser.whatm15[here]
	}
	result, value := parser.dm15(input, here)
	parser.wherem15[here] = result
	parser.whatm15[here] = value
	return result, value
}

// "+"
func (parser Parser) dm15(input []rune, here int) (Result, string) {
	if here+1 > len(input) || string(input[here:here+1]) != "+" {
		return Failure(Expected{Token: "+"}), ""
	}
	return Success(here + 1), "+"
}

var wherem16 = map[int]Result{}
var whatm16 = map[int]struct {
	V0 float64
	V1 string
	V2 float64
}{}

func (parser Parser) m16(input []rune, here int) (Result, struct {
	V0 float64
	V1 string
	V2 float64
}) {
	if result, ok := parser.wherem16[here]; ok {
		return result, parser.whatm16[here]
	}
	result, value := parser.dm16(input, here)
	parser.wherem16[here] = result
	parser.whatm16[here] = value
	return result, value
}

// root number "+" root sum
func (parser Parser) dm16(input []rune, here int) (Result, struct {
	V0 float64
	V1 string
	V2 float64
}) {
	result := struct {
		V0 float64
		V1 string
		V2 float64
	}{}
	if next, value := parser.m12(input, here); next.Ok {
		here = next.At
		result.V0 = value
	} else {
		return next, struct {
			V0 float64
			V1 string
			V2 float64
		}{}
	}
	if next, value := parser.m15(input, here); next.Ok {
		here = next.At
		result.V1 = value
	} else {
		return next, struct {
			V0 float64
			V1 string
			V2 float64
		}{}
	}
	if next, value := parser.m14(input, here); next.Ok {
		here = next.At
		result.V2 = value
	} else {
		return next, struct {
			V0 float64
			V1 string
			V2 float64
		}{}
	}
	return Success(here), result
}

var wherem17 = map[int]Result{}
var whatm17 = map[int]float64{}

func (parser Parser) m17(input []rune, here int) (Result, float64) {
	if result, ok := parser.wherem17[here]; ok {
		return result, parser.whatm17[here]
	}
	result, value := parser.dm17(input, here)
	parser.wherem17[here] = result
	parser.whatm17[here] = value
	return result, value
}

// struct{V0 float64;V1 string;V2 float64;} go { float64 }
func (parser Parser) dm17(input []rune, here int) (Result, float64) {
	check, value := parser.m16(input, here)
	if !check.Ok {
		var zero float64
		return check, zero
	}
	answer := func(arg struct {
		V0 float64
		V1 string
		V2 float64
	}) float64 {
		return arg.V0 + arg.V2
	}(value)
	return check, answer
}

var wherem18 = map[int]Result{}
var whatm18 = map[int]float64{}

func (parser Parser) m18(input []rune, here int) (Result, float64) {
	if result, ok := parser.wherem18[here]; ok {
		return result, parser.whatm18[here]
	}
	result, value := parser.dm18(input, here)
	parser.wherem18[here] = result
	parser.whatm18[here] = value
	return result, value
}

// (struct{V0 float64;V1 string;V2 float64;} go { float64 } / root number)
func (parser Parser) dm18(input []rune, here int) (Result, float64) {
	notes := []Reject{}
	if next, value := parser.m17(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	if next, value := parser.m12(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	var zero float64
	return Failure(notes...), zero
}

func (parser Parser) m19(input []rune, here int) (Result, float64) {
	return parser.m14(input, here)
}

var wherem2 = map[int]Result{}
var whatm2 = map[int]float64{}

func (parser Parser) m2(input []rune, here int) (Result, float64) {
	if result, ok := parser.wherem2[here]; ok {
		return result, parser.whatm2[here]
	}
	result, value := parser.dm2(input, here)
	parser.wherem2[here] = result
	parser.whatm2[here] = value
	return result, value
}

// string go { float64 }
func (parser Parser) dm2(input []rune, here int) (Result, float64) {
	check, value := parser.m1(input, here)
	if !check.Ok {
		var zero float64
		return check, zero
	}
	answer := func(arg string) float64 {
		return 1
	}(value)
	return check, answer
}

func (parser Parser) m3(input []rune, here int) (Result, float64) {
	return parser.m5(input, here)
}

var wherem4 = map[int]Result{}
var whatm4 = map[int]string{}

func (parser Parser) m4(input []rune, here int) (Result, string) {
	if result, ok := parser.wherem4[here]; ok {
		return result, parser.whatm4[here]
	}
	result, value := parser.dm4(input, here)
	parser.wherem4[here] = result
	parser.whatm4[here] = value
	return result, value
}

// "two"
func (parser Parser) dm4(input []rune, here int) (Result, string) {
	if here+3 > len(input) || string(input[here:here+3]) != "two" {
		return Failure(Expected{Token: "two"}), ""
	}
	return Success(here + 3), "two"
}

var wherem5 = map[int]Result{}
var whatm5 = map[int]float64{}

func (parser Parser) m5(input []rune, here int) (Result, float64) {
	if result, ok := parser.wherem5[here]; ok {
		return result, parser.whatm5[here]
	}
	result, value := parser.dm5(input, here)
	parser.wherem5[here] = result
	parser.whatm5[here] = value
	return result, value
}

// string go { float64 }
func (parser Parser) dm5(input []rune, here int) (Result, float64) {
	check, value := parser.m4(input, here)
	if !check.Ok {
		var zero float64
		return check, zero
	}
	answer := func(arg string) float64 {
		return 2
	}(value)
	return check, answer
}

func (parser Parser) m6(input []rune, here int) (Result, float64) {
	return parser.m8(input, here)
}

var wherem7 = map[int]Result{}
var whatm7 = map[int]string{}

func (parser Parser) m7(input []rune, here int) (Result, string) {
	if result, ok := parser.wherem7[here]; ok {
		return result, parser.whatm7[here]
	}
	result, value := parser.dm7(input, here)
	parser.wherem7[here] = result
	parser.whatm7[here] = value
	return result, value
}

// "three"
func (parser Parser) dm7(input []rune, here int) (Result, string) {
	if here+5 > len(input) || string(input[here:here+5]) != "three" {
		return Failure(Expected{Token: "three"}), ""
	}
	return Success(here + 5), "three"
}

var wherem8 = map[int]Result{}
var whatm8 = map[int]float64{}

func (parser Parser) m8(input []rune, here int) (Result, float64) {
	if result, ok := parser.wherem8[here]; ok {
		return result, parser.whatm8[here]
	}
	result, value := parser.dm8(input, here)
	parser.wherem8[here] = result
	parser.whatm8[here] = value
	return result, value
}

// string go { float64 }
func (parser Parser) dm8(input []rune, here int) (Result, float64) {
	check, value := parser.m7(input, here)
	if !check.Ok {
		var zero float64
		return check, zero
	}
	answer := func(arg string) float64 {
		return 3
	}(value)
	return check, answer
}

func (parser Parser) m9(input []rune, here int) (Result, float64) {
	return parser.m11(input, here)
}
