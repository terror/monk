package main

import (
  "bufio"
  "fmt"
  "io"
  "strings"
)

const PROMPT = ">> "

func Repl(in io.Reader, out io.Writer) {
  scanner := bufio.NewScanner(in)
  env := NewEnvironment()

  for {
    fmt.Print(PROMPT)

    scanned := scanner.Scan()

    if !scanned {
      return
    }

    line := scanner.Text()

    // Allow exiting the REPL with 'exit' or 'quit'
    if line == "exit" || line == "quit" {
      fmt.Println("Goodbye!")
      return
    }

    // Ignore empty lines
    if strings.TrimSpace(line) == "" {
      continue
    }

    parser := NewParser(NewLexer(line))

    program := parser.Parse()

    if len(parser.Errors()) != 0 {
      for _, message := range parser.Errors() {
        fmt.Println(message)
      }
      continue
    }

    evaluated := Eval(program, env)

    if evaluated != nil {
      io.WriteString(out, evaluated.Inspect())
      io.WriteString(out, "\n")
    }
  }
}
