package main

import (
	"fmt"

	"github.com/nathan-fenner/go-peg-tree/core"
)

func main() {
	state := core.NewState()
	state.DefineRoot("identifier", core.Regex{Regex: `[a-zA-Z_][a-zA-Z_0-9]*`})
	fmt.Print(state.Generate("self"))
}
