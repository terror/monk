package main

import (
  "fmt"
  "os"
  "strings"
)

func EvalFile(filename string) (Object, error) {
  data, err := os.ReadFile(filename)
  if err != nil {
    return nil, fmt.Errorf("error reading file: %w", err)
  }

  env := NewEnvironment()
  parser := NewParser(NewLexer(string(data)))
  program := parser.Parse()

  if len(parser.Errors()) != 0 {
    var errMsg strings.Builder
    errMsg.WriteString(fmt.Sprintf("parser errors in file %s:\n", filename))
    for _, msg := range parser.Errors() {
      errMsg.WriteString(fmt.Sprintf("\t%s\n", msg))
    }
    return nil, fmt.Errorf(errMsg.String())
  }

  result := Eval(program, env)
  return result, nil
}

func main() {
  args := os.Args[1:]

  if len(args) == 0 {
    fmt.Println("Monk programming language REPL")
    fmt.Println("Type in commands to evaluate them")
    Repl(os.Stdin, os.Stdout)
  } else {
    filename := args[0]

    if !strings.HasSuffix(filename, ".monk") {
      fmt.Printf("Error: File must have .monk extension\n")
      os.Exit(1)
    }

    result, err := EvalFile(filename)
    if err != nil {
      fmt.Printf("Error: %s\n", err)
      os.Exit(1)
    }

    if result != nil && result != NULL_LIT {
      fmt.Println(result.Inspect())
    }
  }
}
