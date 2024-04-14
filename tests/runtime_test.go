package tests

import (
	"github.com/adam-bunce/morpheus/backend"
	"testing"
)

func TestRuntimeSymbolTable(t *testing.T) {
	runtime := backend.Runtime{
		SymbolTable: map[string]backend.Data{},
	}

	runtime.SymbolTable["boolean_var"] = backend.BooleanData{
		Value:   false,
		Literal: "false",
	}

	runtime.SymbolTable["greeting"] = backend.StringData{
		Value:   "hello world",
		Literal: "\"hello world\"",
	}

	runtime.SymbolTable["ten"] = backend.IntData{
		Value:   10,
		Literal: "10",
	}

	if runtime.SymbolTable["boolean_var"].(backend.BooleanData).Value != false {
		t.Fatalf("symbol table should contain `%s` with value of `%t` got: %v", "boolean_var", true, runtime)
	}

	if runtime.SymbolTable["greeting"].(backend.StringData).Value != "hello world" {
		t.Fatalf("symbol table should contain `%s` with value of `%s` got: %v", "greeting", "hello world", runtime)
	}

	if runtime.SymbolTable["ten"].(backend.IntData).Value != 10 {
		t.Fatalf("symbol table should contain `%s` with value of `%d` got: %v", "ten", 10, runtime)
	}
}

func TestRuntimeSubscope(t *testing.T) {
	runtime := backend.Runtime{
		SymbolTable: map[string]backend.Data{},
	}

	runtime.SymbolTable["boolean_var"] = backend.BooleanData{
		Value:   false,
		Literal: "false",
	}

	runtime.SymbolTable["greeting"] = backend.StringData{
		Value:   "hello world",
		Literal: "\"hello world\"",
	}

	runtime.SymbolTable["ten"] = backend.IntData{
		Value:   10,
		Literal: "10",
	}

	newRuntime := runtime.SubScope(map[string]backend.Data{
		"greeting": backend.IntData{
			Value:   1337,
			Literal: "1337",
		},
	})

	// new runtime should contain exact same, but with greeting being intdata now
	if newRuntime.SymbolTable["boolean_var"].(backend.BooleanData).Value != false {
		t.Fatalf("symbol table should contain `%s` with value of `%t` got: %v", "boolean_var", true, newRuntime)
	}

	if newRuntime.SymbolTable["greeting"].(backend.IntData).Value != 1337 {
		t.Fatalf("symbol table should contain `%s` with value of `%s` got: %v", "greeting", "hello world", newRuntime)
	}

	if newRuntime.SymbolTable["ten"].(backend.IntData).Value != 10 {
		t.Fatalf("symbol table should contain `%s` with value of `%d` got: %v", "ten", 10, newRuntime)
	}

	// parent should have old values
	if runtime.SymbolTable["boolean_var"].(backend.BooleanData).Value != false {
		t.Fatalf("symbol table should contain `%s` with value of `%t` got: %v", "boolean_var", true, runtime)
	}

	if runtime.SymbolTable["greeting"].(backend.StringData).Value != "hello world" {
		t.Fatalf("symbol table should contain `%s` with value of `%s` got: %v", "greeting", "hello world", runtime)
	}

	if runtime.SymbolTable["ten"].(backend.IntData).Value != 10 {
		t.Fatalf("symbol table should contain `%s` with value of `%d` got: %v", "ten", 10, runtime)
	}

}
