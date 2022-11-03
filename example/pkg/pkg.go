package pkg

import "fmt"

//lint:exhauststruct
type SpecialType struct {
	Hello string
	World string
}

func init() {
	v := SpecialType{}
	fmt.Println(v)
}
