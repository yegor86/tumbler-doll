package workflow

import (
	"log"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
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

type (
	// Pipeline represents the main Jenkins pipeline structure
	Pipeline struct {
		Agent  *Agent   `"pipeline" "{" "agent" @@`
		Stages []*Stage `"stages" "{" @@+ "}"`
		Close  string   `"}"`
	}

	// Agent represents the agent block in a Jenkinsfile
	Agent struct {
		None   bool    `( "none" )?`
		Docker *Docker `( "{" "docker" @@ "}" )?`
	}

	Docker struct {
		Image QuotedString `@String`
	}

	Parallel []*Stage

	// Stage represents a stage block within stages
	Stage struct {
		Name     QuotedString `"stage" "(" @String ")" "{"`
		Agent    *Agent       `( "agent" @@ )?`
		Steps    []*Step      `( "steps" "{" @@+ "}" )?`
		FailFast *bool        `( "failFast" @Bool )?`
		Parallel Parallel     `( "parallel" "{" @@+ "}" )?`
		Close    string       `"}"`
	}

	// Step represents individual steps within a stage
	Step struct {
		SingleKV *SingleKVCommand `@@ |`
		MultiKV  *MultiKVCommand  `@@`
	}

	SingleKVCommand struct {
		Command string       `@Ident`
		Value   QuotedString `@String`
	}

	MultiKVCommand struct {
		Command string  `@Ident`
		Params  []Param `@@ ("," @@)*`
	}

	Param struct {
		Key   string       `@Ident ":"`
		Value QuotedString `@String`
	}
)

type QuotedString string

// Capture method strips quotes from the Image field
func (o *QuotedString) Capture(values []string) error {
	*o = QuotedString(strings.Trim(values[0], `"'`))
	return nil
}

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
