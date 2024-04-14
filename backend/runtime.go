package backend

import (
	"fmt"
	"github.com/adam-bunce/morpheus/util"
	"github.com/lithdew/casso"
	"strings"
)

type Runtime struct {
	SymbolTable map[string]Data
	Solver      *casso.Solver
}

func NewRuntime() Runtime {
	return Runtime{
		SymbolTable: map[string]Data{},
		Solver:      casso.NewSolver(),
	}
}

func (r Runtime) String() string {
	var sb strings.Builder

	sb.WriteString("[")

	for name, data := range r.SymbolTable {
		sb.WriteString(fmt.Sprintf(" %v=%s", name, data))
	}

	sb.WriteString(" ]")

	return sb.String()
}

// SubScope returns a copy of the parents runtime with new bindings
func (r Runtime) SubScope(bindings map[string]Data) Runtime {
	newRuntime := Runtime{
		SymbolTable: util.DeepCopyMap(r.SymbolTable),
	}

	for name, data := range bindings {
		newRuntime.SymbolTable[name] = data
	}

	return newRuntime
}
