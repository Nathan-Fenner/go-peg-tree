package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nathan-fenner/go-peg-tree/core"
)

type Build interface {
	Build(map[string]string) (core.Peg, error)
}

type BuildRoot struct {
	Root string
}

func (build BuildRoot) Build(roots map[string]string) (core.Peg, error) {
	if returns, ok := roots[build.Root]; ok {
		return core.Root{build.Root, returns}, nil
	}
	return nil, fmt.Errorf("root `%s` is not defined", build.Root)
}

type BuildStar struct {
	Argument Build
}

func (build BuildStar) Build(roots map[string]string) (core.Peg, error) {
	contents, err := build.Argument.Build(roots)
	if err != nil {
		return nil, err
	}
	return core.Star{contents}, nil
}

type BuildPlus struct {
	Argument Build
}

func (build BuildPlus) Build(roots map[string]string) (core.Peg, error) {
	contents, err := build.Argument.Build(roots)
	if err != nil {
		return nil, err
	}
	return core.Plus{contents}, nil
}

type BuildOptional struct {
	Argument Build
}

func (build BuildOptional) Build(roots map[string]string) (core.Peg, error) {
	contents, err := build.Argument.Build(roots)
	if err != nil {
		return nil, err
	}
	return core.Optional{contents}, nil
}

func buildUnit(b Build, suffix string) Build {
	switch suffix {
	case "*":
		return BuildStar{b}
	case "+":
		return BuildPlus{b}
	case "?":
		return BuildOptional{b}
	}
	panic(`buildUnit must be given "*" or "?" or "+"`)
}

type BuildSequence []Build

type ErrorSequence []error

func (err ErrorSequence) Error() string {
	s := ""
	for i := range err {
		s += err[i].Error() + "\n"
	}
	return s
}

func (build BuildSequence) Build(roots map[string]string) (core.Peg, error) {
	if len(build) == 0 {
		panic("sequences must be non-empty")
	}
	result := make([]core.Peg, len(build))
	errs := ErrorSequence{}
	for i := range build {
		peg, err := build[i].Build(roots)
		result[i] = peg
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return nil, errs
	}
	return core.Sequence(result), nil
}

type BuildAlternate []Build

func (build BuildAlternate) Build(roots map[string]string) (core.Peg, error) {
	if len(build) == 0 {
		panic("alternates but be non-empty")
	}
	result := make([]core.Peg, len(build))
	errs := ErrorSequence{}
	correctType := ""
	correctIndex := -1
	for i := range build {
		peg, err := build[i].Build(roots)
		result[i] = peg
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if correctIndex == -1 {
			correctType = peg.TypeName()
			correctIndex = i
		} else if peg.TypeName() != correctType {
			errs = append(errs, fmt.Errorf("alternate members should have consistent types, but %d :: %s while %d :: %s", correctIndex, correctType, i, result[i].TypeName()))
		}
	}
	if len(errs) != 0 {
		return nil, errs
	}
	return core.Alternate(result), nil
}

type BuildRegex string

func (build BuildRegex) Build(roots map[string]string) (core.Peg, error) {
	pattern := string(build)
	_, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return core.Regex{pattern}, nil
}

type BuildGo struct {
	Argument   Build
	Returns    string
	Expression string
}

func (build BuildGo) Build(roots map[string]string) (core.Peg, error) {
	contents, err := build.Argument.Build(roots)
	if err != nil {
		return nil, err
	}
	return core.Go{contents, build.Returns, build.Expression}, nil
}

func buildGo(build Build, buildGo *BuildGo) Build {
	if buildGo == nil {
		return build
	}
	return BuildGo{build, buildGo.Returns, buildGo.Expression}
}

type Rule struct {
	Name    string
	Returns string
	Right   Build
}

func unescapeString(s string) string {
	s = strings.Replace(s, `\"`, `"`, -1)
	s = strings.Replace(s, `\n`, "\n", -1)
	s = strings.Replace(s, `\t`, "\t", -1)
	s = strings.Replace(s, `\v`, "\v", -1)
	s = strings.Replace(s, `\b`, "\b", -1)
	s = strings.Replace(s, `\\`, "\\", -1)
	return s
}
