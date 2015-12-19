package main

import (
	"fmt"

	"github.com/nathan-fenner/go-peg-tree/core"
)

func main() {
	state := core.NewState()
	state.DefineRoot(
		"identifier",
		core.Go{
			core.Sequence{core.Root{"space", "string"}, core.Regex{Regex: `[\p{L}_][\p{L}\d_-]*`}},
			"string",
			"arg.V1",
		},
	)
	state.DefineRoot(
		"keyword",
		core.Not{core.Regex{`[\p{L}\d_]`}},
	)
	state.DefineRoot("space", core.Regex{`\s*`})
	state.DefineRoot("mandatory-space", core.Regex{`\s+`})

	state.DefineRoot("string-backtick", core.Go{
		core.Regex{"`[^`]*`"},
		"string",
		"arg[1:len(arg)-1]",
	})
	state.DefineRoot("string-quote", core.Go{
		core.Regex{`"([^\\"\n]|\\["ntvb\\])*"`}, // TODO: other escapes
		"string",
		"unescapeString(arg[1:len(arg)-1])",
	})
	state.DefineRoot("string-literal", core.Alternate{
		core.Root{"string-backtick", "string"},
		core.Root{"string-quote", "string"},
	})

	state.DefineRoot("type-identifier", core.Root{"identifier", "string"})
	state.DefineRoot("type-head", core.Alternate{
		core.Literal("*"),  // TODO: strip preceding spaces
		core.Literal("[]"), // TODO: strip preceding spaces
		core.Contents{core.Sequence{core.Literal("["), core.Regex{`\d+`}, core.Literal("]")}},                                // TODO: strip preceding spaces
		core.Contents{core.Sequence{core.Literal("map"), core.Literal("["), core.Root{"type", "string"}, core.Literal("]")}}, // TODO: strip preceding spaces
		core.Go{core.Sequence{core.Literal("chan"), core.Root{"keyword", "struct{}"}}, "string", `"chan "`},
		core.Literal("<-chan "), // TODO: use keyword
		core.Literal("chan<- "), // TODO: use keyword
	})
	state.DefineRoot("type", core.Contents{core.Sequence{core.Star{core.Root{"type-head", "string"}}, core.Root{"identifier", "string"}}})

	state.DefineRoot("peg-root", core.Go{
		core.Root{"identifier", "string"},
		"Build",
		"BuildRoot{arg}",
	})

	state.DefineRoot("peg-regex", core.Go{
		core.Sequence{
			core.Root{"space", "string"},
			core.Literal("regex"),
			core.Root{"keyword", "string"},
			core.Root{"space", "string"},
			core.Root{"string-literal", "string"},
		},
		"Build",
		"BuildRegex{arg[4]}",
	})

	state.DefineRoot(
		"peg-atom",
		core.Alternate{
			core.Root{"peg-root", "Build"},
			core.Root{"peg-regex", "Build"},
			core.Go{
				core.Sequence{core.Root{"space", "string"}, core.Literal("("), core.Root{"peg-expression", "Build"}, core.Root{"space", "string"}, core.Literal(")")},
				"Build",
				"arg.V2",
			},
		},
	)

	state.DefineRoot(
		"peg-unit-suffix",
		core.Go{
			core.Sequence{
				core.Root{"space", "string"},
				core.Alternate{
					core.Literal("*"),
					core.Literal("+"),
					core.Literal("?"),
				},
			},
			"string",
			"arg.V1",
		},
	)

	state.DefineRoot(
		"peg-unit",
		core.Go{
			core.Sequence{
				core.Root{"peg-atom", "Build"},
				core.Root{"peg-unit-suffix", "string"},
			},
			"Build",
			"buildUnit(arg.V0, arg.V1)",
		},
	)

	state.DefineRoot(
		"peg-sequence",
		core.Go{
			core.Plus{core.Root{"peg-unit", "Build"}},
			"Build",
			"BuildSequence(arg)",
		},
	)

	state.DefineRoot(
		"peg-go",
		core.Go{
			core.Sequence{
				core.Root{"peg-sequence", "Build"},
				core.Optional{
					core.Go{
						core.Sequence{
							core.Root{"space", "string"},
							core.Literal("go"),
							core.Root{"keyword", "struct{}"},
							core.Root{"type", "string"},
							core.Root{"space", "string"},
							core.Literal("{"),
							core.Contents{core.Regex{`[^{}]+`}},
							core.Literal("}"),
						},
						"BuildGo",
						"BuildGo{nil, arg.V3, arg.V6}",
					},
				},
			},
			"Build",
			"buildGo(arg.V0, arg.V1)",
		},
	)

	state.DefineRoot(
		"peg-alternate",
		core.Go{
			core.Sequence{
				core.Root{"peg-go", "Build"},
				core.Star{
					core.Go{
						core.Sequence{
							core.Root{"space", "string"},
							core.Literal("|"),
							core.Root{"peg-go", "Build"},
						},
						"Build",
						"arg.V2",
					},
				},
			},
			"Build",
			"BuildAlternate(append([]Build{arg.V0}, arg.V1...))",
		},
	)

	state.DefineRoot("peg-expression", core.Root{"peg-alternate", "Build"})

	state.DefineRoot(
		"rule",
		core.Go{
			core.Sequence{
				core.Root{"identifier", "string"},
				core.Root{"type", "string"},
				core.Root{"space", "string"},
				core.Literal("<-"),
				core.Root{"peg-expression", "Build"},
				core.Root{"space", "string"},
				core.Literal(";"),
			},
			"Rule",
			"Rule{arg.V0, arg.V1, arg.V4}",
		},
	)

	state.DefineRoot(
		"Rules",
		core.Plus{core.Root{"rule", "Rule"}},
	)

	fmt.Print(state.Generate("main"))
}
