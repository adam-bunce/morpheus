package main

import (
	"fmt"
	exec "github.com/adam-bunce/morpheus/execute"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./morpheus <program.mph>")
		os.Exit(1)
	}
	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("couldn't find %s...", fileName)
		os.Exit(1)
	}
	defer file.Close()

	program, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("failed to read %s\n", fileName)
		os.Exit(1)
	}

	exec.RunProgram(string(program))
}
