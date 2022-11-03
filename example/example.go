package main

import (
	"fmt"

	"github.com/josebalius/exhauststruct/example/pkg"
)

func main() {
	x := i(0)
	t := myStruct{field: ""}
	v := pkg.SpecialType{}
	a := pkg.AnotherType{}
	fmt.Println(t)
	fmt.Println(x)
	fmt.Println(v)
	fmt.Println(a)
}

// test does
type test struct{}

// myStruct is a struct
//lint:exhauststruct
type myStruct struct {
	field  string
	field2 string
	field3 int
}

// Something else
type i int
