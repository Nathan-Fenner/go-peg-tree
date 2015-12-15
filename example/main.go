package main

import (
	"fmt"

	"github.com/nathan-fenner/go-peg-tree/core"
)

func main() {
	state := core.NewState()
	state.DefineRoot("one", core.Go{core.Literal("one"), "float64", "1"})
	state.DefineRoot("two", core.Go{core.Literal("two"), "float64", "2"})
	state.DefineRoot("three", core.Go{core.Literal("three"), "float64", "3"})
	state.DefineRoot("four", core.Go{core.Literal("four"), "float64", "4"})
	state.DefineRoot("number", core.Alternate{
		core.Root{"one", "float64"},
		core.Root{"two", "float64"},
		core.Root{"three", "float64"},
		core.Root{"four", "float64"},
	})

	state.DefineRoot("sum", core.Alternate{
		core.Go{
			Argument: core.Sequence{
				core.Root{"number", "float64"},
				core.Literal("+"),
				core.Root{"sum", "float64"},
			},
			Returns:    "float64",
			Expression: "arg.V0 + arg.V2",
		},
		core.Root{"number", "float64"},
	})
	state.DefineRoot("Expression", core.Root{"sum", "float64"})

	fmt.Print(state.Generate("arithmetic"))
}
