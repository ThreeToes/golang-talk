package presentforcer
// go mod tidy keeps deleting the tools package out of go.mod, which breaks stuff
// this is just here to force the mod system to keep it since it's not respecting //indirect

import (
	"fmt"
	"golang.org/x/tools/present"
)

func OverlyAggressiveModTidy() {
	fmt.Println(present.Author{})
}
