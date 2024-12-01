package dsl

import (
	"log"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/yegor86/tumbler-doll/internal/workflow"
)

type (
	DslParser struct{}
)

// Define the lexer rules for Jenkinsfile syntax
var lexerRules = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Keyword", Pattern: `\b(pipeline|agent|docker|stages|stage|steps|none|failFast)\b`},
	{Name: "String", Pattern: `'([^']*)'|"([^"]*)"`},
	{Name: "Bool", Pattern: `true|false`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Punctuation", Pattern: `[{}()]`},
	{Name: "whitespace", Pattern: `\s+`},
	{Name: "comment", Pattern: `\/\/[^\n]*`},
	{Name: "Colon", Pattern: `:`},
	{Name: "Comma", Pattern: `,`},
})

func (*DslParser) Parse(dslFile string) (*workflow.Pipeline, error) {
	parser := participle.MustBuild[workflow.Pipeline](
		participle.Lexer(lexerRules),
	)

	pipeline, err := parser.ParseString("", dslFile)
	if err != nil {
		log.Printf("Error parsing Groovy DSL: %v\n", err)
		return nil, err
	}
	return pipeline, nil
}
