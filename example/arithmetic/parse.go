package arithmetic

import "fmt"

func ParseExpression(input string) (float64, error) {
	check, value := m19([]rune(input), 0)
	if check.Ok {
		return value, nil
	}
	var zero float64
	return zero, fmt.Errorf("%s", check.Explain())
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

func m0(input []rune, here int) (Result, float64) {
	return m2(input, here)
}

var wherem1 = map[int]Result{}
var whatm1 = map[int]string{}

func m1(input []rune, here int) (Result, string) {
	if result, ok := wherem1[here]; ok {
		return result, whatm1[here]
	}
	result, value := dm1(input, here)
	wherem1[here] = result
	whatm1[here] = value
	return result, value
}

// "one"
func dm1(input []rune, here int) (Result, string) {
	if here+3 > len(input) || string(input[here:here+3]) != "one" {
		return Failure(Expected{Token: "one"}), ""
	}
	return Success(here + 3), "one"
}

var wherem10 = map[int]Result{}
var whatm10 = map[int]string{}

func m10(input []rune, here int) (Result, string) {
	if result, ok := wherem10[here]; ok {
		return result, whatm10[here]
	}
	result, value := dm10(input, here)
	wherem10[here] = result
	whatm10[here] = value
	return result, value
}

// "four"
func dm10(input []rune, here int) (Result, string) {
	if here+4 > len(input) || string(input[here:here+4]) != "four" {
		return Failure(Expected{Token: "four"}), ""
	}
	return Success(here + 4), "four"
}

var wherem11 = map[int]Result{}
var whatm11 = map[int]float64{}

func m11(input []rune, here int) (Result, float64) {
	if result, ok := wherem11[here]; ok {
		return result, whatm11[here]
	}
	result, value := dm11(input, here)
	wherem11[here] = result
	whatm11[here] = value
	return result, value
}

// string go { float64 }
func dm11(input []rune, here int) (Result, float64) {
	check, value := m10(input, here)
	if !check.Ok {
		var zero float64
		return check, zero
	}
	answer := func(arg string) float64 {
		return 4
	}(value)
	return check, answer
}

func m12(input []rune, here int) (Result, float64) {
	return m13(input, here)
}

var wherem13 = map[int]Result{}
var whatm13 = map[int]float64{}

func m13(input []rune, here int) (Result, float64) {
	if result, ok := wherem13[here]; ok {
		return result, whatm13[here]
	}
	result, value := dm13(input, here)
	wherem13[here] = result
	whatm13[here] = value
	return result, value
}

// (root one / root two / root three / root four)
func dm13(input []rune, here int) (Result, float64) {
	notes := []Reject{}
	if next, value := m0(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	if next, value := m3(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	if next, value := m6(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	if next, value := m9(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	var zero float64
	return Failure(notes...), zero
}

func m14(input []rune, here int) (Result, float64) {
	return m18(input, here)
}

var wherem15 = map[int]Result{}
var whatm15 = map[int]string{}

func m15(input []rune, here int) (Result, string) {
	if result, ok := wherem15[here]; ok {
		return result, whatm15[here]
	}
	result, value := dm15(input, here)
	wherem15[here] = result
	whatm15[here] = value
	return result, value
}

// "+"
func dm15(input []rune, here int) (Result, string) {
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

func m16(input []rune, here int) (Result, struct {
	V0 float64
	V1 string
	V2 float64
}) {
	if result, ok := wherem16[here]; ok {
		return result, whatm16[here]
	}
	result, value := dm16(input, here)
	wherem16[here] = result
	whatm16[here] = value
	return result, value
}

// root number "+" root sum
func dm16(input []rune, here int) (Result, struct {
	V0 float64
	V1 string
	V2 float64
}) {
	result := struct {
		V0 float64
		V1 string
		V2 float64
	}{}
	if next, value := m12(input, here); next.Ok {
		here = next.At
		result.V0 = value
	} else {
		return next, struct {
			V0 float64
			V1 string
			V2 float64
		}{}
	}
	if next, value := m15(input, here); next.Ok {
		here = next.At
		result.V1 = value
	} else {
		return next, struct {
			V0 float64
			V1 string
			V2 float64
		}{}
	}
	if next, value := m14(input, here); next.Ok {
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

func m17(input []rune, here int) (Result, float64) {
	if result, ok := wherem17[here]; ok {
		return result, whatm17[here]
	}
	result, value := dm17(input, here)
	wherem17[here] = result
	whatm17[here] = value
	return result, value
}

// struct{V0 float64;V1 string;V2 float64;} go { float64 }
func dm17(input []rune, here int) (Result, float64) {
	check, value := m16(input, here)
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

func m18(input []rune, here int) (Result, float64) {
	if result, ok := wherem18[here]; ok {
		return result, whatm18[here]
	}
	result, value := dm18(input, here)
	wherem18[here] = result
	whatm18[here] = value
	return result, value
}

// (struct{V0 float64;V1 string;V2 float64;} go { float64 } / root number)
func dm18(input []rune, here int) (Result, float64) {
	notes := []Reject{}
	if next, value := m17(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	if next, value := m12(input, here); next.Ok {
		return next, value
	} else {
		notes = append(notes, next.Expected...)
	}
	var zero float64
	return Failure(notes...), zero
}

func m19(input []rune, here int) (Result, float64) {
	return m14(input, here)
}

var wherem2 = map[int]Result{}
var whatm2 = map[int]float64{}

func m2(input []rune, here int) (Result, float64) {
	if result, ok := wherem2[here]; ok {
		return result, whatm2[here]
	}
	result, value := dm2(input, here)
	wherem2[here] = result
	whatm2[here] = value
	return result, value
}

// string go { float64 }
func dm2(input []rune, here int) (Result, float64) {
	check, value := m1(input, here)
	if !check.Ok {
		var zero float64
		return check, zero
	}
	answer := func(arg string) float64 {
		return 1
	}(value)
	return check, answer
}

func m3(input []rune, here int) (Result, float64) {
	return m5(input, here)
}

var wherem4 = map[int]Result{}
var whatm4 = map[int]string{}

func m4(input []rune, here int) (Result, string) {
	if result, ok := wherem4[here]; ok {
		return result, whatm4[here]
	}
	result, value := dm4(input, here)
	wherem4[here] = result
	whatm4[here] = value
	return result, value
}

// "two"
func dm4(input []rune, here int) (Result, string) {
	if here+3 > len(input) || string(input[here:here+3]) != "two" {
		return Failure(Expected{Token: "two"}), ""
	}
	return Success(here + 3), "two"
}

var wherem5 = map[int]Result{}
var whatm5 = map[int]float64{}

func m5(input []rune, here int) (Result, float64) {
	if result, ok := wherem5[here]; ok {
		return result, whatm5[here]
	}
	result, value := dm5(input, here)
	wherem5[here] = result
	whatm5[here] = value
	return result, value
}

// string go { float64 }
func dm5(input []rune, here int) (Result, float64) {
	check, value := m4(input, here)
	if !check.Ok {
		var zero float64
		return check, zero
	}
	answer := func(arg string) float64 {
		return 2
	}(value)
	return check, answer
}

func m6(input []rune, here int) (Result, float64) {
	return m8(input, here)
}

var wherem7 = map[int]Result{}
var whatm7 = map[int]string{}

func m7(input []rune, here int) (Result, string) {
	if result, ok := wherem7[here]; ok {
		return result, whatm7[here]
	}
	result, value := dm7(input, here)
	wherem7[here] = result
	whatm7[here] = value
	return result, value
}

// "three"
func dm7(input []rune, here int) (Result, string) {
	if here+5 > len(input) || string(input[here:here+5]) != "three" {
		return Failure(Expected{Token: "three"}), ""
	}
	return Success(here + 5), "three"
}

var wherem8 = map[int]Result{}
var whatm8 = map[int]float64{}

func m8(input []rune, here int) (Result, float64) {
	if result, ok := wherem8[here]; ok {
		return result, whatm8[here]
	}
	result, value := dm8(input, here)
	wherem8[here] = result
	whatm8[here] = value
	return result, value
}

// string go { float64 }
func dm8(input []rune, here int) (Result, float64) {
	check, value := m7(input, here)
	if !check.Ok {
		var zero float64
		return check, zero
	}
	answer := func(arg string) float64 {
		return 3
	}(value)
	return check, answer
}

func m9(input []rune, here int) (Result, float64) {
	return m11(input, here)
}
