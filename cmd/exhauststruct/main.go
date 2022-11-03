package main

import (
	"github.com/josebalius/exhauststruct/exhauststruct"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(exhauststruct.Analyzer)
}
