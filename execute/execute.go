package execute

import (
	"github.com/adam-bunce/morpheus/backend"
	parser "github.com/adam-bunce/morpheus/generated"
	"github.com/antlr4-go/antlr/v4"
)

func RunProgram(source string) backend.Runtime {
	cs := antlr.NewInputStream(source)
	lexer := parser.NewmorpheusLexer(cs)
	tokens := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewmorpheusParser(tokens)
	rt := backend.NewRuntime()

	result := p.Program()
	result.GetStatements().Eval(rt)

	return rt
}
