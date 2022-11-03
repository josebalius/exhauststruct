package main

import (
	"github.com/josebalius/exhauststruct/exhauststruct"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		exhauststruct.Analyzer,
	}
}

var AnalyzerPlugin analyzerPlugin
