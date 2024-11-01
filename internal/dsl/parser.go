package dsl

import (
	"log"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Define the lexer rules for Jenkinsfile syntax
var lexerRules = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Keyword", Pattern: `\b(pipeline|agent|docker|image|stages|stage|steps|sh|echo|none)\b`},
	{Name: "String", Pattern: `'(.*?)'`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Punctuation", Pattern: `[{}()]`},
	{Name: "whitespace", Pattern: `\s+`},
})

type (
	DslParser struct{}

	// Pipeline represents the main Jenkins pipeline structure
	Pipeline struct {
		Agent  *Agent   `"pipeline" "{" "agent" ("none" | "{" @@ "}")`
		Stages []*Stage `"stages" "{" @@+ "}"`
		Close  string   `"}"`
	}

	// Agent represents the agent block in a Jenkinsfile
	Agent struct {
		Docker *Docker `"docker" ("{" @@ "}") | "docker" @@`
	}

	// Docker represents the Docker configuration within an agent
	Docker struct {
		Image string `"image" @String | @String`
	}

	// Stage represents a stage block within stages
	Stage struct {
		Name  string  `"stage" "(" @String ")" "{"`
		Agent *Agent  `"agent" ("none" | "{" @@ "}")`
		Steps []*Step `"steps" "{" @@+ "}"`
		Close string  `"}"`
	}

	// Step represents individual steps within a stage
	Step struct {
		Command string `"sh" @String | @String`
	}
)

func (*DslParser) Parse(dslFile string) (*Pipeline, error) {
	parser := participle.MustBuild[Pipeline](
		participle.Lexer(lexerRules),
	)

	pipeline, err := parser.ParseString("", dslFile)
	if err != nil {
		log.Printf("Error parsing Groovy DSL: %v\n", err)
		return nil, err
	}
	return pipeline, nil
}
